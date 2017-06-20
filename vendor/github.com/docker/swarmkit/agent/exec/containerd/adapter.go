package containerd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	containersapi "github.com/containerd/containerd/api/services/containers"
	contentapi "github.com/containerd/containerd/api/services/content"
	"github.com/containerd/containerd/api/services/execution"
	imagesapi "github.com/containerd/containerd/api/services/images"
	"github.com/containerd/containerd/api/types/mount"
	"github.com/containerd/containerd/api/types/task"
	"github.com/containerd/containerd/archive"
	"github.com/containerd/containerd/archive/compression"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	contentservice "github.com/containerd/containerd/services/content"
	imagesservice "github.com/containerd/containerd/services/images"
	"github.com/containerd/fifo"
	dockermount "github.com/docker/docker/pkg/mount"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	"github.com/docker/swarmkit/log"
	protobuf "github.com/gogo/protobuf/types"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	mountPropagationReverseMap = map[api.Mount_BindOptions_MountPropagation]string{
		api.MountPropagationPrivate:  "private",
		api.MountPropagationRPrivate: "rprivate",
		api.MountPropagationShared:   "shared",
		api.MountPropagationRShared:  "rshared",
		api.MountPropagationRSlave:   "slave",
		api.MountPropagationSlave:    "rslave",
	}
)

// containerAdapter conducts remote operations for a container. All calls
// are mostly naked calls to the client API, seeded with information from
// containerConfig.
type containerAdapter struct {
	conn              *grpc.ClientConn
	taskClient        execution.TasksClient
	container         *api.ContainerSpec
	task              *api.Task
	secrets           exec.SecretGetter
	dir               string
	resolvedImageName string
	deleteResponse    *execution.DeleteResponse
}

func newContainerAdapter(conn *grpc.ClientConn, containerDir string, task *api.Task, secrets exec.SecretGetter) (*containerAdapter, error) {
	container := task.Spec.GetContainer()
	if container == nil {
		return nil, exec.ErrRuntimeUnsupported
	}

	dir := filepath.Join(containerDir, task.ID)

	return &containerAdapter{
		conn:       conn,
		taskClient: execution.NewTasksClient(conn),
		container:  container,
		task:       task,
		secrets:    secrets,
		dir:        dir,
	}, nil
}

func (c *containerAdapter) applyLayer(ctx context.Context, cs content.Store, rootfs string, layer digest.Digest) error {
	blob, err := cs.Reader(ctx, layer)
	if err != nil {
		return err
	}

	rd, err := compression.DecompressStream(blob)
	if err != nil {
		return err
	}

	_, err = archive.Apply(ctx, rootfs, rd)

	blob.Close()
	return err
}

// github.com/containerd/containerd cmd/ctr/utils.go, dropped stdin handling
func prepareStdio(stdout, stderr string, console bool) (wg *sync.WaitGroup, err error) {
	wg = &sync.WaitGroup{}
	ctx := context.Background()

	f, err := fifo.OpenFifo(ctx, stdout, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700)
	if err != nil {
		return nil, err
	}
	defer func(c io.Closer) {
		if err != nil {
			c.Close()
		}
	}(f)
	wg.Add(1)
	go func(r io.ReadCloser) {
		io.Copy(os.Stdout, r)
		r.Close()
		wg.Done()
	}(f)

	f, err = fifo.OpenFifo(ctx, stderr, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700)
	if err != nil {
		return nil, err
	}
	defer func(c io.Closer) {
		if err != nil {
			c.Close()
		}
	}(f)
	if !console {
		wg.Add(1)
		go func(r io.ReadCloser) {
			io.Copy(os.Stderr, r)
			r.Close()
			wg.Done()
		}(f)
	}

	return wg, nil
}

func (c *containerAdapter) pullImage(ctx context.Context) error {
	options := docker.ResolverOptions{}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig:       &tls.Config{},
		ExpectContinueTimeout: 5 * time.Second,
	}

	options.Client = &http.Client{
		Transport: tr,
	}

	resolver := docker.NewResolver(options)

	name, desc, err := resolver.Resolve(ctx, c.container.Image)
	if err != nil {
		return errors.Wrap(err, "failed to resolve ref")
	}
	fetcher, err := resolver.Fetcher(ctx, name)
	if err != nil {
		return errors.Wrap(err, "failed to resolve fetcher")
	}
	c.resolvedImageName = name

	content := contentservice.NewStoreFromClient(contentapi.NewContentClient(c.conn))
	imageStore := imagesservice.NewStoreFromClient(imagesapi.NewImagesClient(c.conn))

	if err := imageStore.Put(ctx, name, desc); err != nil {
		return errors.Wrap(err, "put image")
	}

	return images.Dispatch(ctx,
		images.Handlers(
			remotes.FetchHandler(content, fetcher),
			images.ChildrenHandler(content)),
		desc)
}

func (c *containerAdapter) makeAnonVolume(ctx context.Context, target string) (specs.Mount, error) {
	source := filepath.Join(c.dir, "anon-volumes", target)
	if err := os.MkdirAll(source, 0755); err != nil {
		return specs.Mount{}, err
	}

	return specs.Mount{
		Destination: target,
		Type:        "bind",
		Source:      source,
		Options:     []string{"rbind", "rprivate", "rw"},
	}, nil
}

// Somewhat like docker/docker/daemon/oci_linux.go:setMounts
func (c *containerAdapter) setMounts(ctx context.Context, s *specs.Spec, mounts []api.Mount, volumes map[string]struct{}) error {

	userMounts := make(map[string]struct{})
	for _, m := range mounts {
		userMounts[m.Target] = struct{}{}
	}

	// Filter out mounts that are overridden by user supplied mounts
	var defaultMounts []specs.Mount
	_, mountDev := userMounts["/dev"]
	for _, m := range s.Mounts {
		if _, ok := userMounts[m.Destination]; !ok {
			if mountDev && strings.HasPrefix(m.Destination, "/dev/") {
				continue
			}
			defaultMounts = append(defaultMounts, m)
		}
	}

	s.Mounts = defaultMounts
	for _, m := range mounts {
		if !filepath.IsAbs(m.Target) {
			return errors.Errorf("mount %s is not absolute", m.Target)
		}

		for _, cm := range s.Mounts {
			if cm.Destination == m.Target {
				return errors.Errorf("duplicate mount point '%s'", m.Target)
			}
		}

		delete(volumes, m.Target) // volume is no longer anon

		switch m.Type {
		case api.MountTypeTmpfs:
			opts := []string{"noexec", "nosuid", "nodev", "rprivate"}
			if m.TmpfsOptions != nil {
				if m.TmpfsOptions.SizeBytes <= 0 {
					return errors.New("invalid tmpfs size give")
				}
				opts = append(opts, fmt.Sprintf("size=%d", m.TmpfsOptions.SizeBytes))
				opts = append(opts, fmt.Sprintf("mode=%o", m.TmpfsOptions.Mode))
			}
			if m.ReadOnly {
				opts = append(opts, "ro")
			} else {
				opts = append(opts, "rw")
			}

			opts, err := dockermount.MergeTmpfsOptions(opts)
			if err != nil {
				return err
			}

			s.Mounts = append(s.Mounts, specs.Mount{
				Destination: m.Target,
				Type:        "tmpfs",
				Source:      "tmpfs",
				Options:     opts,
			})

		case api.MountTypeVolume:
			if m.Source != "" {
				return errors.Errorf("non-anon volume mounts not implemented, ignoring %v", m)
			}
			if m.VolumeOptions != nil {
				return errors.Errorf("volume mount VolumeOptions not implemented, ignoring %v", m)
			}

			mt, err := c.makeAnonVolume(ctx, m.Target)
			if err != nil {
				return err
			}

			s.Mounts = append(s.Mounts, mt)
			continue

		case api.MountTypeBind:
			opts := []string{"rbind"}
			if m.ReadOnly {
				opts = append(opts, "ro")
			} else {
				opts = append(opts, "rw")
			}

			propagation := "rprivate"
			if m.BindOptions != nil {
				if p, ok := mountPropagationReverseMap[m.BindOptions.Propagation]; ok {
					propagation = p
				} else {
					log.G(ctx).Warningf("unknown bind mount propagation,  using %q", propagation)
				}
			}
			opts = append(opts, propagation)

			mt := specs.Mount{
				Destination: m.Target,
				Type:        "bind",
				Source:      m.Source,
				Options:     opts,
			}

			s.Mounts = append(s.Mounts, mt)
			continue
		}
	}

	for v := range volumes {
		mt, err := c.makeAnonVolume(ctx, v)
		if err != nil {
			return err
		}

		s.Mounts = append(s.Mounts, mt)
	}
	return nil
}

func (c *containerAdapter) spec(ctx context.Context, config *ocispec.ImageConfig, rootfs string) (*specs.Spec, error) {
	caps := []string{
		"CAP_CHOWN",
		"CAP_DAC_OVERRIDE",
		"CAP_FSETID",
		"CAP_FOWNER",
		"CAP_MKNOD",
		"CAP_NET_RAW",
		"CAP_SETGID",
		"CAP_SETUID",
		"CAP_SETFCAP",
		"CAP_SETPCAP",
		"CAP_NET_BIND_SERVICE",
		"CAP_SYS_CHROOT",
		"CAP_KILL",
		"CAP_AUDIT_WRITE",
	}

	// Need github.com/docker/docker/oci.DefaultSpec()
	spec := specs.Spec{
		Version: "1.0.0-rc2-dev",
		Root: specs.Root{
			Path: rootfs,
		},
		Mounts: []specs.Mount{
			{
				Destination: "/proc",
				Type:        "proc",
				Source:      "proc",
				Options:     []string{"nosuid", "noexec", "nodev"},
			},
			{
				Destination: "/dev",
				Type:        "tmpfs",
				Source:      "tmpfs",
				Options:     []string{"nosuid", "strictatime", "mode=755"},
			},
			{
				Destination: "/dev/pts",
				Type:        "devpts",
				Source:      "devpts",
				Options:     []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"},
			},
			{
				Destination: "/sys",
				Type:        "sysfs",
				Source:      "sysfs",
				Options:     []string{"nosuid", "noexec", "nodev", "ro"},
			},
			{
				Destination: "/sys/fs/cgroup",
				Type:        "cgroup",
				Source:      "cgroup",
				Options:     []string{"ro", "nosuid", "noexec", "nodev"},
			},
			{
				Destination: "/dev/mqueue",
				Type:        "mqueue",
				Source:      "mqueue",
				Options:     []string{"nosuid", "noexec", "nodev"},
			},
		},
		Process: specs.Process{
			Cwd: "/",
			Capabilities: &specs.LinuxCapabilities{
				Bounding:    caps,
				Effective:   caps,
				Inheritable: caps,
				Permitted:   caps,
				Ambient:     caps,
			},
			NoNewPrivileges: true,
			Terminal:        false,
		},
		Linux: &specs.Linux{
			Namespaces: []specs.LinuxNamespace{
				{Type: "mount"},
				{Type: "network"},
				{Type: "uts"},
				{Type: "pid"},
				{Type: "ipc"},
			},
		},
	}

	spec.Platform = specs.Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	if config.WorkingDir != "" {
		spec.Process.Cwd = config.WorkingDir
	}
	spec.Process.Env = config.Env

	var args []string
	if len(c.container.Args) > 0 {
		args = c.container.Args
	} else {
		args = config.Cmd
	}

	if len(c.container.Command) > 0 {
		spec.Process.Args = append(c.container.Command, args...)
	} else {
		spec.Process.Args = append(config.Entrypoint, args...)
	}

	log.G(ctx).Debugf("Process args: %v", spec.Process.Args)
	if err := c.setMounts(ctx, &spec, c.container.Mounts, config.Volumes); err != nil {
		return nil, errors.Wrap(err, "failed to set mounts")
	}
	sort.Sort(mounts(spec.Mounts))

	return &spec, nil
}

func (c *containerAdapter) create(ctx context.Context) error {
	if c.resolvedImageName == "" {
		return errors.New("image has not been pulled")
	}

	containers := containersapi.NewContainersClient(c.conn)
	cs := contentservice.NewStoreFromClient(contentapi.NewContentClient(c.conn))
	imageStore := imagesservice.NewStoreFromClient(imagesapi.NewImagesClient(c.conn))

	image, err := imageStore.Get(ctx, c.resolvedImageName)
	if err != nil {
		return errors.Wrap(err, "image get")
	}

	mbytes, err := content.ReadBlob(ctx, cs, image.Target.Digest)
	if err != nil {
		return err
	}

	rootfs := filepath.Join(c.dir, "rootfs")
	// TODO(ijc) support ControllerLogs interface
	stdin := "/dev/null"
	stdout := filepath.Join(c.dir, "stdout")
	stderr := filepath.Join(c.dir, "stderr")

	if err := os.MkdirAll(rootfs, 0755); err != nil {
		return err
	}

	var config ocispec.Image

	var manifest ocispec.Manifest
	if err := json.Unmarshal(mbytes, &manifest); err != nil {
		return errors.Wrap(err, "unmarshalling image manifest")
	}

	bytes, err := content.ReadBlob(ctx, cs, manifest.Config.Digest)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return errors.Wrap(err, "unmarshalling image config")
	}

	for _, layer := range manifest.Layers {
		if err := c.applyLayer(ctx, cs, rootfs, layer.Digest); err != nil {
			return errors.Wrapf(err, "failed to apply layer %s", layer.Digest.String())
		}
	}

	spec, err := c.spec(ctx, &config.Config, rootfs)
	if err != nil {
		return err
	}

	_, err = prepareStdio(stdout, stderr, spec.Process.Terminal)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(spec, "    ", "    ")
	if err != nil {
		return err
	}

	cid := naming.Task(c.task)
	_, err = containers.Create(ctx, &containersapi.CreateContainerRequest{
		Container: containersapi.Container{
			ID: cid,
			Spec: &protobuf.Any{
				TypeUrl: specs.Version,
				Value:   data,
			},
			Runtime: "linux",
		},
	})
	if err != nil {
		return errors.Wrap(err, "creating container")
	}

	_, err = c.taskClient.Create(ctx, &execution.CreateRequest{
		ContainerID: cid,
		Rootfs:      []*mount.Mount{},
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
		Terminal:    spec.Process.Terminal,
	})
	if err != nil {
		return errors.Wrap(err, "creating task")
	}

	return nil
}

func (c *containerAdapter) start(ctx context.Context) error {
	_, err := c.taskClient.Start(ctx, &execution.StartRequest{
		ContainerID: naming.Task(c.task),
	})
	return err
}

func (c *containerAdapter) eventStream(ctx context.Context, id string) (<-chan task.Event, <-chan error, error) {

	var (
		evtch = make(chan task.Event)
		errch = make(chan error)
	)

	return evtch, errch, nil
}

// events issues a call to the events API and returns a channel with all
// events. The stream of events can be shutdown by cancelling the context.
//
// A chan struct{} is returned that will be closed if the event processing
// fails and needs to be restarted.
func (c *containerAdapter) events(ctx context.Context, opts ...grpc.CallOption) (<-chan task.Event, <-chan struct{}, error) {
	id := naming.Task(c.task)

	l := log.G(ctx).WithFields(logrus.Fields{
		"ID": id,
	})

	// TODO(stevvooe): Move this to a single, global event dispatch. For
	// now, we create a connection per container.
	var (
		eventsq = make(chan task.Event)
		closed  = make(chan struct{})
	)

	l.Debugf("waiting on events")

	cl, err := c.taskClient.Events(ctx, &execution.EventsRequest{}, opts...)
	if err != nil {
		l.WithError(err).Errorf("failed to start event stream")
		return nil, nil, err
	}

	go func() {
		defer close(closed)

		for {
			evt, err := cl.Recv()
			if err != nil {
				l.WithError(err).Error("fatal error from events stream")
				return
			}
			if evt.ID != id {
				l.Debugf("Event for a different container %s", evt.ID)
				continue
			}

			select {
			case eventsq <- *evt:
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventsq, closed, nil
}

func (c *containerAdapter) inspect(ctx context.Context) (task.Task, error) {
	id := naming.Task(c.task)
	rsp, err := c.taskClient.Info(ctx, &execution.InfoRequest{ContainerID: id})
	if err != nil {
		return task.Task{}, err
	}
	return *rsp.Task, nil
}

func (c *containerAdapter) shutdown(ctx context.Context) (uint32, error) {
	id := naming.Task(c.task)
	l := log.G(ctx).WithFields(logrus.Fields{
		"ID": id,
	})

	if c.deleteResponse == nil {
		var err error
		l.Debug("Deleting")

		rsp, err := c.taskClient.Delete(ctx, &execution.DeleteRequest{ContainerID: id})
		if err != nil {
			return 0, err
		}
		l.Debugf("Status=%d", rsp.ExitStatus)
		c.deleteResponse = rsp

		containers := containersapi.NewContainersClient(c.conn)
		_, err = containers.Delete(ctx, &containersapi.DeleteContainerRequest{
			ID: id,
		})
		if err != nil {
			l.WithError(err).Warnf("failed to delete container")
		}
	}

	return c.deleteResponse.ExitStatus, nil
}

func (c *containerAdapter) terminate(ctx context.Context) error {
	id := naming.Task(c.task)
	l := log.G(ctx).WithFields(logrus.Fields{
		"ID": id,
	})
	l.Debug("Terminate")
	return errors.New("terminate not implemented")
}

func (c *containerAdapter) remove(ctx context.Context) error {
	id := naming.Task(c.task)
	l := log.G(ctx).WithFields(logrus.Fields{
		"ID": id,
	})
	l.Debug("Remove")
	return os.RemoveAll(c.dir)
}

func isContainerCreateNameConflict(err error) bool {
	// container ".*" already exists
	splits := strings.SplitN(err.Error(), "\"", 3)
	return splits[0] == "container " && splits[2] == " already exists"
}

func isUnknownContainer(err error) bool {
	return strings.Contains(err.Error(), "container does not exist")
}

// For sort.Sort
type mounts []specs.Mount

// Len returns the number of mounts. Used in sorting.
func (m mounts) Len() int {
	return len(m)
}

// Less returns true if the number of parts (a/b/c would be 3 parts) in the
// mount indexed by parameter 1 is less than that of the mount indexed by
// parameter 2. Used in sorting.
func (m mounts) Less(i, j int) bool {
	return m.parts(i) < m.parts(j)
}

// Swap swaps two items in an array of mounts. Used in sorting
func (m mounts) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// parts returns the number of parts in the destination of a mount. Used in sorting.
func (m mounts) parts(i int) int {
	return strings.Count(filepath.Clean(m[i].Destination), string(os.PathSeparator))
}

package main

import (
	"bytes"
	"errors"
	"github.com/blang/semver"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	minimumDockerEngineVersion = "1.12.3"
	dockerAPIMinVersion        = "1.24"
	dockerURL                  = "unix:///var/run/docker.sock"
	dockerGetURL               = "http://get.docker.com"
	dockerRemoteAPIPort        = "2375"
)

var startCommands = [...]string{"systemctl start docker.service", "service docker start", "start docker"}
var enableCommands = [...]string{"systemctl enable docker.service", "chkconfig docker on", "update-rc.d docker defaults"}

func (c *ampManager) getDockerServerVersion() (string, error) {
	// first try the Docker api
	if c.docker == nil {
		defaultHeaders := map[string]string{"User-Agent": "amp-bootstrap"}
		c.docker, _ = dockerclient.NewClient(dockerURL, dockerAPIMinVersion, nil, defaultHeaders)
	}
	if c.docker != nil {
		version, err := c.docker.ServerVersion(context.Background())
		if err == nil {
			return version.Version, nil
		}
	}
	c.printf(colInfo, "Unable to contact the Docker server through the API, switch to the Docker cli\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	versionCmd := [...]string{"docker", "version", "-f", "{{ .Server.Version }}"}
	cmd := exec.Command(versionCmd[0], versionCmd[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		c.printf(colInfo, "%s failed: %s\n", versionCmd, stderr.String())
		return "0.0.0", err
	}
	c.printf(colInfo, "Installed Docker engine version: %s\n", out.String())
	return strings.TrimRight(out.String(), "\n"), nil
}

func (c *ampManager) enableDockerEngine() error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	enabled := false
	c.printf(colInfo, "Enabling Docker engine... \n")
	for _, enableCmd := range enableCommands {
		enableCmdSplitted := strings.Split(enableCmd, " ")
		cmd := exec.Command(enableCmdSplitted[0], enableCmdSplitted[1:]...)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			continue
		}
		enabled = true
	}
	if enabled == false {
		c.printf(colInfo, "all methods failed\n")
		return errors.New("Unable to enable the Docker engine")
	}
	c.printf(colInfo, "done\n")
	return nil
}

// update Docker configuration with pre-requisites
// for now only on systemd compatible systems
func (c *ampManager) configureDockerEngine() error {
	var execstart string
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("systemctl", "show", "--property=ExecStart docker")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`ExecStart={.*argv\[\]=([^;]*);.*`)
	match := re.FindStringSubmatch(out.String())
	if match == nil || len(match) < 2 {
		// can't find an existing start command, rewrite it
		execstart = "/usr/bin/dockerd"
	} else {
		execstart = match[1]
	}
	// add options to the start command
	execstart = execstart + " -H 0.0.0.0:" + dockerRemoteAPIPort + " -H /var/run/docker.sock"

	_, err = os.Stat("/etc/systemd/system")
	if err == nil {
		err := os.MkdirAll("/etc/systemd/system/docker.service.d", 0755)
		if err != nil {
			return err
		}
		file, err := os.Create("/etc/systemd/system/docker.service.d/override.conf")
		if err != nil {
			return err
		}
		defer file.Close()
		content := `[Service]
ExecStart=-
ExecStart=` + execstart + `
LimitNOFILE=-
LimitNOFILE=infinity`
		n, err := file.WriteString(content)
		if err != nil {
			return err
		}
		if n < len(content) {
			return errors.New("unable to override Docker systemd configuration")
		}
	} else {
		return err
	}
	return nil
}

// Try different methods to start the Docker engine
func (c *ampManager) startDockerEngine() error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	c.printf(colInfo, "Starting Docker engine... \n")
	started := false
	for _, startCmd := range startCommands {
		startCmdSplitted := strings.Split(startCmd, " ")
		cmd := exec.Command(startCmdSplitted[0], startCmdSplitted[1:]...)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			// try next method
			continue
		}
		started = true
	}
	if started == false {
		c.printf(colInfo, "all methods failed\n")
		return errors.New("Unable to start the Docker engine")
	}
	c.printf(colInfo, "done\n")
	return nil
}

// Use the official script to install Docker
func (c *ampManager) installDockerEngine() error {
	resp, err := http.Get(dockerGetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	tmpfile, err := ioutil.TempFile("", "get-docker")
	if err != nil {
		c.fatalf("%v\n", err)
	}

	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(body); err != nil {
		c.fatalf("%v\n", err)
	}
	if err := tmpfile.Close(); err != nil {
		c.fatalf("%v\n", err)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("/bin/sh", tmpfile.Name())
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	c.printf(colInfo, "Installing Docker engine... \n")
	err = cmd.Run()
	if err != nil {
		c.printf(colInfo, "failed\n")
		c.printf(colInfo, out.String())
		c.printf(colInfo, stderr.String())
		return err
	}
	c.printf(colInfo, "done\n")
	return nil
}

func (c *ampManager) validateInstalledDockerEngineVersion(version string) (bool, error) {
	expected, err := semver.Make(minimumDockerEngineVersion)
	if err != nil {
		c.printf(colInfo, "Minimum version is %s\n", minimumDockerEngineVersion)
		return false, errors.New("Unable to parse minimum Docker engine version")
	}
	observed, err := semver.Make(version)
	if err != nil {
		c.printf(colInfo, "Version is %s\n", version)
		return false, errors.New("Unable to parse Docker engine version")
	}
	return observed.GTE(expected), nil
}

func (c *ampManager) isSwarmInit() bool {
	var out bytes.Buffer
	var stderr bytes.Buffer
	checkCmd := "docker node inspect self"
	checkCmdSplitted := strings.Split(checkCmd, " ")
	cmd := exec.Command(checkCmdSplitted[0], checkCmdSplitted[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			c.printf(colInfo, "Docker swarm mode is not set\n")
		} else {
			c.printf(colInfo, "Failed to check Swarm status: %s\n", err)
		}
		return false
	}
	return true
}

func (c *ampManager) joinSwarm() error {
	if c.docker == nil {
		defaultHeaders := map[string]string{"User-Agent": "amp-bootstrap"}
		dclient, err := dockerclient.NewClient(dockerURL, dockerAPIMinVersion, nil, defaultHeaders)
		if err != nil {
			return err
		}
		c.docker = dclient
	}
	if createSwarm == false {
		if swarmJoinToken == "" {
			return errors.New("worker join token not set")
		}
		remoteAddresses := []string{swarmManagerHost + swarmManagerPort}
		req := swarm.JoinRequest{ListenAddr: "0.0.0.0:" + swarmManagerPort, RemoteAddrs: remoteAddresses, JoinToken: swarmJoinToken}

		return c.docker.SwarmJoin(context.Background(), req)
	}
	if manager == false {
		return errors.New("won't initialize the Swarm on a worker node")
	}

	// if multiple IPs, one should be specified
	ifaces, err := net.Interfaces()
	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// @todo: decide if the IP can be used for the Swarm network (private or public)
			break
		}
	}
	if ip == nil {
		return errors.New("No interface available for Swarm init")
	}
	req := swarm.InitRequest{ListenAddr: "0.0.0.0:" + swarmManagerPort, AdvertiseAddr: ip.String()}
	c.printf(colInfo, "Swarm init on ip %s\n", ip.String())
	_, err = c.docker.SwarmInit(context.Background(), req)
	if err == nil {
		c.printf(colInfo, "done\n")
	}
	return err
}

func (c *ampManager) pullAmpImage(version string) error {
	image := "appcelerator/amp:" + version
	reader, err := c.docker.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	data := make([]byte, 1000, 1000)
	for {
		_, err := reader.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
	}
	c.printf(colInfo, "AMP image successfully pulled\n")
	return err
}

package config

import (
	"regexp"

	"github.com/appcelerator/amp/api/rpc/types"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/docker/client"
	"github.com/docker/docker/pkg/term"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MaxConfigSize is the maximum byte length of the `Config.Spec.Data` field.
const MaxConfigSize = 500 * 1024 // 500KB

type Server struct {
	Docker *docker.Docker
}

// TODO: these methods (and the ones in secret) need to actually pass `ctx` along

// CreateConfig creates and return a `CreateConfigResponse` with a `Config` based
// on the provided `CreateConfigRequest.ConfigSpec`.
// - Returns `InvalidArgument` if the `CreateConfigRequest.ConfigSpec` is malformed,
//   or if the secret data is too long or contains invalid characters.
// - Returns an error if the creation fails.
// From: api/control.proto
func (s *Server) CreateConfig(ctx context.Context, request *CreateConfigRequest) (*CreateConfigResponse, error) {
	if err := validateConfigSpec(request.Spec); err != nil {
		return nil, err
	}

	// TODO: fix vendored packages
	// all this because we can't use the same types right now due to vendor conflicts
	from := request.GetSpec()
	data := from.GetData()
	annotations := from.GetAnnotations()
	name := annotations.GetName()
	labels := annotations.GetLabels()

	stdin, stdout, stderr := term.StdStreams()
	cli := client.NewDockerCli(stdin, stdout, stderr)

	id, err := client.ConfigCreate(cli, name, labels, data)
	if err != nil {
		return nil, err
	}

	resp := &CreateConfigResponse{
		Config: &Config{
			Id:   id,
			Spec: from,
		},
	}
	return resp, nil
}

// ListConfigs returns a `ListConfigResponse` with a list all non-internal `Config`s being
// managed, or all secrets matching any name in `ListConfigsRequest.Names`, any
// name prefix in `ListConfigsRequest.NamePrefixes`, any id in
// `ListConfigsRequest.ConfigIDs`, or any id prefix in `ListConfigsRequest.IDPrefixes`.
// - Returns an error if listing fails.
// From: api/control.proto
func (s *Server) ListConfigs(ctx context.Context, request *ListConfigsRequest) (*ListConfigsResponse, error) {
	stdin, stdout, stderr := term.StdStreams()
	cli := client.NewDockerCli(stdin, stdout, stderr)

	secrets, err := client.ConfigList(cli)
	if err != nil {
		return nil, err
	}

	resp := &ListConfigsResponse{}
	for _, secret := range secrets {
		s := &Config{
			Spec: &ConfigSpec{
				Annotations: &types.Annotations{
					Name: secret,
				},
			},
		}
		resp.Configs = append(resp.Configs, s)
	}
	return resp, nil
}

// RemoveConfig removes the secret referenced by `RemoveConfigRequest.ID`.
// - Returns `InvalidArgument` if `RemoveConfigRequest.ID` is empty.
// - Returns `NotFound` if the a secret named `RemoveConfigRequest.ID` is not found.
// - Returns an error if the deletion fails.
// From: api/control.proto
func (s *Server) RemoveConfig(ctx context.Context, request *RemoveConfigRequest) (*RemoveConfigResponse, error) {
	stdin, stdout, stderr := term.StdStreams()
	cli := client.NewDockerCli(stdin, stdout, stderr)

	id := request.GetConfigId()

	if err := client.ConfigRemove(cli, id); err != nil {
		return nil, err
	}
	return &RemoveConfigResponse{ConfigId: id}, nil
}

func validateConfigSpec(spec *ConfigSpec) error {
	if spec == nil {
		return status.Errorf(codes.InvalidArgument, "invalid argument")
	}
	if err := validateConfigOrConfigAnnotations(spec.Annotations); err != nil {
		return err
	}

	if len(spec.Data) >= MaxConfigSize || len(spec.Data) < 1 {
		return status.Errorf(codes.InvalidArgument, "config data must be larger than 0 and less than %d bytes", MaxConfigSize)
	}
	return nil
}

func validateConfigOrConfigAnnotations(m *types.Annotations) error {
	if m.Name == "" {
		return status.Errorf(codes.InvalidArgument, "name must be provided")
	} else if len(m.Name) > 64 || !isValidConfigOrConfigName.MatchString(m.Name) {
		// if the name doesn't match the regex
		return status.Errorf(codes.InvalidArgument,
			"invalid name, only 64 [a-zA-Z0-9-_.] characters allowed, and the start and end character must be [a-zA-Z0-9]")
	}
	return nil
}

// configs and secrets have different naming requirements from tasks and services
var isValidConfigOrConfigName = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[a-zA-Z0-9-_.]*[a-zA-Z0-9])?$`)

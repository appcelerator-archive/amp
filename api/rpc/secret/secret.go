package secret

import (
	"fmt"
	"regexp"

	"github.com/appcelerator/amp/api/rpc/types"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/docker/client"
	"github.com/docker/docker/pkg/term"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MaxSecretSize is the maximum byte length of the `Secret.Spec.Data` field.
const MaxSecretSize = 500 * 1024 // 500KB

type Server struct {
	Docker *docker.Docker
}

// CreateSecret creates and return a `CreateSecretResponse` with a `Secret` based
// on the provided `CreateSecretRequest.SecretSpec`.
// - Returns `InvalidArgument` if the `CreateSecretRequest.SecretSpec` is malformed,
//   or if the secret data is too long or contains invalid characters.
// - Returns an error if the creation fails.
// From: api/control.proto
func (s *Server) CreateSecret(ctx context.Context, request *CreateSecretRequest) (*CreateSecretResponse, error) {
	fmt.Printf("CreateSecretRequest: %+v\n", request)
	if err := validateSecretSpec(request.Spec); err != nil {
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

	id, err := client.SecretCreate(cli, name, labels, data)
	if err != nil {
		return nil, err
	}

	resp := &CreateSecretResponse{
		Secret: &Secret{
			Id:   id,
			Spec: from,
		},
	}
	return resp, nil
}

// ListSecrets returns a `ListSecretResponse` with a list all non-internal `Secret`s being
// managed, or all secrets matching any name in `ListSecretsRequest.Names`, any
// name prefix in `ListSecretsRequest.NamePrefixes`, any id in
// `ListSecretsRequest.SecretIDs`, or any id prefix in `ListSecretsRequest.IDPrefixes`.
// - Returns an error if listing fails.
// From: api/control.proto
func (s *Server) ListSecrets(ctx context.Context, request *ListSecretsRequest) (*ListSecretsResponse, error) {
	fmt.Printf("ListSecrets: %+v\n", request)

	stdin, stdout, stderr := term.StdStreams()
	cli := client.NewDockerCli(stdin, stdout, stderr)

	secrets, err := client.SecretList(cli)
	if err != nil {
		return nil, err
	}

	resp := &ListSecretsResponse{}
	for _, secret := range secrets {
		s := &Secret{
			Spec: &SecretSpec{
				Annotations: &types.Annotations{
					Name: secret,
				},
			},
		}
		resp.Secrets = append(resp.Secrets, s)
	}
	return resp, nil
}

// RemoveSecret removes the secret referenced by `RemoveSecretRequest.ID`.
// - Returns `InvalidArgument` if `RemoveSecretRequest.ID` is empty.
// - Returns `NotFound` if the a secret named `RemoveSecretRequest.ID` is not found.
// - Returns an error if the deletion fails.
// From: api/control.proto
func (s *Server) RemoveSecret(ctx context.Context, request *RemoveSecretRequest) (*RemoveSecretResponse, error) {
	fmt.Printf("RemoveSecret: %+v\n", request)

	stdin, stdout, stderr := term.StdStreams()
	cli := client.NewDockerCli(stdin, stdout, stderr)

	id := request.GetSecretId()

	if err := client.SecretRemove(cli, id); err != nil {
		return nil, err
	}

	return &RemoveSecretResponse{SecretId: id}, nil
}

func validateSecretSpec(spec *SecretSpec) error {
	if spec == nil {
		return status.Errorf(codes.InvalidArgument, "invalid argument")
	}
	if err := validateConfigOrSecretAnnotations(spec.Annotations); err != nil {
		return err
	}

	if len(spec.Data) >= MaxSecretSize || len(spec.Data) < 1 {
		return status.Errorf(codes.InvalidArgument, "secret data must be larger than 0 and less than %d bytes", MaxSecretSize)
	}
	return nil
}

func validateConfigOrSecretAnnotations(m *types.Annotations) error {
	if m.Name == "" {
		return status.Errorf(codes.InvalidArgument, "name must be provided")
	} else if len(m.Name) > 64 || !isValidConfigOrSecretName.MatchString(m.Name) {
		// if the name doesn't match the regex
		return status.Errorf(codes.InvalidArgument,
			"invalid name, only 64 [a-zA-Z0-9-_.] characters allowed, and the start and end character must be [a-zA-Z0-9]")
	}
	return nil
}

// configs and secrets have different naming requirements from tasks and services
var isValidConfigOrSecretName = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[a-zA-Z0-9-_.]*[a-zA-Z0-9])?$`)

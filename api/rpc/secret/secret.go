package secret

import (
	"fmt"
	"regexp"

	"github.com/appcelerator/amp/api/rpc/types"
	"github.com/appcelerator/amp/pkg/docker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MaxSecretSize is the maximum byte length of the `Secret.Spec.Data` field.
const MaxSecretSize = 500 * 1024 // 500KB

type Server struct {
	Docker   *docker.Docker
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

	resp := &CreateSecretResponse{
		Secret: &Secret{
			Id:   "foo",
			Spec: request.GetSpec(),
		},
	}
	return resp, nil
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


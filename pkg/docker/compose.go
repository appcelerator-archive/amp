package docker

import (
	"fmt"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	"github.com/appcelerator/amp/docker/cli/cli/compose/loader"
	composetypes "github.com/appcelerator/amp/docker/cli/cli/compose/types"
	"golang.org/x/net/context"
)

// ComposeParse parses a compose file
func (d *Docker) ComposeParse(ctx context.Context, composeFile []byte) (*composetypes.Config, error) {
	var details composetypes.ConfigDetails
	var err error

	// WorkingDir
	details.WorkingDir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Parsing compose file
	config, err := loader.ParseYAML(composeFile)
	if err != nil {
		return nil, err
	}
	details.ConfigFiles = []composetypes.ConfigFile{{
		Filename: "filename",
		Config:   config,
	}}

	// Environment
	env := os.Environ()
	details.Environment = make(map[string]string, len(env))
	for _, s := range env {
		if !strings.Contains(s, "=") {
			return nil, fmt.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		details.Environment[kv[0]] = kv[1]
	}

	return loader.Load(details)
}

// ComposeIsAuthorized checks if the given compose file is authorized
func (d *Docker) ComposeIsAuthorized(compose *composetypes.Config) bool {
	for _, reservedSecret := range constants.Secrets {
		if _, exists := compose.Secrets[reservedSecret]; exists {
			return false
		}
	}

	for _, service := range compose.Services {
		for _, reservedSecret := range constants.Secrets {
			for _, secret := range service.Secrets {
				if strings.EqualFold(secret.Source, reservedSecret) {
					return false
				}
			}
		}
		for _, reservedLabel := range constants.Labels {
			if _, exists := service.Labels[reservedLabel]; exists {
				return false
			}
			if _, exists := service.Deploy.Labels[reservedLabel]; exists {
				return false
			}
		}
	}
	return true
}

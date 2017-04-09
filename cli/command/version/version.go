package version

import (
	"bytes"
	"runtime"

	"fmt"

	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cli"
	"github.com/docker/docker/pkg/templates"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Version struct {
	Client *ClientVersionInfo
	Server *version.Info
}

// AmplifierOK Checks if AMP is connected to Amplifier
func (v Version) IsConnected() bool {
	return v.Server != nil
}

type ClientVersionInfo struct {
	Version   string
	Build     string
	Address   string
	GoVersion string
	Os        string
	Arch      string
}

var versionTemplate = `amp:
 Version:       {{.Client.Version}}
 Build:         {{.Client.Build}}
 Address:       {{.Client.Address}}
 Go version:    {{.Client.GoVersion}}
 OS/Arch:       {{.Client.Os}}/{{.Client.Arch}}

amplifier:      {{if .IsConnected}}
 Version:       {{.Server.Version}}
 Build:         {{.Server.Build}}
 Go version:    {{.Server.GoVersion}}
 OS/Arch:       {{.Server.Os}}/{{.Server.Arch}}{{else}}not connected{{end}}`

// NewVersionCommand returns a new instance of the version command.
func NewVersionCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Show amp version information",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showVersion(c)
		},
	}
}

// Print version info of client and server (if connected).
func showVersion(c cli.Interface) error {
	tmpl, err := templates.Parse(versionTemplate)
	if err != nil {
		return fmt.Errorf("template parsing error: %v\n", err)
	}

	v := Version{
		Client: &ClientVersionInfo{
			Version:   c.Version(),
			Build:     c.Build(),
			Address:   c.Address(),
			GoVersion: runtime.Version(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
	}

	conn, err := c.ClientConn()
	if err == nil {
		client := version.NewVersionClient(conn)
		reply, err := client.Get(context.Background(), &version.GetRequest{})
		if err != nil {
			return fmt.Errorf("%s", grpc.ErrorDesc(err))
		}
		v.Server = &version.Info{
			Version:   reply.Info.Version,
			Build:     reply.Info.Build,
			GoVersion: reply.Info.GoVersion,
			Os:        reply.Info.Os,
			Arch:      reply.Info.Arch,
		}
	}

	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, v); err != nil {
		return fmt.Errorf("executing templating error: %v\n", err)
	}
	c.Console().Println(doc.String())
	return err
}

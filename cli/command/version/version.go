package version

import (
	"bytes"
	"fmt"
	"runtime"

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
	Server    string
	GoVersion string
	OS        string
	Arch      string
}

var template = `Client:
 Version:       {{.Client.Version}}
 Build:         {{.Client.Build}}
 Server:        {{.Client.Server}}
 Go version:    {{.Client.GoVersion}}
 OS/Arch:       {{.Client.OS}}/{{.Client.Arch}}

Server:         {{if .IsConnected}}
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
	tmpl, err := templates.Parse(template)
	if err != nil {
		return fmt.Errorf("template parsing error: %v\n", err)
	}

	v := Version{
		Client: &ClientVersionInfo{
			Version:   c.Version(),
			Build:     c.Build(),
			Server:    c.Server(),
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
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

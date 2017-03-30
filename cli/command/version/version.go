package version

import (
	"bytes"
	"runtime"

	"github.com/appcelerator/amp/cli"
	"github.com/docker/docker/pkg/templates"
	"github.com/spf13/cobra"
)

type Version struct {
	Client *ClientVersionInfo
	Server interface{} // should be:  *ServerVersionInfo from rpc/version/version.proto ...
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
 Address:       {{.Server.Address}}
 Go version:    {{.Server.GoVersion}}
 OS/Arch:       {{.Server.Os}}/{{.Server.Arch}}{{else}}not connected{{end}}`

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

	templateFormat := versionTemplate
	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		c.Console().Fatalf("template parsing error: %v\n", err)
	}
	var doc bytes.Buffer

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

	// TODO: rpc service needs to be refactored ... it should return a ServerVersionInfo as protobuf-generated struct
	// the service shouldn't know anything about the client

	//request := &version.ListRequest{}
	//if err = AMP.Connect(); err == nil {
	//	client := version.NewVersionClient(c.ClientConn())
	//	reply, err := client.List(context.Background(), request)
	//	if err != nil {
	//		mgr.Fatal(grpc.ErrorDesc(err))
	//	}
	//	vd.Amplifier = &version.Details{
	//		Version:   reply.Reply.Version,
	//		Port:      reply.Reply.Port,
	//		GoVersion: reply.Reply.Goversion,
	//		Os:        reply.Reply.Os,
	//		Arch:      reply.Reply.Arch,
	//	}
	//}

	if err := tmpl.Execute(&doc, v); err != nil {
		c.Console().Fatalf("executing templating error: %v\n", err)
	}

	c.Console().Println(doc.String())
	return err
}

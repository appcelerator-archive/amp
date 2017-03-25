package main

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/appcelerator/amp/api/rpc/version"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/docker/docker/pkg/templates"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var versionTemplate = `AMP:
 Version:       {{.AMP.Version}}
 Build:         {{.AMP.Build}}
 ConfigAddr:    {{.AMP.ConfigAddr}}
 Go version:    {{.AMP.GoVersion}}
 OS/Arch:       {{.AMP.Os}}/{{.AMP.Arch}}{{if .AmplifierOK}}

Amplifier:
 Version:       {{.Amplifier.Version}}
 Port:          {{.Amplifier.Port}}
 Go version:    {{.AMP.GoVersion}}
 OS/Arch:       {{.AMP.Os}}/{{.AMP.Arch}}{{end}}`

// VersionCmd represents the amp version
var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display the version info for AMP and Amplifier",
	Example: " ",
	RunE: func(cmd *cobra.Command, args []string) error {
		return list(AMP)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

// Lists version info of AMP and Amplifier
func list(amp *cli.AMP) error {

	templateFormat := versionTemplate
	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		mgr.Fatal("template parsing error: %v", err)
	}
	var doc bytes.Buffer

	vd := version.Config{
		AMP: &version.Details{
			Version:    Version,
			Build:      Build,
			ConfigAddr: amp.Configuration.AmpAddress,
			GoVersion:  runtime.Version(),
			Os:         runtime.GOOS,
			Arch:       runtime.GOARCH,
		},
	}

	request := &version.ListRequest{}
	if err = AMP.Connect(); err == nil {
		client := version.NewVersionClient(amp.Conn)
		reply, err := client.List(context.Background(), request)
		if err != nil {
			mgr.Fatal(grpc.ErrorDesc(err))
		}
		vd.Amplifier = &version.Details{
			Version:   reply.Reply.Version,
			Port:      reply.Reply.Port,
			GoVersion: reply.Reply.Goversion,
			Os:        reply.Reply.Os,
			Arch:      reply.Reply.Arch,
		}
	}

	if err := tmpl.Execute(&doc, vd); err != nil {
		mgr.Fatal("executing templating error: %v", err)
	}

	fmt.Println(doc.String())
	return err
}

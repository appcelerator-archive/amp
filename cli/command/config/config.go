package config

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AMPConfig struct {
	Config *cli.Configuration
}

var t = `AMP Configuration:
 Server:        {{.Config.Server}}`

// NewConfigCommand returns a new instance of the config command.
func NewConfigCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Display amp configuration",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig(c)
		},
	}
}

func showConfig(c cli.Interface) error {
	conf := &AMPConfig{
		&cli.Configuration{
			Server: c.Server(),
		},
	}

	tmpl, err := template.New("t").Parse(t)
	if err != nil {
		return fmt.Errorf("template parsing error: %v\n", err)
	}

	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, conf); err != nil {
		return fmt.Errorf("template execution error: %v\n", err)
	}
	if viper.ConfigFileUsed() != "" {
		c.Console().Println("Configuration file:", viper.ConfigFileUsed())
	} else {
		c.Console().Println("Configuration file: none")
	}
	c.Console().Println(doc.String())
	return nil
}

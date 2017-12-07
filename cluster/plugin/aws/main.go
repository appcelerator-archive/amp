package main

import (
	"fmt"
	"os"

	"github.com/appcelerator/amp/cluster/plugin/aws/cmd"
	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
)

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s - Build: %s\n", Version, Build)
}

func main() {
	rootCmd := cmd.NewRootCommand()
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "version of the plugin",
		Run:   version,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/appcelerator/amp/cluster/envoy"
)

func foo(cmd *cobra.Command, args []string) {
	log.Println(envoy.Foo())
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "envoy",
		Short: "run commands in target cluster",
		// If needed
		// PersistentPreRun: initEnvoy,
	}

	fooCmd := &cobra.Command{
		Use:   "foo",
		Short: "foo bar",
		Run:   foo,
	}

	rootCmd.AddCommand(fooCmd)

	_ = rootCmd.Execute()
}

package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ampctl",
		Short: "Run commands in target amp cluster",
	}

	rootCmd.AddCommand(NewChecksCommand())
	rootCmd.AddCommand(NewInstallCommand())
	rootCmd.AddCommand(NewMonitorCommand())

	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

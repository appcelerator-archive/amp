package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/appcelerator/amp/api/client"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an image in the amp registry",
	Long:  `Push an image to the amp registry`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RegistryPush(AMP, cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
}

// RegistryPush displays resource usage statistcs
func RegistryPush(amp *client.AMP, cmd *cobra.Command, args []string) error {
	_, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}

	image := args[0]

	if amp.Verbose() {
		fmt.Println("Execute registry push command with:")
		fmt.Printf("image: %s\n", image)
	}

	if err = validateRegistryImage(image); err != nil {
		return err
	}
	cmdexe := exec.Command("docker", "push", image)
	cmdexe.Stdout = os.Stdout
	cmdexe.Stderr = os.Stderr
	err = cmdexe.Run()
	if err != nil {
		return err
	}
	return err
}

func validateRegistryImage(image string) error {
	if image == "" {
		return errors.New("Need a valid image name")
	}
	return nil
}

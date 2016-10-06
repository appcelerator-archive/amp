package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/appcelerator/amp/api/client"
	"github.com/spf13/cobra"
)

// Registry is the main command for attaching registry subcommands.
var RegCmd = &cobra.Command{
	Use:   "registry operations",
	Short: "Registry operations",
	Long:  `Manage registry-related operations.`,
}

var (
	domain = "local.appcelerator.io"
	pushCmd = &cobra.Command{
		Use:   "push [image]",
		Short: "Push an image to the amp registry",
		Long:  `Push an image to the amp registry`,
		Run: func(cmd *cobra.Command, args []string) {
			err := RegistryPush(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	reglsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List the amp registry images",
		Long:  `List the amp registry images`,
		Run: func(cmd *cobra.Command, args []string) {
			err := RegistryLs(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}	
)



func init() {
	RootCmd.AddCommand(RegCmd)
	pushCmd.Flags().StringVar(&domain, "domain", domain, "The amp domain")
	RegCmd.AddCommand(pushCmd)
	RegCmd.AddCommand(reglsCmd)
}

// RegistryPush displays resource usage statistics
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
	taggedImage := image
	if !strings.HasPrefix(image, "registry."+domain) {
		nn := strings.Index(image, "/")
		if (nn < 0) {
			return fmt.Errorf("Invalid image name %s", image)
		}
		taggedImage = "registry."+domain+"/"+image[nn+1:]
		fmt.Printf("Tag image from %s to %s\n", image, taggedImage)
		cmdexe := exec.Command("docker", "tag", image, taggedImage)
		cmdexe.Stdout = os.Stdout
		cmdexe.Stderr = os.Stderr
		err = cmdexe.Run()
		if err != nil {
			return err
		}
	}
	fmt.Printf("push image %s\n", taggedImage)
	cmdexe := exec.Command("docker", "push", taggedImage)
	cmdexe.Stdout = os.Stdout
	cmdexe.Stderr = os.Stderr
	err = cmdexe.Run()
	if err != nil {
		return err
	}
	return err
}

// RegistryPush displays resource usage statistics: localhost:5000/v2/_catalog
func RegistryLs(amp *client.AMP, cmd *cobra.Command, args []string) error {
	_, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}
	resp, err := http.Get("http://registry."+domain+"/v2/_catalog")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return err
}

func validateRegistryImage(image string) error {
	if image == "" {
		return errors.New("Need a valid image name")
	}
	return nil
}

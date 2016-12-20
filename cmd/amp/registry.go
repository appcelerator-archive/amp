package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appcelerator/amp/api/client"
	distreference "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/reference"
	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"regexp"
)

// RegCmd is the main command for attaching registry subcommands.
var RegCmd = &cobra.Command{
	Use:   "registry",
	Short: "Registry operations",
	Long:  `Manage registry-related operations.`,
}

var (
	endpoint string
	domain   string
	insecure bool
	pushCmd  = &cobra.Command{
		Use:   "push [image]",
		Short: "Push an image to the amp registry",
		Long:  `Push an image to the amp registry.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RegistryPush(AMP, cmd, args)
		},
	}
	reglsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List the amp registry images",
		Long:  `List the amp registry images.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RegistryLs(AMP, cmd, args)
		},
	}
)

func init() {
	RootCmd.AddCommand(RegCmd)
	RegCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "i", true, "Insecure registry")
	RegCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "local.appcelerator.io", "The amp registry domain (hostname or IP)")
	RegCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "The amp registry endpoint (hostname or IP), overrides the domain option")
	RegCmd.AddCommand(pushCmd)
	RegCmd.AddCommand(reglsCmd)
}

// registryEndpoint returns the registry endpoint
func registryEndpoint() (ep string) {
	if endpoint != "" {
		ep = endpoint
		return
	}
	ep = "registry." + domain
	return
}

// RegistryPush displays resource usage statistics
func RegistryPush(amp *client.AMP, cmd *cobra.Command, args []string) error {
	defaultHeaders := map[string]string{"User-Agent": "amp-cli"}
	dclient, err := docker.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = amp.GetAuthorizedContext()
	if err != nil {
		return err
	}
	// @todo: read the .dockercfg file for authentication, or use credentials from amp.yaml
	ac := types.AuthConfig{Username: "none"}
	jsonString, err := json.Marshal(ac)
	if err != nil {
		return errors.New("failed to marshal authconfig")
	}
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(jsonString)))
	base64.URLEncoding.Encode(dst, jsonString)
	authConfig := string(dst)
	imagePushOptions := types.ImagePushOptions{RegistryAuth: authConfig}

	image := args[0]
	distributionRef, err := distreference.ParseNamed(image)
	if err != nil {
		return fmt.Errorf("error parsing reference: %q is not a valid repository/tag", image)
	}
	if _, isCanonical := distributionRef.(distreference.Canonical); isCanonical {
		return errors.New("refusing to create a tag with a digest reference")
	}
	tag := reference.GetTagFromNamedRef(distributionRef)
	hostname, name := distreference.SplitHostname(distributionRef)

	if amp.Verbose() {
		fmt.Printf("Registry push request with:\n  image: %s\n", image)
	}

	taggedImage := image
	if hostname != registryEndpoint() {
		taggedImage = registryEndpoint() + "/" + name + ":" + tag
		fmt.Printf("Tag image from %s to %s\n", image, taggedImage)
		if err := dclient.ImageTag(ctx, image, taggedImage); err != nil {
			return err
		}
	}
	fmt.Printf("Push image %s\n", taggedImage)
	resp, err := dclient.ImagePush(ctx, taggedImage, imagePushOptions)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`: digest: sha256:`)
	if !re.Match(body) {
		fmt.Print(string(body))
		return errors.New("push failed")
	}
	return nil
}

// RegistryLs lists images
func RegistryLs(amp *client.AMP, cmd *cobra.Command, args []string) error {
	_, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}
	var protocol string
	if insecure {
		protocol = "http"
	} else {
		protocol = "https"
	}
	resp, err := http.Get(protocol + "://" + registryEndpoint() + "/v2/_catalog")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return err
}

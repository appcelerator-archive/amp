package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/appcelerator/amp/cmd/amp/cli"
	distreference "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/reference"
	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"regexp"
)

// RegCmd is the main command for attaching registry sub-commands.
var RegCmd = &cobra.Command{
	Use:   "registry",
	Short: "Registry operations",
}

var (
	endpoint string
	domain   string
	insecure bool
	pushCmd  = &cobra.Command{
		Use:     "push",
		Short:   "Push an image to the AMP registry",
		Example: "amp registry push sample/test-registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RegistryPush(AMP, args)
		},
	}
	reglsCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List the AMP registry images",
		Example: "amp registry ls \namp registry ls -q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RegistryLs(AMP)
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

	reglsCmd.Flags().BoolP("quiet", "q", false, "Only display IDs")
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
func RegistryPush(amp *cli.AMP, args []string) error {
	defaultHeaders := map[string]string{"User-Agent": "amp-cli"}
	dclient, err := docker.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	ctx := context.Background()
	_, err = amp.GetAuthorizedContext()
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	// @todo: read the .dockercfg file for authentication, or use credentials from amp.yaml
	ac := types.AuthConfig{Username: "none"}
	jsonString, err := json.Marshal(ac)
	if err != nil {
		mgr.Error("failed to marshal authconfig")
	}
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(jsonString)))
	base64.URLEncoding.Encode(dst, jsonString)
	authConfig := string(dst)
	imagePushOptions := types.ImagePushOptions{RegistryAuth: authConfig}

	image := args[0]
	distributionRef, err := distreference.ParseNamed(image)
	if err != nil {
		mgr.Error("error parsing reference: %q is not a valid repository/tag", image)
	}
	if _, isCanonical := distributionRef.(distreference.Canonical); isCanonical {
		mgr.Error("refusing to create a tag with a digest reference")
	}
	tag := reference.GetTagFromNamedRef(distributionRef)
	hostname, name := distreference.SplitHostname(distributionRef)

	if amp.Verbose() {
		mgr.Info("Registry push request with:\n  image: %s\n", image)
	}

	taggedImage := image
	if hostname != registryEndpoint() {
		taggedImage = registryEndpoint() + "/" + name + ":" + tag
		mgr.Regular("Tag image from %s to %s\n", image, taggedImage)
		if err := dclient.ImageTag(ctx, image, taggedImage); err != nil {
			mgr.Error(grpc.ErrorDesc(err))
		}
	}
	fmt.Printf("Push image %s\n", taggedImage)
	resp, err := dclient.ImagePush(ctx, taggedImage, imagePushOptions)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	re := regexp.MustCompile(`: digest: sha256:`)
	if !re.Match(body) {
		fmt.Print(string(body))
		//return errors.New("push failed")
		mgr.Error("push failed")
	}
	return nil
}

// RegistryLs lists images
func RegistryLs(amp *cli.AMP) error {
	_, err := amp.GetAuthorizedContext()
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	var protocol string
	if insecure {
		protocol = "http"
	} else {
		protocol = "https"
	}
	resp, err := http.Get(protocol + "://" + registryEndpoint() + "/v2/_catalog")
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return err
}

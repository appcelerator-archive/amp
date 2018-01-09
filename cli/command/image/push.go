package image

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/api/rpc/image"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

const (
	defaultURL     = "unix:///var/run/docker.sock"
	defaultVersion = "1.30"
)

var dockerClient *docker.Client

// NewPushCommand returns a new instance of the push command for pushing an image
func NewPushCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "push SRC_IMAGE [DST_IMAGE]",
		Short:   "Push an image to the cluster registry",
		PreRunE: cli.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return push(c, cmd, args)
		},
	}
}

func tempImageName() string {
	prefix := "tmp/amp"
	repo := make([]byte, 16)
	tag := "local"
	rand.Read(repo)
	return fmt.Sprintf("%s%s:%s", prefix, hex.EncodeToString(repo), tag)
}

func push(c cli.Interface, cmd *cobra.Command, args []string) error {
	dockerClient, err := docker.NewClient(defaultURL, defaultVersion, nil, nil)
	if err != nil {
		return err
	}
	source := args[0]
	var dest string
	if len(args) == 1 {
		dest = source
	} else {
		dest = args[1]
	}
	// tagging the image with a random name before saving it (the name will be loaded too)
	tmpImage := tempImageName()
	if err = dockerClient.ImageTag(context.Background(), source, tmpImage); err != nil {
		return fmt.Errorf("error tagging image [%s] to [%s]", source, tmpImage)
	}
	defer dockerClient.ImageRemove(context.Background(), tmpImage, types.ImageRemoveOptions{})
	var data []byte
	saveResp, err := dockerClient.ImageSave(context.Background(), []string{tmpImage})
	if err != nil {
		return fmt.Errorf("error saving image [%s]: %s", tmpImage, err.Error())
	}
	defer saveResp.Close()
	data, err = ioutil.ReadAll(saveResp)
	if err != nil {
		return fmt.Errorf("error reading docker save output for image [%s]: %s", tmpImage, err.Error())
	}

	conn := c.ClientConn()
	client := image.NewImageClient(conn)
	request := &image.PushRequest{
		Name: dest,
		Data: data,
	}
	reply, err := client.ImagePush(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New("ImagePush API call failed: " + s.Message())
		}
		return fmt.Errorf("error pushing image [%s]: %s", dest, err)
	}
	fmt.Println(reply.GetDigest())

	return nil
}

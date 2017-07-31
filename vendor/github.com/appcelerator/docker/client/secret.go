package client

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type createOptions struct {
	name   string
	file   string
	labels opts.ListOpts
}

// TODO: this is intended to be used as the client from the amplifier API, so
// this needs to be modified to take the file contents sent from the amp cli,
// not read the file.
func SecretCreate(dockerCli command.Cli, options createOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	var in io.Reader = dockerCli.In()
	if options.file != "-" {
		file, err := system.OpenSequential(options.file)
		if err != nil {
			return err
		}
		in = file
		defer file.Close()
	}

	secretData, err := ioutil.ReadAll(in)
	if err != nil {
		return errors.Errorf("Error reading content from %q: %v", options.file, err)
	}

	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name:   options.name,
			Labels: opts.ConvertKVStringsToMap(options.labels.GetAll()),
		},
		Data: secretData,
	}

	r, err := client.SecretCreate(ctx, spec)
	if err != nil {
		return err
	}

	fmt.Fprintln(dockerCli.Out(), r.ID)
	return nil
}

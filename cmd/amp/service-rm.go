package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	// TODO: add support for 'remove' alias
	serviceRmCmd = &cobra.Command{
		Use:   "rm [OPTIONS] SERVICE [SERVICE...]",
		Short: "Remove one or more services",
		Long:  `Remove one or more services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return serviceRm(AMP, cmd, args)
		},
	}

	// services to remove
	services []string
)

func init() {
	ServiceCmd.AddCommand(serviceRmCmd)
}

func serviceRm(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// TODO use standard errors and print usage
		log.Fatal("\"amp service rm\" requires at least 1 argument(s)")
	}

	client := service.NewServiceClient(amp.Conn)

	for _, ident := range args {
		req := &service.RemoveRequest{
			Ident: ident,
		}

		resp, err := client.Remove(context.Background(), req)
		if err != nil {
			pattern := "daemon: "
			errstr := fmt.Sprintf("%s", err)
			index := strings.LastIndex(errstr, pattern)
			if index > -1 {
				errmsg := fmt.Sprintf("Error: %s", errstr[index+len(pattern):])
				fmt.Println(errmsg)
			} else {
				fmt.Printf("Error: %s\n", err)
			}
		} else {
			fmt.Println(resp.Ident)
		}
	}

	return nil
}

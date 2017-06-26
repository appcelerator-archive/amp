package org

import (
	"errors"
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewOrgRemoveCommand returns a new instance of the remove organization command.
func NewOrgRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm ORGANIZATION(S)",
		Short:   "Remove organization",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrg(c, args)
		},
	}
}

func removeOrg(c cli.Interface, args []string) error {
	var errs []string
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)

	for _, org := range args {
		requestDelete := &account.DeleteOrganizationRequest{
			Name: org,
		}
		if _, err := client.DeleteOrganization(context.Background(), requestDelete); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(org)
		// Check if user is logged in on behalf of an org
		token, err := cli.ReadToken(c.Server())
		if err != nil {
			return err
		}
		pToken, _ := jwt.ParseWithClaims(token, &auth.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte{}, nil
		})
		if claims, ok := pToken.Claims.(*auth.AuthClaims); ok {
			if claims.ActiveOrganization == org {
				// Switch back to the users account
				client := account.NewAccountClient(conn)
				requestSwitch := &account.SwitchRequest{
					Account: claims.AccountName,
				}
				headers := metadata.MD{}
				_, err := client.Switch(context.Background(), requestSwitch, grpc.Header(&headers))
				if err != nil {
					return fmt.Errorf("%s", grpc.ErrorDesc(err))
				}
				if err := cli.SaveToken(headers, c.Server()); err != nil {
					return fmt.Errorf("%s", grpc.ErrorDesc(err))
				}
				c.Console().Printf("Switched back to account: %s\n", claims.AccountName)
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	if err := cli.RemoveFile(c.Server()); err != nil {
		return err
	}
	return nil
}

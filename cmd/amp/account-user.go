package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
)

// UserCmd is the main command for attaching user sub-commands.
var (
	listUserCmd = &cobra.Command{
		Use:   "ls",
		Short: "List user",
		Long:  `The list command lists all available users.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(AMP)
		},
	}

	deleteUserCmd = &cobra.Command{
		Use:     "del",
		Short:   "Delete user",
		Long:    `The delete command deletes a user.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteUser(AMP, cmd)
		},
	}

	getUserCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user info",
		Long:  `The get command retrieves details of a user.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getUser(AMP, cmd)
		},
	}
)

func init() {
	UserCmd.AddCommand(listUserCmd)
	UserCmd.AddCommand(deleteUserCmd)
	UserCmd.AddCommand(getUserCmd)

	getUserCmd.Flags().StringVar(&name, "name", name, "Account Name")
	deleteUserCmd.Flags().StringVar(&name, "name", name, "Account Name")
}

// listUser validates the input command line arguments and lists all users
// by invoking the corresponding rpc/storage method
func listUser(amp *cli.AMP) (err error) {
	request := &account.ListUsersRequest{}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.ListUsers(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\t")
	for _, user := range reply.Users {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Email)
	}
	w.Flush()
	return nil
}

// deleteUser validates the input command line arguments and deletes a user
// by invoking the corresponding rpc/storage method
func deleteUser(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		name = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		name = GetName()
	}

	request := &account.DeleteUserRequest{
		Name: name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.DeleteUser(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully deleted user.")
	return nil
}

// getUser validates the input command line arguments and retrieves info of a user
// by invoking the corresponding rpc/storage method
func getUser(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		name = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		name = GetName()
	}

	request := &account.GetUserRequest{
		Name: name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.GetUser(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\tVERIFIED\tCREATED\t")
	fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", reply.User.Name, reply.User.Email, reply.User.IsVerified, ConvertTime(reply.User.CreateDt))
	w.Flush()
	return nil
}

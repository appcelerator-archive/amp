package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// OrgCmd is the main command for attaching organization sub-commands.
var (
	listOrgCmd = &cobra.Command{
		Use:     "list",
		Short:   "List organization",
		Long:    `The list command lists all available organizations.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrg(AMP)
		},
	}

	createOrgCmd = &cobra.Command{
		Use:   "create",
		Short: "Create organization",
		Long:  `The create command creates an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrg(AMP)
		},
	}

	deleteOrgCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete organization",
		Long:    `The delete command deletes an organization.`,
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteOrg(AMP)
		},
	}

	getOrgCmd = &cobra.Command{
		Use:   "get",
		Short: "Get organization info",
		Long:  `The get command retrieves details of an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getOrg(AMP)
		},
	}

	addOrgCmd = &cobra.Command{
		Use:   "add",
		Short: "Add members to organization",
		Long:  `The add command adds members to an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addMem(AMP)
		},
	}

	removeOrgCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Remove members from organization",
		Long:    `The remove command removes from an organization.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeMem(AMP)
		},
	}
)

func init() {
	OrgCmd.AddCommand(listOrgCmd)
	OrgCmd.AddCommand(createOrgCmd)
	OrgCmd.AddCommand(deleteOrgCmd)
	OrgCmd.AddCommand(getOrgCmd)
	OrgCmd.AddCommand(addOrgCmd)
	OrgCmd.AddCommand(removeOrgCmd)
}

// listOrg validates the input command line arguments and lists available organizations
// by invoking the corresponding rpc/storage method
func listOrg(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will list organizations.")
	request := &account.ListOrganizationsRequest{}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.ListOrganizations(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	fmt.Fprintln(w, "------------\t-----\t-------\t")
	for _, org := range reply.Organizations {
		orgCreate, err := strconv.ParseInt(strconv.FormatInt(org.CreateDt, 10), 10, 64)
		if err != nil {
			panic(err)
		}
		orgCreateTime := time.Unix(orgCreate, 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", org.Name, org.Email, orgCreateTime)
	}
	w.Flush()
	return nil
}

// createOrg validates the input command line arguments and creates an organization
// by invoking the corresponding rpc/storage method
func createOrg(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will create an organization.")
	orgName := getOrgName()
	email := getEmailAddress()
	request := &account.CreateOrganizationRequest{
		Name:  orgName,
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.CreateOrganization(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Hi %s! Please check your email to complete the signup process.", orgName)
	return nil
}

// deleteOrg validates the input command line arguments and deletes an organization
// by invoking the corresponding rpc/storage method
func deleteOrg(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will delete an organization.")
	orgName := getOrgName()
	request := &account.DeleteOrganizationRequest{
		Name: orgName,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.DeleteOrganization(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "The organization has been deleted successfully.")
	return nil
}

// getOrg validates the input command line arguments and retrieves info of an organization
// by invoking the corresponding rpc/storage method
func getOrg(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will get details of an organization.")
	orgName := getOrgName()
	request := &account.GetOrganizationRequest{
		Name: orgName,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.GetOrganization(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	fmt.Fprintln(w, "------------\t-----\t-------\t")
	orgCreate, err := strconv.ParseInt(strconv.FormatInt(reply.Organization.CreateDt, 10), 10, 64)
	if err != nil {
		panic(err)
	}
	orgCreateTime := time.Unix(orgCreate, 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t", reply.Organization.Name, reply.Organization.Email, orgCreateTime)

	fmt.Fprintln(w, "MEMBER NAME\tROLE\t")
	fmt.Fprintln(w, "-----------\t----\t")
	for _, mem := range reply.Organization.Members {
		fmt.Fprintf(w, "%s\t%s\t\n", mem.Name, mem.Role)
	}
	fmt.Fprintln(w, "TEAM NAME\tCREATED\t")
	fmt.Fprintln(w, "---------\t-------\t")
	for _, team := range reply.Organization.Teams {
		teamCreate, err := strconv.ParseInt(strconv.FormatInt(team.CreateDt, 10), 10, 64)
		if err != nil {
			panic(err)
		}
		teamCreateTime := time.Unix(teamCreate, 0)
		fmt.Fprintf(w, "%s\t%s\t", team.Name, teamCreateTime)
	}
	w.Flush()
	return nil
}

// addMem validates the input command line arguments and adds members to an organization
// by invoking the corresponding rpc/storage method
func addMem(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will add members to an organization.")
	orgName := getOrgName()
	name := getUserName()
	request := &account.AddUserToOrganizationRequest{
		OrganizationName: orgName,
		UserName:         name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.AddUserToOrganization(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Member(s) have been added to organization successfully.")
	return nil
}

// removeMem validates the input command line arguments and removes members from an organization
// by invoking the corresponding rpc/storage method
func removeMem(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will remove members from an organization.")
	orgName := getOrgName()
	name := getUserName()
	request := &account.RemoveUserFromOrganizationRequest{
		OrganizationName: orgName,
		UserName:         name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.RemoveUserFromOrganization(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Member(s) have been removed from organization successfully.")
	return nil
}

func getOrgName() (org string) {
	fmt.Print("organization name: ")
	fmt.Scanln(&org)
	org = strings.TrimSpace(org)
	err := schema.CheckName(org)
	if err != nil {
		manager.printf(colWarn, err.Error())
		return getOrgName()
	}
	return
}

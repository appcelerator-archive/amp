package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

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
		Use:   "ls",
		Short: "List organization",
		Long:  `The list command lists all available organizations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrg(AMP, cmd)
		},
	}

	createOrgCmd = &cobra.Command{
		Use:   "create",
		Short: "Create organization",
		Long:  `The create command creates an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrg(AMP, cmd)
		},
	}

	deleteOrgCmd = &cobra.Command{
		Use:     "del",
		Short:   "Delete organization",
		Long:    `The delete command deletes an organization.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteOrg(AMP, cmd)
		},
	}

	getOrgCmd = &cobra.Command{
		Use:   "get",
		Short: "Get organization info",
		Long:  `The get command retrieves details of an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getOrg(AMP, cmd)
		},
	}

	memOrgCmd = &cobra.Command{
		Use:   "member",
		Short: "Member-related operations in an organization",
		Long:  `The member command manages member-related operations for an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return memberOrg()
		},
	}

	addOrgMemCmd = &cobra.Command{
		Use:   "add",
		Short: "Add members to organization",
		Long:  `The add command adds members to an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMem(AMP, cmd)
		},
	}

	changeOrgMemRoleCmd = &cobra.Command{
		Use:   "role owner|member",
		Short: "Change role of organization member",
		Long:  `The role command changes the role of an organization member.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return changeOrgMem(AMP, cmd)
		},
	}

	remOrgMemCmd = &cobra.Command{
		Use:     "del",
		Short:   "Remove members from organization",
		Long:    `The remove command removes member from an organization.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return remOrgMem(AMP, cmd)
		},
	}

	listOrgMemCmd = &cobra.Command{
		Use:   "ls",
		Short: "List members of organization",
		Long:  `The list command lists members of an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrgMem(AMP, cmd)
		},
	}

	organization string
	member       string
	role         string
)

func init() {
	OrgCmd.AddCommand(listOrgCmd)
	OrgCmd.AddCommand(createOrgCmd)
	OrgCmd.AddCommand(deleteOrgCmd)
	OrgCmd.AddCommand(getOrgCmd)
	OrgCmd.AddCommand(memOrgCmd)
	memOrgCmd.AddCommand(addOrgMemCmd)
	memOrgCmd.AddCommand(remOrgMemCmd)
	memOrgCmd.AddCommand(listOrgMemCmd)
	memOrgCmd.AddCommand(changeOrgMemRoleCmd)

	listOrgCmd.Flags().BoolP("quiet", "q", false, "Only display Organization Name")

	createOrgCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	createOrgCmd.Flags().StringVar(&email, "email", email, "Email ID")

	deleteOrgCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")

	getOrgCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")

	addOrgMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	addOrgMemCmd.Flags().StringVar(&member, "member", member, "Member Name")

	changeOrgMemRoleCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	changeOrgMemRoleCmd.Flags().StringVar(&member, "member", member, "Member Name")
	changeOrgMemRoleCmd.Flags().StringVar(&role, "role", role, "Role Name")

	remOrgMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	remOrgMemCmd.Flags().StringVar(&member, "member", member, "Member Name")

	listOrgMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	listOrgMemCmd.Flags().BoolP("quiet", "q", false, "Only display Member Name")
}

// listOrg validates the input command line arguments and lists available organizations
// by invoking the corresponding rpc/storage method
func listOrg(amp *cli.AMP, cmd *cobra.Command) (err error) {
	request := &account.ListOrganizationsRequest{}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.ListOrganizations(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}

	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		return fmt.Errorf("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, org := range reply.Organizations {
			fmt.Println(org.Name)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	for _, org := range reply.Organizations {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", org.Name, org.Email, ConvertTime(org.CreateDt))
	}
	w.Flush()
	return nil
}

// createOrg validates the input command line arguments and creates an organization
// by invoking the corresponding rpc/storage method
func createOrg(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = GetEmailAddress()
	}

	request := &account.CreateOrganizationRequest{
		Name:  organization,
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.CreateOrganization(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "The organization %s has been successfully created.", organization)
	return nil
}

// deleteOrg validates the input command line arguments and deletes an organization
// by invoking the corresponding rpc/storage method
func deleteOrg(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}

	request := &account.DeleteOrganizationRequest{
		Name: organization,
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
func getOrg(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}
	request := &account.GetOrganizationRequest{
		Name: organization,
	}

	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.GetOrganization(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	fmt.Fprintf(w, "%s\t%s\t%s\n", reply.Organization.Name, reply.Organization.Email, ConvertTime(reply.Organization.CreateDt))
	w.Flush()
	return nil
}

// memberOrg validates the input command line arguments and retrieves info about members of an organization
// by invoking the corresponding rpc/storage method
func memberOrg() (err error) {
	manager.printf(colWarn, "Choose a command for member operations.\nUse amp org member -h for help.")
	return nil
}

// addOrgMem validates the input command line arguments and adds members to an organization
// by invoking the corresponding rpc/storage method
func addOrgMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}
	if cmd.Flag("member").Changed {
		member = cmd.Flag("member").Value.String()
	} else {
		fmt.Print("member name: ")
		member = GetName()
	}

	request := &account.AddUserToOrganizationRequest{
		OrganizationName: organization,
		UserName:         member,
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

// remOrgMem validates the input command line arguments and removes members from an organization
// by invoking the corresponding rpc/storage method
func remOrgMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}
	if cmd.Flag("member").Changed {
		member = cmd.Flag("member").Value.String()
	} else {
		fmt.Print("member name: ")
		member = GetName()
	}

	request := &account.RemoveUserFromOrganizationRequest{
		OrganizationName: organization,
		UserName:         member,
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

func changeOrgMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}
	if cmd.Flag("member").Changed {
		member = cmd.Flag("member").Value.String()
	} else {
		fmt.Print("member name: ")
		member = GetName()
	}
	if cmd.Flag("role").Changed {
		role = cmd.Flag("role").Value.String()
	} else {
		fmt.Print("organization role: ")
		fmt.Scanln(&role)
	}

	orgRole := schema.OrganizationRole_ORGANIZATION_MEMBER
	switch role {
	case "owner":
		orgRole = schema.OrganizationRole_ORGANIZATION_OWNER
	case "member":
		orgRole = schema.OrganizationRole_ORGANIZATION_MEMBER
	default:
		manager.fatalf("invalid organization role: %s. Please specify 'owner' or 'member' as role value.", role)
	}

	request := &account.ChangeOrganizationMemberRoleRequest{
		OrganizationName: organization,
		UserName:         member,
		Role:             orgRole,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.ChangeOrganizationMemberRole(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Role has been changed successfully.")
	return nil
}

// listOrgMem validates the input command line arguments and lists all members of an organization
// by invoking the corresponding rpc/storage method
func listOrgMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("org").Changed {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = GetName()
	}

	request := &account.GetOrganizationRequest{
		Name: organization,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.GetOrganization(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}

	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		return fmt.Errorf("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, member := range reply.Organization.Members {
			fmt.Println(member.Name)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tROLE\t")
	for _, user := range reply.Organization.Members {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Role)
	}
	w.Flush()
	return nil
}

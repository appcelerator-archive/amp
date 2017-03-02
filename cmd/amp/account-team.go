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

// TeamCmd is the main command for attaching team sub-commands.
var (
	listTeamCmd = &cobra.Command{
		Use:     "list",
		Short:   "List team",
		Long:    `The list command lists all available teams in an organization.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(AMP, cmd)
		},
	}

	createTeamCmd = &cobra.Command{
		Use:   "create",
		Short: "Create team",
		Long:  `The create command creates a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(AMP, cmd)
		},
	}

	deleteTeamCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete team",
		Long:    `The delete command deletes a team in an organization.`,
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTeam(AMP, cmd)
		},
	}

	getTeamCmd = &cobra.Command{
		Use:   "get",
		Short: "Get team info",
		Long:  `The get command retrieves details of a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(AMP, cmd)
		},
	}

	memTeamCmd = &cobra.Command{
		Use:   "member",
		Short: "Member-related operations in a team",
		Long:  `The member command manages member-related operations of a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return memberTeam()
		},
	}

	addTeamMemCmd = &cobra.Command{
		Use:   "add",
		Short: "Add members to team",
		Long:  `The add command adds members to a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(AMP, cmd)
		},
	}

	remTeamMemCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Remove members from team",
		Long:    `The remove command removes members from a team in an organization.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeamMem(AMP, cmd)
		},
	}

	listTeamMemCmd = &cobra.Command{
		Use:   "ls",
		Short: "List members of team",
		Long:  `The list command lists members of a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamMem(AMP, cmd)
		},
	}

	team string
)

func init() {
	TeamCmd.AddCommand(listTeamCmd)
	TeamCmd.AddCommand(createTeamCmd)
	TeamCmd.AddCommand(deleteTeamCmd)
	TeamCmd.AddCommand(getTeamCmd)
	TeamCmd.AddCommand(memTeamCmd)
	memTeamCmd.AddCommand(addTeamMemCmd)
	memTeamCmd.AddCommand(remTeamMemCmd)
	memTeamCmd.AddCommand(listTeamMemCmd)

	listTeamCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")

	createTeamCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	createTeamCmd.Flags().StringVar(&team, "team", team, "Team Name")

	deleteTeamCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	deleteTeamCmd.Flags().StringVar(&team, "team", team, "Team Name")

	getTeamCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	getTeamCmd.Flags().StringVar(&team, "team", team, "Team Name")

	addTeamMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	addTeamMemCmd.Flags().StringVar(&team, "team", team, "Team Name")
	addTeamMemCmd.Flags().StringVar(&member, "member", member, "Member Name")

	remTeamMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	remTeamMemCmd.Flags().StringVar(&team, "team", team, "Team Name")
	remTeamMemCmd.Flags().StringVar(&member, "member", member, "Member Name")

	listTeamMemCmd.Flags().StringVar(&organization, "org", organization, "Organization Name")
	listTeamMemCmd.Flags().StringVar(&team, "team", team, "Team Name")
}

// listTeam validates the input command line arguments and lists available teams
// by invoking the corresponding rpc/storage method
func listTeam(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will list teams in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}

	request := &account.ListTeamsRequest{
		OrganizationName: organization,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.ListTeams(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	fmt.Fprintln(w, "----\t-------\t")
	for _, team := range reply.Teams {
		teamCreate, err := strconv.ParseInt(strconv.FormatInt(team.CreateDt, 10), 10, 64)
		if err != nil {
			panic(err)
		}
		teamCreateTime := time.Unix(teamCreate, 0)
		fmt.Fprintf(w, "%s\t%s\n", team.Name, teamCreateTime)
	}
	w.Flush()
	return nil
}

// createTeam validates the input command line arguments and creates a team in an organization
// by invoking the corresponding rpc/storage method
func createTeam(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will create a team in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}

	request := &account.CreateTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.CreateTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully created team %s in organization %s.", team, organization)
	return nil
}

// deleteTeam validates the input command line arguments and deletes a team in an organization
// by invoking the corresponding rpc/storage method
func deleteTeam(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will delete a team from an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}

	request := &account.DeleteTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.DeleteTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully deleted team %s from organization %s.", team, organization)
	return nil
}

// getTeam validates the input command line arguments and retrieves info of a team in an organization
// by invoking the corresponding rpc/storage method
func getTeam(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will get details of a team in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}

	request := &account.GetTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.GetTeam(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	fmt.Fprintln(w, "----\t-------\t")
	teamCreate, err := strconv.ParseInt(strconv.FormatInt(reply.Team.CreateDt, 10), 10, 64)
	if err != nil {
		panic(err)
	}
	teamCreateTime := time.Unix(teamCreate, 0)
	fmt.Fprintf(w, "%s\t%s\n", reply.Team.Name, teamCreateTime)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "MEMBER NAME\tROLE\t")
	fmt.Fprintln(w, "-----------\t----\t")
	for _, mem := range reply.Team.Members {
		fmt.Fprintf(w, "%s\t%s\n", mem.Name, mem.Role)
	}
	w.Flush()
	return nil
}

// memberTeam validates the input command line arguments and retrieves info about members of a team in an organization
// by invoking the corresponding rpc/storage method
func memberTeam() (err error) {
	manager.printf(colWarn, "Choose a command for member operations.\nUse amp team member -h for help.")
	return nil
}

// addTeamMem validates the input command line arguments and adds members to a team
// by invoking the corresponding rpc/storage method
func addTeamMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will add members to a team in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}
	if cmd.Flags().Changed("member") {
		member = cmd.Flag("member").Value.String()
	} else {
		member = getUserName()
	}

	request := &account.AddUserToTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
		UserName:         member,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.AddUserToTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Member(s) have been added to team %s successfully.", team)
	return nil
}

// removeTeamMem validates the input command line arguments and removes members from a team
// by invoking the corresponding rpc/storage method
func removeTeamMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will remove members from a team in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}
	if cmd.Flags().Changed("member") {
		member = cmd.Flag("member").Value.String()
	} else {
		member = getUserName()
	}

	request := &account.RemoveUserFromTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
		UserName:         member,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.RemoveUserFromTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Member(s) have been removed from team %s successfully.", team)
	return nil
}

// listMem validates the input command line arguments and lists members of a team
// by invoking the corresponding rpc/storage method
func listTeamMem(amp *cli.AMP, cmd *cobra.Command) (err error) {
	manager.printf(colRegular, "This will list members of a team in an organization.")
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		organization = getOrgName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		team = getTeamName()
	}

	request := &account.GetTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.GetTeam(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USERNAME\tROLE\t")
	fmt.Fprintln(w, "--------\t----\t")
	for _, user := range reply.Team.Members {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Role)
	}
	w.Flush()
	return nil
}

func getTeamName() (team string) {
	fmt.Print("team name: ")
	fmt.Scanln(&team)
	team = strings.TrimSpace(team)
	err := schema.CheckName(team)
	if err != nil {
		manager.printf(colWarn, err.Error())
		return getTeamName()
	}
	return
}

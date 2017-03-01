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
			return listTeam(AMP)
		},
	}

	createTeamCmd = &cobra.Command{
		Use:   "create",
		Short: "Create team",
		Long:  `The create command creates a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(AMP)
		},
	}

	deleteTeamCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete team",
		Long:    `The delete command deletes a team in an organization.`,
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTeam(AMP)
		},
	}

	getTeamCmd = &cobra.Command{
		Use:   "get",
		Short: "Get team info",
		Long:  `The get command retrieves details of a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(AMP)
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

	addMemCmd = &cobra.Command{
		Use:   "add",
		Short: "Add members to team",
		Long:  `The add command adds members to a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(AMP)
		},
	}

	removeMemCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Remove members from team",
		Long:    `The remove command removes members from a team in an organization.`,
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeamMem(AMP)
		},
	}

	listMemCmd = &cobra.Command{
		Use:   "ls",
		Short: "List members of team",
		Long:  `The list command lists members of a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeamMem(AMP)
		},
	}
)

func init() {
	TeamCmd.AddCommand(listTeamCmd)
	TeamCmd.AddCommand(createTeamCmd)
	TeamCmd.AddCommand(deleteTeamCmd)
	TeamCmd.AddCommand(getTeamCmd)
	TeamCmd.AddCommand(memTeamCmd)
	memTeamCmd.AddCommand(addMemCmd)
	memTeamCmd.AddCommand(removeMemCmd)
	memTeamCmd.AddCommand(listMemCmd)
}

// listTeam validates the input command line arguments and lists available teams
// by invoking the corresponding rpc/storage method
func listTeam(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will list teams in an organization.")
	orgName := getOrgName()
	request := &account.ListTeamsRequest{
		OrganizationName: orgName,
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
func createTeam(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will create a team in an organization.")
	orgName := getOrgName()
	team := getTeamName()
	request := &account.CreateTeamRequest{
		OrganizationName: orgName,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.CreateTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully created team %s in organization %s.", team, orgName)
	return nil
}

// deleteTeam validates the input command line arguments and deletes a team in an organization
// by invoking the corresponding rpc/storage method
func deleteTeam(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will delete a team from an organization.")
	orgName := getOrgName()
	team := getTeamName()
	request := &account.DeleteTeamRequest{
		OrganizationName: orgName,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.DeleteTeam(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully deleted team %s from organization %s.", team, orgName)
	return nil
}

// getTeam validates the input command line arguments and retrieves info of a team in an organization
// by invoking the corresponding rpc/storage method
func getTeam(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will get details of a team in an organization.")
	orgName := getOrgName()
	team := getTeamName()
	request := &account.GetTeamRequest{
		OrganizationName: orgName,
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
	manager.printf(colWarn, "Choose a command for member operations.\nUse amp account team member -h for help.")
	return nil
}

// addTeamMem validates the input command line arguments and adds members to a team
// by invoking the corresponding rpc/storage method
func addTeamMem(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will add members to a team in an organization.")
	orgName := getOrgName()
	team := getTeamName()
	name := getUserName()
	request := &account.AddUserToTeamRequest{
		OrganizationName: orgName,
		TeamName:         team,
		UserName:         name,
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
func removeTeamMem(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will remove members from a team in an organization.")
	orgName := getOrgName()
	team := getTeamName()
	name := getUserName()
	request := &account.RemoveUserFromTeamRequest{
		OrganizationName: orgName,
		TeamName:         team,
		UserName:         name,
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
func listTeamMem(amp *cli.AMP) (err error) {
	manager.printf(colRegular, "This will list members of a team in an organization.")
	orgName := getOrgName()
	team := getTeamName()
	request := &account.GetTeamRequest{
		OrganizationName: orgName,
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

package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// TeamCmd is the main command for attaching team sub-commands.
var (
	listTeamCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List team",
		Example: "amp team ls --org=dummyorg",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(AMP, cmd)
		},
	}

	createTeamCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create team",
		Example: "amp team create --org=randomorg --team=coolteam ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(AMP, cmd)
		},
	}

	deleteTeamCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove team",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTeam(AMP, cmd)
		},
	}

	getTeamCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get team info",
		Example: "amp team get --org=fakeorg --team=funteam",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTeam(AMP, cmd)
		},
	}

	memTeamCmd = &cobra.Command{
		Use:   "member",
		Short: "Member-related operations in a team",
		RunE: func(cmd *cobra.Command, args []string) error {
			return memberTeam(AMP, cmd)
		},
	}

	addTeamMemCmd = &cobra.Command{
		Use:     "add",
		Short:   "Add members to team",
		Example: "amp team member add --org=fakeorg --team=funteam --member=rachel",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMem(AMP, cmd)
		},
	}

	remTeamMemCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove members from team",
		Example: "amp team member rm --org=randomorg --team=coolteam --member=joey \namp team member del --org=randomorg --team=coolteam --member=joey",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeamMem(AMP, cmd)
		},
	}

	listTeamMemCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List members of team",
		Example: "amp team member ls --org=dummyorg --team=geekteam",
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
	listTeamCmd.Flags().BoolP("quiet", "q", false, "Only display Team Name")

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
	listTeamMemCmd.Flags().BoolP("quiet", "q", false, "Only display Team Name")
}

// listTeam validates the input command line arguments and lists available teams
// by invoking the corresponding rpc/storage method
func listTeam(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}

	request := &account.ListTeamsRequest{
		OrganizationName: organization,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.ListTeams(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		return fmt.Errorf("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, team := range reply.Teams {
			fmt.Println(team.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	for _, team := range reply.Teams {
		fmt.Fprintf(w, "%s\t%s\n", team.Name, convertTime(team.CreateDt))
	}
	w.Flush()
	return nil
}

// createTeam validates the input command line arguments and creates a team in an organization
// by invoking the corresponding rpc/storage method
func createTeam(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}

	request := &account.CreateTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.CreateTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Successfully created team %s in organization %s.", team, organization)
	return nil
}

// deleteTeam validates the input command line arguments and deletes a team in an organization
// by invoking the corresponding rpc/storage method
func deleteTeam(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}

	request := &account.DeleteTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.DeleteTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Successfully deleted team %s from organization %s.", team, organization)
	return nil
}

// getTeam validates the input command line arguments and retrieves info of a team in an organization
// by invoking the corresponding rpc/storage method
func getTeam(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}

	request := &account.GetTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.GetTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "TEAM\tCREATED\t")
	fmt.Fprintf(w, "%s\t%s\n", reply.Team.Name, convertTime(reply.Team.CreateDt))
	w.Flush()
	return nil
}

// memberTeam validates the input command line arguments and retrieves info about members of a team in an organization
// by invoking the corresponding rpc/storage method
func memberTeam(amp *cli.AMP, cmd *cobra.Command) (err error) {
	mgr.Warn("Choose a command for member operations.\nUse amp team member -h for help.")
	return nil
}

// addTeamMem validates the input command line arguments and adds members to a team
// by invoking the corresponding rpc/storage method
func addTeamMem(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}
	if cmd.Flags().Changed("member") {
		member = cmd.Flag("member").Value.String()
	} else {
		fmt.Print("member name: ")
		member = getName()
	}

	request := &account.AddUserToTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
		UserName:         member,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.AddUserToTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Member(s) have been added to team %s successfully.", team)
	return nil
}

// removeTeamMem validates the input command line arguments and removes members from a team
// by invoking the corresponding rpc/storage method
func removeTeamMem(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}
	if cmd.Flags().Changed("member") {
		member = cmd.Flag("member").Value.String()
	} else {
		fmt.Print("member name: ")
		member = getName()
	}

	request := &account.RemoveUserFromTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
		UserName:         member,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.RemoveUserFromTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Member(s) have been removed from team %s successfully.", team)
	return nil
}

// listTeamMem validates the input command line arguments and lists members of a team
// by invoking the corresponding rpc/storage method
func listTeamMem(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flags().Changed("org") {
		organization = cmd.Flag("org").Value.String()
	} else {
		fmt.Print("organization: ")
		organization = getName()
	}
	if cmd.Flags().Changed("team") {
		team = cmd.Flag("team").Value.String()
	} else {
		fmt.Print("team: ")
		team = getName()
	}

	request := &account.GetTeamRequest{
		OrganizationName: organization,
		TeamName:         team,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.GetTeam(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		mgr.Error("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, member := range reply.Team.Members {
			fmt.Println(member.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tROLE\t")
	for _, user := range reply.Team.Members {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Role)
	}
	w.Flush()
	return nil
}

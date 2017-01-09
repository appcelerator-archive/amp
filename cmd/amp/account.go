package main

import (
	"github.com/spf13/cobra"
)

var (
	// Interactive
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "log in to amp",
	}

	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "account operations",
	}
	// The rest of the commands are subcommands

	// Interactive
	signUpCmd = &cobra.Command{
		Use:   "signup",
		Short: "Create a new account and login",
	}

	// Interactive
	verifyCmd = &cobra.Command{
		Use:   "verify [CODE]",
		Short: "verify email using code",
	}

	switchRoleCmd = &cobra.Command{
		Use:   "switch [ORGANIZATION]",
		Short: "Switch primary organization",
	}

	createOrganizationCmd = &cobra.Command{
		Use:   "create organization [NAME] [EMAIL]",
		Short: "Create an organization",
	}

	listUsersCmd = &cobra.Command{
		Use:   "list users [ORGANIZATION] [TEAM]",
		Short: "list users, optionally filter by organization and team",
	}

	listOrganizationsCmd = &cobra.Command{
		Use:   "list organizations",
		Short: "list organizations",
	}

	listTeamsCmd = &cobra.Command{
		Use:   "list teams [ORGANIZATION]",
		Short: "list teams by organization",
	}

	listPermissionsCmd = &cobra.Command{
		Use:   "list permissions [ORGANIZATION] [TEAM]",
		Short: "list permissions by team",
	}

	infoCmd = &cobra.Command{
		Use:   "info [name]",
		Short: "list account information",
	}

	editCmd = &cobra.Command{
		Use:   "edit [name]",
		Short: "edits account information",
	}

	deleteCmd = &cobra.Command{
		Use:   "delete [name]",
		Short: "deletes account",
	}

	addOrganizationMembersCmd = &cobra.Command{
		Use:   "add organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Add users to an organization",
	}

	addTeamMembersCmd = &cobra.Command{
		Use:   "add team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Add users to a team",
	}

	removeOrganizationMembersCmd = &cobra.Command{
		Use:   "remove organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Remove users from an organization",
	}

	removeTeamMembersCmd = &cobra.Command{
		Use:   "remove team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Remove users from a team",
	}

	grantPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM]",
		Short: "Grant permission to a team for a resource",
	}

	editPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM] ",
		Short: "Edit a permission",
	}

	revokePermissionCmd = &cobra.Command{
		Use:   "revoke [ORGANIZATION] [TEAM] [RESOURCE_ID]",
		Short: "Revoke permission to a team for a resource",
	}

	transferOwnershipCmd = &cobra.Command{
		Use:   "transfer [RESOURCE_ID] [ORGANIZATION] [TEAM]",
		Short: "Transfer ownership of a resource to an organization or team",
	}
)

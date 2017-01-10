package main

import "github.com/spf13/cobra"

var (
	// Interactive
	loginStubCmd = &cobra.Command{
		Use:   "login",
		Short: "log in to amp",
	}
	/*
		     I am a user
		     I want to login

		     $ amp login
			   username: {username}
			   password: *{password}*

		     result: session token set
	*/

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
		/*
			I am an account owner
			I want information about a specific account

			$ amp account info <account-name>

			Result: the account information is listed. There could be two different responses.
			a) A user viewing their own account information:
				* account name
				* email
				* Belongs to <organization>
					* Belongs to <team>
				* Billing information
				* Other Settings

			b) A user viewing a different account information
			* account name
			* email
			* Belongs to <organization>
				* Belongs to <team>
		*/
	}

	editCmd = &cobra.Command{
		Use:   "edit [name]",
		Short: "edits account information",
		/*
			I am an account owner
			I want to edit information about a specific account

			$ amp account edit <account-name>

			Result: the account information is modified

		*/
	}

	deleteCmd = &cobra.Command{
		Use:   "delete [name]",
		Short: "deletes account",
		/*
			I am an account owner
			I want to remove a specific account

			$ amp account delete <account-name>

			Result: the account information is deleted

		*/
	}

	addOrganizationMembersCmd = &cobra.Command{
		Use:   "add organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Add users to an organization",
		/*
			I am an owner team member for an organization account.
			I want to add members (comma-separated names if multiple) to the organization.

			$ amp account add organization members <org-name> <member-account-name(s)>

			Result: the member(s) are added to the organization

		*/
	}

	addTeamMembersCmd = &cobra.Command{
		Use:   "add team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Add users to a team",
		/*
			I am an owner team member for an organization account.
			I want to add team members (comma-separated names if multiple) to the organization.

			$ amp account add team members <org-name> <team-name> <member-account-name(s)>

			Result: the member(s) are added to the team

		*/
	}

	removeOrganizationMembersCmd = &cobra.Command{
		Use:   "remove organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Remove users from an organization",
		/*
			I am an owner team member for an organization account.
			I want to remove members (comma-separated names if multiple) from the organization.

			$ amp account remove organization members <org-name> <member-account-name(s)>

			Result: the specified member(s) are removed from the organization

		*/
	}

	removeTeamMembersCmd = &cobra.Command{
		Use:   "remove team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Remove users from a team",
		/*
			I am an owner team member for an organization account.
			I want to remove team members  (comma-separated names if multiple) from a team in the organization.

			$ amp account remove team members <org-name> <team-name> <member-account-name(s)>

			Result: the specified member(s) are removed from the team

		*/
	}

	grantPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM]",
		Short: "Grant permission to a team for a resource",
		/*
			I am an account owner.
			I want to grant permission to a team (which is part of an organization) for a resource.
			Level can be Read, Write, Read/Write, Delete

			$ amp account grant <resource-id> <level> <org-name> <team-name>

			Result: the team is granted permission for the requested resource

		*/
	}

	editPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM] ",
		Short: "Edit a permission for a resource",
		/*
			I am an account owner.
			I want to edit permissions of a team (which is part of an organization) for a resource.
			Level can be Read, Write, Read/Write, Delete

			$ amp account grant <resource-id> <level> <org-name> <team-name>

			Result: the permissions of the team are modified on the requested resource

		*/
	}

	revokePermissionCmd = &cobra.Command{
		Use:   "revoke [RESOURCE_ID] [ORGANIZATION] [TEAM]",
		Short: "Revoke permission from a team for a resource",
		/*
			I am an account owner.
			I want to revoke permission from a team (which is part of an organization) for a resource.

			$ amp account revoke <resource-id> <org-name> <team-name>

			Result: the permission is revoked from the team for the requested resource

		*/
	}

	transferOwnershipCmd = &cobra.Command{
		Use:   "transfer [RESOURCE_ID] [ORGANIZATION]",
		Short: "Transfer ownership of a resource to a different organization",
		/*
			I am a resource owner.
			I want to transfer a resource to a different organization.

			$ amp account transfer <resource-id> <org-name>

			Result: the account is transferred to the specified organization

		*/
	}
)

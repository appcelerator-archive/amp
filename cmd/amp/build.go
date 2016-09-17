package main

// import (
// 	"fmt"
// 	"io"
// 	"strconv"

// 	"github.com/spf13/cobra"
// 	"google.golang.org/grpc"
// )

// var buildCmd = &cobra.Command{
// 	Use:   "build",
// 	Short: "Manage the amp build service",
// 	Long:  `Register projects, list builds etc . . . in the amp build service.`,
// }

// func init() {
// 	pingCmd := &cobra.Command{
// 		Use:   "ping",
// 		Short: "Ping the amp build service",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("pong")
// 		},
// 	}

// 	registerCmd := &cobra.Command{
// 		Use:   "register [repo to register]",
// 		Short: "Register a repo as a project",
// 		Long:  `Register a repo as a project on the amp build service`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) != 1 {
// 				fmt.Println("must specify repo to register")
// 				return
// 			}

// 			repo := args[0]
// 			project, err := AMP.RegisterProject(repo)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			fmt.Println("registered https://build.amp.appcelerator.io/p/" + project.Owner + "/" + project.Name)
// 		},
// 	}

// 	removeCmd := &cobra.Command{
// 		Use:   "remove [repo to remove]",
// 		Short: "Remove a repo as a project",
// 		Long:  `Remove a repo as a project on the amp build service`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) != 1 {
// 				fmt.Println("must specify repo to remove")
// 				return
// 			}

// 			repo := args[0]
// 			project, err := AMP.RemoveProject(repo)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			fmt.Println("removed https://build.amp.appcelerator.io/p/" + project.Owner + "/" + project.Name)
// 		},
// 	}

// 	listProjectsCmd := &cobra.Command{
// 		Use:   "listprojects",
// 		Short: "Lists projects",
// 		Long:  `Lists projects on the amp build service`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			flags := cmd.Flags()
// 			latest, err := flags.GetBool("latest")
// 			if err != nil {
// 				fmt.Println("could not process flags")
// 				return
// 			}

// 			quiet, err := flags.GetBool("quiet")
// 			if err != nil {
// 				fmt.Println("could not process flags")
// 				return
// 			}

// 			organization, err := flags.GetString("organization")
// 			if err != nil {
// 				fmt.Println("could not process flags")
// 				return
// 			}

// 			projects, err := AMP.ListProjects(organization, latest)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			for _, p := range projects.Projects {
// 				if quiet {
// 					fmt.Println(p.Owner + "/" + p.Name)
// 				} else {
// 					fmt.Println("https://build.amp.appcelerator.io/p/"+p.Owner+"/"+p.Name, p.Status)
// 				}
// 			}
// 		},
// 	}
// 	listProjectsCmd.Flags().StringP("organization", "o", "", "filter projects by organization")
// 	listProjectsCmd.Flags().BoolP("latest", "l", false, "only return the latest project")
// 	listProjectsCmd.Flags().BoolP("quiet", "q", false, "only display the project id")

// 	listBuildsCmd := &cobra.Command{
// 		Use:   "listbuilds [project to list builds from]",
// 		Short: "Lists builds",
// 		Long:  `Lists the builds for a given project`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) != 1 {
// 				fmt.Println("must specify project to list build from")
// 				return
// 			}

// 			repo := args[0]
// 			flags := cmd.Flags()
// 			latest, err := flags.GetBool("latest")
// 			if err != nil {
// 				fmt.Println("could not process flags")
// 				return
// 			}

// 			quiet, err := flags.GetBool("quiet")
// 			if err != nil {
// 				fmt.Println("could not process flags")
// 				return
// 			}

// 			builds, err := AMP.ListBuilds(repo, latest)
// 			for i, b := range builds.Builds {
// 				if quiet {
// 					fmt.Println(b.Owner + "/" + b.Name + "/" + b.Sha)
// 				} else {
// 					fmt.Println("https://build.amp.appcelerator.io/p/"+b.Owner+"/"+b.Name+"/"+strconv.Itoa(i), b.CommitMessage, b.Status)
// 				}
// 				if latest {
// 					return
// 				}
// 			}
// 		},
// 	}
// 	listBuildsCmd.Flags().BoolP("latest", "l", false, "only return the latest build")
// 	listBuildsCmd.Flags().BoolP("quiet", "q", false, "only display the build id")

// 	logsCmd := &cobra.Command{
// 		Use:   "logs [build to print the logs from]",
// 		Short: "Prints the logs for a build",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) != 1 {
// 				fmt.Println("must specify build to print logs from")
// 				return
// 			}

// 			buildid := args[0]
// 			logs, err := AMP.BuildLog(buildid)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			for {
// 				log, err := logs.Recv()
// 				if err == io.EOF {
// 					break
// 				}
// 				if err != nil {
// 					if grpc.ErrorDesc(err) == "EOF" {
// 						break
// 					} else {
// 						fmt.Println(err)
// 						return
// 					}
// 				}
// 				fmt.Print(log.Message)
// 			}
// 		},
// 	}

// 	rebuildCmd := &cobra.Command{
// 		Use:   "rebuild [build to rebuild]",
// 		Short: "Triggers a rebuild of a build",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) != 1 {
// 				fmt.Println("must specify build to rebuild")
// 				return
// 			}

// 			buildid := args[0]
// 			build, err := AMP.Rebuild(buildid)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			fmt.Println("rebuilding https://build.amp.appcelerator.io/p/" + build.Owner + "/" + build.Name + "/" + build.Sha)
// 		},
// 	}

// 	buildCmd.AddCommand(pingCmd)
// 	buildCmd.AddCommand(registerCmd)
// 	buildCmd.AddCommand(removeCmd)
// 	buildCmd.AddCommand(listProjectsCmd)
// 	buildCmd.AddCommand(listBuildsCmd)
// 	buildCmd.AddCommand(logsCmd)
// 	buildCmd.AddCommand(rebuildCmd)

// 	RootCmd.AddCommand(buildCmd)
// }

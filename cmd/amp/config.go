package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

func init() {
	// configCmd represents the Config command
	configCmd := &cobra.Command{
		Use:     "config",
		Short:   "Display or update the current configuration",
		Example: "AmpAddress",
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				// Display the configuration (not all has to be displayed)
				w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
				fmt.Fprintf(w, "%s\n", "PARAMETER\tVALUE")
				fmt.Fprintf(w, "%s\t%s\n", "AmpAddress", Config.AmpAddress)
				fmt.Fprintf(w, "%s\t%s\n", "ServerPort", Config.ServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "AdminServerPort", Config.AdminServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "CmdTheme", Config.CmdTheme)
				fmt.Fprintf(w, "%s\t%t\n", "Verbose", Config.Verbose)
				w.Flush()
			case 1:
				// Display key
				s := structs.New(Config)
				f, ok := s.FieldOk(strings.Title(args[0]))
				if !ok {
					//log.Fatalf("Field %s not found", strings.Title(args[0]))
					mgr.Fatal("field %s not found", strings.Title(args[0]))
				}
				fmt.Println(f.Value())
			case 2:
				// Update key
				s := structs.New(Config)
				f, ok := s.FieldOk(strings.Title(args[0]))
				if !ok {
					//log.Fatalf("Field %s not found", strings.Title(args[0]))
					mgr.Fatal("field %s not found", strings.Title(args[0]))
				}
				switch f.Kind().String() {
				case "bool":
					b, err := strconv.ParseBool(args[1])
					if err != nil {
						//log.Fatalf("Could not parse %s as bool", args[1])
						mgr.Fatal("could not parse %s as bool", args[1])
					}
					f.Set(b)
				case "string":
					f.Set(args[1])
				default:
					//log.Fatal("Unsupported field type")
					mgr.Fatal("unsupported field type")
				}
				err := cli.SaveConfiguration(Config)
				if err != nil {
					//log.Fatal("Failed to save config")
					mgr.Fatal("failed to save config")
				}
				fmt.Println(f.Value())
			default:
				//log.Fatal("Too many arguments")
				mgr.Fatal("too many arguments")
			}
		},
	}
	RootCmd.AddCommand(configCmd)
}

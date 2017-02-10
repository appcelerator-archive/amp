package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	conf "github.com/appcelerator/amp/pkg/config"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

func init() {
	// configCmd represents the Config command
	configCmd := &cobra.Command{
		Use:   "config [KEY] [VALUE]",
		Short: "Display or update the current configuration",
		Long: `The config command displays/updates the current configuration.
No arguments: display the current configuration.
One argument: display the configuration key value.
Two arguments: set the key to the value.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				// Display the configuration (not all has to be displayed)
				w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
				fmt.Fprintf(w, "%s\n", "PARAMETER\tVALUE")
				fmt.Fprintf(w, "%s\t%s\n", "AmpAddress", Config.AmpAddress)
				fmt.Fprintf(w, "%s\t%s\n", "ServerPort", Config.ServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "AdminServerPort", Config.AdminServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "WebMailServerPort", Config.WebMailServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "CmdTheme", Config.CmdTheme)
				fmt.Fprintf(w, "%s\t%t\n", "Verbose", Config.Verbose)
				fmt.Fprintf(w, "%s\t%s\n", "EmailServerAddress", Config.EmailServerAddress)
				fmt.Fprintf(w, "%s\t%s\n", "EmailServerPort", Config.EmailServerPort)
				fmt.Fprintf(w, "%s\t%s\n", "EmailSender", Config.EmailSender)
				w.Flush()
			case 1:
				// Display key
				s := structs.New(Config)
				f, ok := s.FieldOk(strings.Title(args[0]))
				if !ok {
					log.Fatalf("Field %s not found", strings.Title(args[0]))
				}
				fmt.Println(f.Value())
			case 2:
				// Update key
				s := structs.New(Config)
				f, ok := s.FieldOk(strings.Title(args[0]))
				if !ok {
					log.Fatalf("Field %s not found", strings.Title(args[0]))
				}
				switch f.Kind().String() {
				case "bool":
					b, err := strconv.ParseBool(args[1])
					if err != nil {
						log.Fatalf("Could not parse %s as bool", args[1])
					}
					f.Set(b)
				case "string":
					f.Set(args[1])
				default:
					log.Fatal("Unsupported field type")
				}
				err := conf.SaveConfiguration(Config)
				if err != nil {
					log.Fatal("Failed to save config")
				}
				fmt.Println(f.Value())
			default:
				log.Fatal("Too many arguments")
			}
		},
	}
	RootCmd.AddCommand(configCmd)
}

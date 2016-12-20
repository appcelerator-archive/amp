package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

func init() {
	// configCmd represents the Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Display or update the current configuration",
		Long: `No arguments: display the current configuration.
One argument: display the configuration key value.
Two arguments: set the key to the value.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				// Display the full configuration
				j, err := json.MarshalIndent(structs.Map(Config), "", "  ")
				if err != nil {
					fmt.Println("error: ", err)
				} else {
					fmt.Println(string(j))
				}
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
				err := cli.SaveConfiguration(Config)
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

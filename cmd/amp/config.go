package main

import (
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
		Short: "Display the current configuration",
		Long:  `Display the current configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				fmt.Println(Config)
			case 1:
				s := structs.New(Config)
				f, ok := s.FieldOk(strings.Title(args[0]))
				if !ok {
					log.Fatalf("Field %s not found", strings.Title(args[0]))
				}
				fmt.Println(f.Value())
			case 2:
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
					log.Fatal("unsupported field type")
				}
				err := cli.SaveConfiguration(Config)
				if err != nil {
					log.Fatal("Failed to save config")
				}
				fmt.Println(f.Value())
			default:
				log.Fatal("too many arguments")
			}
		},
	}
	RootCmd.AddCommand(configCmd)
}

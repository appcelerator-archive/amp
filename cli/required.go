package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NoArgs checks that the command is not passed any parameters
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("unexpected argument: %s\nSee '%s --help'", args[0], cmd.CommandPath())
}

// ExactArgs returns an error if the exact number of args are not passed
func ExactArgs(num int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == num {
			return nil
		}
		return fmt.Errorf(
			"\"%s\" requires exactly %d argument(s).\nSee '%s --help'",
			cmd.CommandPath(),
			num,
			cmd.CommandPath(),
		)
	}
}

// AtLeastArgs returns an error if the min number of args are not passed
func AtLeastArgs(min int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= min {
			return nil
		}
		return fmt.Errorf(
			"\"%s\" requires at least %d argument(s).\nSee '%s --help'",
			cmd.CommandPath(),
			min,
			cmd.CommandPath(),
		)
	}
}

// RangeArgs returns an error if the min and max number of args are not passed
func RangeArgs(min int, max int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= min && len(args) <= max {
			return nil
		}
		return fmt.Errorf(
			"\"%s\" requires at least %d argument(s) and at most %d argument(s).\nSee '%s --help'",
			cmd.CommandPath(),
			min,
			max,
			cmd.CommandPath(),
		)
	}
}

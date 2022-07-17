package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "semver-cli [command]",
		Long:              "semver-cli is set of utility commands for comparing semver versions and constraints",
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(
		newCheckCommand(),
		newEqualCommand(),
		newGreaterThanCommand(),
		newGreaterThanOrEqualCommand(),
		newLessThanCommand(),
		newLessThanOrEqualCommand(),
		newSelectVersionCommand(),
	)

	return cmd
}

func Execute() {
	cmd := newRootCommand()

	if err := cmd.Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
	}
}
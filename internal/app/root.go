package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/stephenwilliams/semver-cli/internal/pkg/version"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "semver-cli [command]",
		Long:              "semver-cli is set of utility commands for comparing semver versions and constraints",
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Version:           version.Version.String(),
	}

	cmd.AddCommand(
		newCheckCommand(),
		newEqualCommand(),
		newGreaterThanCommand(),
		newGreaterThanOrEqualCommand(),
		newLessThanCommand(),
		newLessThanOrEqualCommand(),
		newSelectVersionCommand(),
		newVersionCommand(),
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

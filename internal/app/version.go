package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stephenwilliams/semver-cli/internal/pkg/outputer"
	"github.com/stephenwilliams/semver-cli/internal/pkg/version"
)

func newVersionCommand() *cobra.Command {
	var short bool
	var output = outputer.OutputGeneric
	cmd := &cobra.Command{
		Use:   "version",
		Short: "prints the version information",
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Println(version.Version.Version)
			} else {
				fmt.Println(output.Marshal(version.Version))
			}
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Print just the version.")
	output.AddFlag(cmd)

	return cmd
}

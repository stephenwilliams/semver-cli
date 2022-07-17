package app

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

func check(c, v string, strict bool) (bool, error) {
	ver, err := newVersion(v, strict)
	if err != nil {
		return false, fmt.Errorf("failed to parse version: %w", err)
	}

	con, err := version.NewConstraint(c)
	if err != nil {
		return false, fmt.Errorf("failed to parse constraint: %w", err)
	}

	return con.Check(ver), nil
}

func newCheckCommand() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:     "check <CONSTRAINT> <VERSION>",
		Aliases: []string{"eq"},
		Short:   "Checks a version matches constraint; exits 0 it does, 1 if not.",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := check(args[0], args[1], strict)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			if !result {
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "enforce version that adheres strictly to SemVer specs")

	return cmd
}

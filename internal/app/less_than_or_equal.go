package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func lessThanOrEqual(a, b string, strict bool) (bool, error) {
	vA, err := newVersion(a, strict)
	if err != nil {
		return false, fmt.Errorf("failed to parse version A: %w", err)
	}

	vB, err := newVersion(b, strict)
	if err != nil {
		return false, fmt.Errorf("failed to parse version B: %w", err)
	}

	return vA.GreaterThanOrEqual(vB), nil
}

func newLessThanOrEqualCommand() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:     "less-than-or-equal <A> <B>",
		Aliases: []string{"gte"},
		Short:   "Checks two versions to see if version A <= version B; exits 0 if it is, 1 if not.",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := lessThanOrEqual(args[0], args[1], strict)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			if !result {
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "enforce versions that adhere strictly to SemVer specs")

	return cmd
}

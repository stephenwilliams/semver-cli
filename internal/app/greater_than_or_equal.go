package app

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

func greaterThanOrEqual(a, b string, strict bool, p *regexp.Regexp) (bool, error) {
	vA, err := newVersion(a, strict, p)
	if err != nil {
		return false, fmt.Errorf("failed to parse version A: %w", err)
	}

	vB, err := newVersion(b, strict, p)
	if err != nil {
		return false, fmt.Errorf("failed to parse version B: %w", err)
	}

	return vA.GreaterThanOrEqual(vB), nil
}

func newGreaterThanOrEqualCommand() *cobra.Command {
	var strict bool
	var pattern string

	cmd := &cobra.Command{
		Use:     "greater-than-or-equal <A> <B>",
		Aliases: []string{"gte"},
		Short:   "Checks two versions to see if version A >= version B; exits 0 if it is, 1 if not.",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var p *regexp.Regexp
			if pattern != "" {
				var err error
				p, err = regexp.Compile(pattern)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "failed to parse pattern: %s", err)
					os.Exit(2)
				}

				if p.SubexpIndex("version") == -1 {
					_, _ = fmt.Fprintln(os.Stderr, "invalid pattern. must have a version group", err)
					os.Exit(2)
				}
			}

			result, err := greaterThanOrEqual(args[0], args[1], strict, p)
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
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "pattern for retrieving a version. Must include a version group.")

	return cmd
}

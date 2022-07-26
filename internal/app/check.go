package app

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

func check(c, v string, strict bool, p *regexp.Regexp) (bool, error) {
	ver, err := newVersion(v, strict, p)
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
	var pattern string

	cmd := &cobra.Command{
		Use:     "check <CONSTRAINT> <VERSION>",
		Aliases: []string{"eq"},
		Short:   "Checks a version matches constraint; exits 0 it does, 1 if not.",
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

			result, err := check(args[0], args[1], strict, p)
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
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "pattern for retrieving a version. Must include a version group.")

	return cmd
}

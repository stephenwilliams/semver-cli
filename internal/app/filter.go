package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"

	"github.com/stephenwilliams/semver-cli/internal/pkg/terminal"
)

func filterVersion(c string, versions []string, strict, ignoreErrors, quiet bool, p *regexp.Regexp) ([]string, error) {
	con, err := version.NewConstraint(c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse constraint: %w", err)
	}

	vers, err := newVersions(versions, strict, ignoreErrors, quiet, p)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, v := range vers {
		if con.Check(v.Version) {
			results = append(results, v.Original)
		}
	}

	return results, nil
}

func newFilterVersionCommand() *cobra.Command {
	var strict, ignoreErrors, quiet bool
	var constraint, versionsInput, separator, pattern string

	cmd := &cobra.Command{
		Use:     "filter-versions [FLAGS]",
		Aliases: []string{"fv"},
		Short:   "Filters versions matching the provided constraint.",
		Long: `Filters versions matching the provided constraint. The versions are provided to stdout.
If none match, exits 1.`,
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

			if separator == "NEWLINE" {
				separator = "\n"
			}

			var data []byte
			if versionsInput == "-" {
				var err error
				data, err = terminal.ReadPipedInput(os.Stdin)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
			} else {
				var err error
				data, err = os.ReadFile(versionsInput)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "failed to read input file: %s\n", err)
					os.Exit(2)
				}
			}

			versions := filterEmptyStrings(strings.Split(string(data), separator))

			filtered, err := filterVersion(constraint, versions, strict, ignoreErrors, quiet, p)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			if len(filtered) == 0 {
				_, _ = fmt.Fprintln(os.Stderr, "no versions matched constraint")
				os.Exit(1)
			}

			fmt.Println(strings.Join(filtered, separator))
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "enforce versions that adhere strictly to SemVer specs")
	cmd.Flags().BoolVar(&ignoreErrors, "ignore-errors", false, "ignore errors when parsing versions. Skips invalid version.")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "when ignoring errors parsing versions, don't output an error to stderr.")
	cmd.Flags().StringVar(&constraint, "constraint", "", "the SemVer constraint to use")
	cmd.Flags().StringVar(&versionsInput, "versions", "-", "The versions to filter. - means piped stdin, otherwise assumes its a file path.")
	cmd.Flags().StringVar(&separator, "separator", "NEWLINE", "The separator between versions passed in the input. Defaults to a new line character.")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "pattern for retrieving a version. Must include a version group.")

	return cmd
}

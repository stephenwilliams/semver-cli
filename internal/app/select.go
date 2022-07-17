package app

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"

	"github.com/stephenwilliams/semver-cli/internal/pkg/terminal"
)

func selectVersion(c string, versions []string, strict, latest bool) (string, error) {
	con, err := version.NewConstraint(c)
	if err != nil {
		return "", fmt.Errorf("failed to parse constraint: %w", err)
	}

	vers, err := newVersions(versions, strict)
	if err != nil {
		return "", err
	}

	if latest {
		sort.Sort(sort.Reverse(vers))
	} else {
		sort.Sort(vers)
	}

	for _, v := range vers {
		if con.Check(v) {
			return v.String(), nil
		}
	}

	return "", nil
}

func newSelectVersionCommand() *cobra.Command {
	var strict, latest bool
	var constraint, versionsInput, separator string

	cmd := &cobra.Command{
		Use:     "select-version [FLAGS]",
		Aliases: []string{"sv"},
		Short:   "Selects a version matching the provided constraint.",
		Long: `Selects a version matching the provided constraint. The selected version is provided to stdout.
If none is selected, exits 1.`,
		Run: func(cmd *cobra.Command, args []string) {
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

			versions := filterVersions(strings.Split(string(data), separator))

			selected, err := selectVersion(constraint, versions, strict, latest)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			if selected == "" {
				_, _ = fmt.Fprintln(os.Stderr, "no version matched constraint")
				os.Exit(1)
			}

			fmt.Println(selected)
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "enforce versions that adhere strictly to SemVer specs")
	cmd.Flags().BoolVar(&latest, "latest", true, "selects the latest (highest) version that matches the constraint")
	cmd.Flags().StringVar(&constraint, "constraint", "", "the SemVer constraint to use")
	cmd.Flags().StringVar(&versionsInput, "versions", "-", "The versions to select from. - means piped stdin, otherwise assumes its a file path.")
	cmd.Flags().StringVar(&separator, "separator", "NEWLINE", "The separator between versions passed in the input. Defaults to a new line character.")

	return cmd
}

func filterVersions(versions []string) []string {
	var results []string

	for _, v := range versions {
		v = strings.TrimSpace(v)

		if v != "" {
			results = append(results, v)
		}
	}

	return results
}

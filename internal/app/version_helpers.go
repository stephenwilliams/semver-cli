package app

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

func newVersion(s string, strict bool) (*version.Version, error) {
	if strict {
		return version.NewSemver(s)
	}

	return version.NewVersion(s)
}

func newVersions(versions []string, strict bool) (version.Collection, error) {
	var results version.Collection

	for i, v := range versions {
		ver, err := newVersion(v, strict)
		if err != nil {
			return nil, fmt.Errorf("version at index %d, '%s', failed to parse: %w", i, v, err)
		}

		results = append(results, ver)
	}

	return results, nil
}

package app

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/go-version"
)

func newVersion(s string, strict bool, p *regexp.Regexp) (*version.Version, error) {
	if p != nil {
		sub := p.FindSubmatch([]byte(s))

		if len(sub) == 0 {
			return nil, fmt.Errorf("input failed to match pattern")
		}

		s = string(sub[p.SubexpIndex("version")])
	}

	if strict {
		return version.NewSemver(s)
	}

	return version.NewVersion(s)
}

func newVersions(versions []string, strict, ignoreErrors bool, p *regexp.Regexp) (OriginalValueVersionCollection, error) {
	var results OriginalValueVersionCollection

	for i, v := range versions {
		ver, err := newVersion(v, strict, p)
		if err != nil {
			wErr := fmt.Errorf("version at index %d, '%s', failed to parse: %w", i, v, err)

			if ignoreErrors {
				fmt.Fprintln(os.Stderr, wErr)
				continue
			}

			return nil, wErr
		}

		results = append(results, &OriginalValueVersion{
			Original: v,
			Version:  ver,
		})
	}

	return results, nil
}

type OriginalValueVersion struct {
	Original string
	Version  *version.Version
}

// OriginalValueVersionCollection is a type that implements the
// sort.Interface interface so that versions can be sorted.
type OriginalValueVersionCollection []*OriginalValueVersion

func (v OriginalValueVersionCollection) Len() int {
	return len(v)
}

func (v OriginalValueVersionCollection) Less(i, j int) bool {
	return v[i].Version.LessThan(v[j].Version)
}

func (v OriginalValueVersionCollection) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

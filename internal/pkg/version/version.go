package version

import (
	"fmt"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	Version Info
)

type Info struct {
	Version string    `json:"version,omitempty"`
	Commit  string    `json:"commit,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

func init() {
	Version = Info{
		Version: version,
		Commit:  commit,
	}

	if date == "" {
		return
	}

	var err error
	Version.Date, err = time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
}

func (i Info) String() string {
	return fmt.Sprintf("semver-cli %s, commit %s, built at %s", i.Version, i.Commit, i.Date)
}

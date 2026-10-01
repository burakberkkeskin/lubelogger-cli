package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic release version (e.g. "v0.1.0").
	// It is injected at compile time via -ldflags.
	Version = "0.1.0-dev"

	// Commit is the git commit SHA.
	Commit = "none"

	// Date is the build timestamp.
	Date = "unknown"
)

// Info contains build and runtime version metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	Date      string `json:"date,omitempty"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the current version info.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a human-readable representation of version info.
func (i Info) String() string {
	if i.Commit != "none" && i.Date != "unknown" {
		return fmt.Sprintf("lubelog %s (%s/%s, %s, commit: %s, date: %s)", i.Version, i.OS, i.Arch, i.GoVersion, i.Commit, i.Date)
	}
	return fmt.Sprintf("lubelog %s (%s/%s, %s)", i.Version, i.OS, i.Arch, i.GoVersion)
}

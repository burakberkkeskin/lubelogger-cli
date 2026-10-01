package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	// Version is the semantic release version (e.g. "v0.1.0").
	// It is injected at compile time via -ldflags.
	Version = "dev"

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

// Get returns the current version info, falling back to runtime/debug VCS info if available.
func Get() Info {
	v := Version
	commit := Commit
	date := Date

	if bi, ok := debug.ReadBuildInfo(); ok {
		if v == "dev" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			v = bi.Main.Version
		}
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if commit == "none" {
					commit = s.Value
					if len(commit) > 7 {
						commit = commit[:7]
					}
				}
			case "vcs.time":
				if date == "unknown" {
					date = s.Value
				}
			case "vcs.modified":
				if s.Value == "true" && !strings.HasSuffix(commit, "-dirty") && commit != "none" {
					commit += "-dirty"
				}
			}
		}
	}

	return Info{
		Version:   v,
		Commit:    commit,
		Date:      date,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a human-readable representation of version info.
func (i Info) String() string {
	var meta []string
	if i.Commit != "none" {
		meta = append(meta, fmt.Sprintf("commit: %s", i.Commit))
	}
	if i.Date != "unknown" {
		meta = append(meta, fmt.Sprintf("date: %s", i.Date))
	}
	if len(meta) > 0 {
		return fmt.Sprintf("lubelog %s (%s/%s, %s, %s)", i.Version, i.OS, i.Arch, i.GoVersion, strings.Join(meta, ", "))
	}
	return fmt.Sprintf("lubelog %s (%s/%s, %s)", i.Version, i.OS, i.Arch, i.GoVersion)
}

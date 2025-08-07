package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time via -ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Info returns formatted version information
func Info() string {
	return fmt.Sprintf("ganalyzer version %s\n  Git commit: %s\n  Built: %s\n  Go version: %s\n  OS/Arch: %s/%s",
		Version, GitCommit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// Short returns just the version string
func Short() string {
	return Version
}

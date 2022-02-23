package version

import "runtime"

var (
	// Version holds the complete version number. Filled in at linking time.
	Version = "1.0.0+unknown"

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	// GoVersion is Go tree's version.
	GoVersion = runtime.Version()
)

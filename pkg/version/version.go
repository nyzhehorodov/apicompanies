package version

import (
	"fmt"
	"runtime"
)

var (
	date       string
	commit     string
	appName    string
	appVersion string
)

func Date() string {
	return date
}

func Commit() string {
	return commit
}

func AppName() string {
	return appName
}

func AppVersion() string {
	return appVersion
}

func String() string {
	return fmt.Sprintf(
		"AppName: %q, AppVersion: %q, BuildDate: %q, GitCommit: %q, GoVersion: %q",
		appName,
		appVersion,
		date,
		commit,
		runtime.Version(),
	)
}

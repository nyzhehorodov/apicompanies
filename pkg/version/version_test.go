// nolint:testpackage,paralleltest // testing private properties
package version

import (
	"testing"
)

func TestDate(t *testing.T) {
	date = "testdate"

	got := Date()
	if date != got {
		t.Fatalf(`expected result %q, got %q`, date, got)
	}
}

func TestCommit(t *testing.T) {
	commit = "testcommit"

	got := Commit()
	if commit != got {
		t.Fatalf(`expected result %q, got %q`, commit, got)
	}
}

func TestAppName(t *testing.T) {
	appName = "appname"

	got := AppName()
	if appName != got {
		t.Fatalf(`expected result %q, got %q`, appName, got)
	}
}

func TestAppVersion(t *testing.T) {
	appVersion = "appversion"

	got := AppVersion()
	if appVersion != got {
		t.Fatalf(`expected result %q, got %q`, appVersion, got)
	}
}

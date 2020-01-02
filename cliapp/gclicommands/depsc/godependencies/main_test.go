package godependencies

import (
	"regexp"
	"testing"
)

func TestMatchGoSrcRelativePath(t *testing.T) {
	t.Parallel()
	var (
		result string
		err    error
	)
	if result, err = MatchGoSrcRelativePath("/Users/goatcore/go", "/Users/goatcore/go/src/github.com/goatcms/goatcms"); err != nil {
		t.Error(err)
		return
	}
	if result != "github.com/goatcms/goatcms" {
		t.Errorf("From /Users/goatcore/go/src/github.com/goatcms/goatcms expected github.com/goatcms/goatcms path and take '%v'", result)
	}
	if _, err = MatchGoSrcRelativePath("/Users/goatcore/go", "/no/in/go/path"); err == nil {
		t.Errorf("expected error when CWD is outside GOPATH")
	}
}

func TestIsIgnoredPath(t *testing.T) {
	t.Parallel()
	var (
		AlwaysIgnored = []*regexp.Regexp{
			regexp.MustCompile("^(.*).golang.org$"),
			regexp.MustCompile("^(.*).golang.org/.*$"),
			regexp.MustCompile("^golang.org$"),
			regexp.MustCompile("^golang.org/.*$"),
		}
	)
	if IsIgnoredPath(AlwaysIgnored, "github.com") {
		t.Errorf("github.com should not be ignored")
	}
	if IsIgnoredPath(AlwaysIgnored, "github.com/goatcms/goatcore") {
		t.Errorf("github.com/goatcms/goatcore should not be ignored")
	}
	if !IsIgnoredPath(AlwaysIgnored, "golang.org") {
		t.Errorf("golang.org should be ignored")
	}
	if !IsIgnoredPath(AlwaysIgnored, "golang.org/x/net") {
		t.Errorf("all paths in golang.org should be ignored")
	}
	if !IsIgnoredPath(AlwaysIgnored, "google.golang.org") {
		t.Errorf("google.golang.org should be ignored")
	}
	if !IsIgnoredPath(AlwaysIgnored, "google.golang.org/appengine") {
		t.Errorf("all paths in google.golang.org should be ignored")
	}
}

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

func TestMapPath(t *testing.T) {
	var (
		result string
	)
	t.Parallel()
	result = MapPath(GOPathMapping, "github.com/goatcms/goatcore")
	if result != "github.com/goatcms/goatcore" {
		t.Errorf("MapPath should don't modyfied unmapping paths")
	}
	result = MapPath(GOPathMapping, "google.golang.org/appengine")
	if result != "github.com/golang/appengine" {
		t.Errorf("google.golang.org/appengine should be mapped to github.com/golang/appengine")
	}
	// golang.org/x/blog — the content and server program for blog.golang.org.
	result = MapPath(GOPathMapping, "golang.org/x/blog")
	if result != "github.com/golang/blog" {
		t.Errorf("golang.org/x/blog should be mapped to github.com/golang/blog")
	}
	// golang.org/x/crypto — additional cryptography packages.
	result = MapPath(GOPathMapping, "golang.org/x/crypto")
	if result != "github.com/golang/crypto" {
		t.Errorf("golang.org/x/crypto should be mapped to github.com/golang/crypto")
	}
	// golang.org/x/exp — experimental code (handle with care).
	result = MapPath(GOPathMapping, "golang.org/x/exp")
	if result != "github.com/golang/exp" {
		t.Errorf("golang.org/x/exp should be mapped to github.com/golang/exp")
	}
	// golang.org/x/image — additional imaging packages.
	result = MapPath(GOPathMapping, "golang.org/x/image")
	if result != "github.com/golang/image" {
		t.Errorf("golang.org/x/image should be mapped to github.com/golang/image")
	}
	// golang.org/x/mobile — libraries and build tools for Go on Android.
	result = MapPath(GOPathMapping, "golang.org/x/mobile")
	if result != "github.com/golang/mobile" {
		t.Errorf("golang.org/x/mobile should be mapped to github.com/golang/image")
	}
	// golang.org/x/net — additional networking packages.
	result = MapPath(GOPathMapping, "golang.org/x/net")
	if result != "github.com/golang/net" {
		t.Errorf("golang.org/x/net should be mapped to github.com/golang/net")
	}
	// golang.org/x/sys — for low-level interactions with the operating system.
	result = MapPath(GOPathMapping, "golang.org/x/sys")
	if result != "github.com/golang/sys" {
		t.Errorf("golang.org/x/sys should be mapped to github.com/golang/net")
	}
	// golang.org/x/talks — the content and server program for talks.golang.org.
	result = MapPath(GOPathMapping, "golang.org/x/talks")
	if result != "github.com/golang/talks" {
		t.Errorf("golang.org/x/talks should be mapped to github.com/golang/talks")
	}
	// golang.org/x/text — packages for working with text.
	result = MapPath(GOPathMapping, "golang.org/x/text")
	if result != "github.com/golang/text" {
		t.Errorf("golang.org/x/text should be mapped to github.com/golang/text")
	}
	// https://godoc.org/cloud.google.com/go/compute
	result = MapPath(GOPathMapping, "cloud.google.com/go/compute")
	if result != "cloud.google.com/go/compute" {
		t.Errorf("cloud.google.com/go/compute should be mapped to github.com/GoogleCloudPlatform/google-cloud-go/compute")
	}
}

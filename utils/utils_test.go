package utils

import (
	"regexp"
	"testing"
)

func TestMatchString(t *testing.T) {
	var pat = regexp.MustCompile("^un simple try")
	s := "it doesn't match"
	if MatchString(pat, s) {
		t.Errorf("String %s should not match regexp %s ", s, pat.String())
	}
	s = " un simple try"
	if MatchString(pat, s) {
		t.Errorf("String %s should not match regexp %s ", s, pat.String())
	}
	s = "un simple try"
	if !MatchString(pat, s) {
		t.Errorf("String %s should match regexp %s ", s, pat.String())
	}

	s = "un simple try extended"
	if !MatchString(pat, s) {
		t.Errorf("String %s should match regexp %s ", s, pat.String())
	}
}

func TestIsCorrrectURL(t *testing.T) {
	s := "http:/www.example.com"
	if IsCorrrectURL(s) {
		t.Errorf("Url '%s' should not be valid", s)
	}
	s = " http://www.example.com"
	if IsCorrrectURL(s) {
		t.Errorf("Url '%s' should not be valid ", s)
	}
	s = "http://www.example.com"
	if !IsCorrrectURL(s) {
		t.Errorf("Url '%s' ia valid ", s)
	}
	s = "https://www.example.com"
	if !IsCorrrectURL(s) {
		t.Errorf("Url '%s' ia valid", s)
	}
	s = "http://www.example.com  "
	if !IsCorrrectURL(s) {
		t.Errorf("Url '%s' ia valid", s)
	}
}

func TestIsUnderURL(t *testing.T) {
	base := "http:/www.example.com/first"
	url := "https://www.redhat.com"
	if IsUnderBaseUrl(base, url) {
		t.Errorf("Url '%s' is not under base %s", url, base)
	}
	url = "https://www.redhat.com/www.example.com/first"
	if IsUnderBaseUrl(base, url) {
		t.Errorf("Url '%s' is not under base %s", url, base)
	}
	url = "http:/www.example.com/second"
	if IsUnderBaseUrl(base, url) {
		t.Errorf("Url '%s' is not under base %s", url, base)
	}
	url = "http:/www.example.com/first/seconf"
	if !IsUnderBaseUrl(base, url) {
		t.Errorf("Url '%s' is under base %s", url, base)
	}
}

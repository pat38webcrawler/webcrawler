package utils

import (
	"net/url"
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

func TestGetRootURL(t *testing.T) {
	url := "https://www.example.com/foo/////"
	rooturl := GetRootURL(url)
	if rooturl != "https://www.example.com" {
		t.Errorf("Root url for '%s' should be  '\"https://www.example.com\" ' and is actually '%s'", url, rooturl)
	}
	url = "https://www.example.com/foo///bar//toto?titi"
	rooturl = GetRootURL(url)
	if rooturl != "https://www.example.com" {
		t.Errorf("Root url for '%s' should be  '\"https://www.example.com\" ' and is actually '%s'", url, rooturl)
	}
	url = "https://www.example.com/foo/////#anchor"
	rooturl = GetRootURL(url)
	if rooturl != "https://www.example.com" {
		t.Errorf("Root url for '%s' should be  '\"https://www.example.com\" ' and is actually '%s'", url, rooturl)
	}
	url = "https://www.example.com/"
	rooturl = GetRootURL(url)
	if rooturl != "https://www.example.com" {
		t.Errorf("Root url for '%s' should be  '\"https://www.example.com\" ' and is actually '%s'", url, rooturl)
	}
	url = "https://www.example.com"
	rooturl = GetRootURL(url)
	if rooturl != "https://www.example.com" {
		t.Errorf("Root url for '%s' should be  '\"https://www.example.com\" ' and is actually '%s'", url, rooturl)
	}
	url = "https:/www.example.com/bar"
	rooturl = GetRootURL(url)
	if rooturl != "https:/www.example.com/bar" {
		t.Errorf("Root url for '%s' should be  '\"https:/www.example.com/bar\" ' and is actually '%s'", url, rooturl)
	}
}

func TestGCleanPath(t *testing.T) {
	baseurl := "https://site.com//first/////"
	base, _ := url.Parse(baseurl)
	CleanPath(base)
	if base.String() != "https://site.com/first" {
		t.Errorf("Url %s was not properly cleaned. Expected 'https://site.com/first', got '%s'", baseurl, base.String())
	}
	baseurl = "http://site.com//first/////"
	base, _ = url.Parse(baseurl)
	CleanPath(base)
	if base.String() != "http://site.com/first" {
		t.Errorf("Url %s was not properly cleaned. Expected 'http://site.com/first', got '%s'", baseurl, base.String())
	}
	baseurl = "http://site.com//first#second"
	base, _ = url.Parse(baseurl)
	CleanPath(base)
	if base.String() != "http://site.com/first#second" {
		t.Errorf("Url %s was not properly cleaned. Expected 'http://site.com/first#second', got '%s'", baseurl, base.String())
	}
	baseurl = "http://site.com///first/second/third?parameter1=value1&parameter2=value2"
	base, _ = url.Parse(baseurl)
	CleanPath(base)
	if base.String() != "http://site.com/first/second/third?parameter1=value1&parameter2=value2" {
		t.Errorf("Url %s was not properly cleaned. Expected 'http://site.com/first/second/third?parameter1=value1&parameter2=value2', got '%s'", baseurl, base.String())
	}
}

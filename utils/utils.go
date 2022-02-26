package utils

import (
	"net/http"
	"regexp"
)

var TextUrlPattern = regexp.MustCompile(`text/html.*`)
var URLPattern = regexp.MustCompile("^https?://.*")

// return true if s contains any match of re regex

func MatchString(re *regexp.Regexp, s string) bool {
	return re.MatchString(s)
}

// return true if content-type associated to resp is of type text/html

func IsTextURL(resp *http.Response) bool {
	return MatchString(TextUrlPattern, resp.Header.Get("Content-Type"))
}

// check if provided Url start with http[s]://

func IsCorrrectURL(url string) bool {
	return MatchString(URLPattern, url)
}

// check if provided url is uner base url

func IsUnderBaseUrl(base string, url string) bool {
	basepattern := regexp.MustCompile("^" + base + "/?.*")
	return MatchString(basepattern, url)
}

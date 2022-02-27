package utils

import (
	"net/http"
	"regexp"
)

var TextUrlPattern = regexp.MustCompile(`text/html.*`)
var URLPattern = regexp.MustCompile("^https?://.*")

// MatchString return true if s contains any match of re regex
func MatchString(re *regexp.Regexp, s string) bool {
	return re.MatchString(s)
}

// IsTextURL return true if content-type associated to resp is of type text/html
func IsTextURL(resp *http.Response) bool {
	return MatchString(TextUrlPattern, resp.Header.Get("Content-Type"))
}

// IsCorrrectURL check if provided Url start with http[s]://
func IsCorrrectURL(url string) bool {
	return MatchString(URLPattern, url)
}

// IsUnderBaseUrl check if provided url is uner base url
func IsUnderBaseUrl(base string, url string) bool {
	basepattern := regexp.MustCompile("^" + base + "/?.*")
	return MatchString(basepattern, url)
}

// GetRootURL get  root url for url (e.g. providing https://www.foo.com/bar will return https://www.foo.bar)
func GetRootURL(url string) string {
	pattern := regexp.MustCompile("^(https?://[^/]*).*")
	rootUrl := pattern.ReplaceAllString(url, "$1")
	return (rootUrl)
}

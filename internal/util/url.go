package util

import (
	"net/url"
	"regexp"
)

// MustURL parses rawurl into a URL structure.
//
// This panics if the error occurs.
func MustURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	regexp.MustCompile("")
	return u
}

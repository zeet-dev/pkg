package utils

import "strings"

// TrimURLScheme removes https:// or http:// scheme from the start of a url string
func TrimURLScheme(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return url
}

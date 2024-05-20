package utils

import "regexp"

var (
	emailRegexp *regexp.Regexp
)

func init() {
	emailRegexp = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)
}

// IsEmailValid checks if the email provided passes the required structure and length. https://golangcode.com/validate-an-email-address/
func IsEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegexp.MatchString(e)
}

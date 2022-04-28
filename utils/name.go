package utils

import (
	"regexp"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	dns1035SpecialChars *regexp.Regexp
	emailRegexp         *regexp.Regexp
)

func init() {
	dns1035SpecialChars = regexp.MustCompile("[^a-zA-Z0-9]+")
	emailRegexp = regexp.MustCompile("(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,}$")
}

func DNS1035Name(input string) string {
	return DNS1035WithPrefix(input, "app-")
}

func DNS1035WithPrefix(input string, prefix string) string {
	output := dns1035SpecialChars.ReplaceAllString(input, "-")
	output = strings.ToLower(output)
	output = strings.Trim(output, "-")
	if len(output) > 0 && unicode.IsDigit(rune(output[0])) {
		output = prefix + output
	}
	if len(output) > validation.DNS1035LabelMaxLength {
		output = output[:validation.DNS1035LabelMaxLength]
	}
	return output
}

// IsEmailValid checks if the email provided passes the required structure and length. https://golangcode.com/validate-an-email-address/
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegexp.MatchString(e)
}

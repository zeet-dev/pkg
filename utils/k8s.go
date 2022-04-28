package utils

import (
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	labelSspecialChars *regexp.Regexp
)

func init() {
	labelSspecialChars = regexp.MustCompile("[^a-zA-Z0-9]+")
}

func ParseMemoryQuantityToGB(s string) (amount float32, err error) {
	q, err := resource.ParseQuantity(s)
	if err != nil {
		return 0, err
	}

	return float32(q.ScaledValue(resource.Mega)) / 1000, nil
}

func K8SLabel(input string) string {
	output := labelSspecialChars.ReplaceAllString(input, "-")
	output = strings.Trim(output, "-")
	if len(output) > validation.LabelValueMaxLength {
		output = output[:validation.LabelValueMaxLength]
	}
	return output
}

package utils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeet-dev/pkg/utils"
)

func TestIsEmailValid(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"example@example.com", true},                   // Valid email
		{"ex", false},                                   // Too short
		{strings.Repeat("a", 246) + "@test.com", false}, // Too long
		{strings.Repeat("b", 245) + "@test.com", true},
		{"invalid-email", false}, // Invalid format
	}

	for _, test := range tests {
		result := utils.IsEmailValid(test.email)
		require.Equal(t, result, test.expected, "Expected %s to be %v", test.email, test.expected)
	}
}

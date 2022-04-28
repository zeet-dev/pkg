package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeet-dev/pkg/utils"
)

func TestName(t *testing.T) {
	var tests = []struct {
		input  string
		expect string
	}{
		{"laksjfl;df", "laksjfl-df"},
		{"asd/faj-__w/1o4(*&(H#@WIFS___", "asd-faj-w-1o4-h-wifs"},
		{"******", ""},
		{"29037492374", "app-29037492374"},
		{"[a]b[a-1239.cXm../././\\", "a-b-a-1239-cxm"},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}

	for _, tt := range tests {
		t.Run("Test name generator", func(t *testing.T) {
			out := utils.DNS1035Name(tt.input)
			assert.Equal(t, tt.expect, out)
		})
	}
}

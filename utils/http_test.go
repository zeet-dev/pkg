package utils_test

import (
	"fmt"
	"testing"

	"github.com/zeet-dev/pkg/utils"
)

func TestHttps(t *testing.T) {
	var tests = []string{
		"https://lol.com",
		"//lol.com",
		"http://lol.com",
		"asdf",
	}

	for _, tt := range tests {
		t.Run("Test name", func(t *testing.T) {
			fmt.Println(utils.MakeHTTPS(tt))
		})
	}
}

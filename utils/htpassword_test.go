package utils_test

import (
	"fmt"
	"testing"

	"github.com/zeet-dev/pkg/utils"
)

func TestPassord(t *testing.T) {
	t.Run("Test password", func(t *testing.T) {
		fmt.Println(utils.Htpasswd("user", "pass"))
		fmt.Println(utils.Htpasswd("", "pass"))
		fmt.Println(utils.Htpasswd("user", ""))
		fmt.Println(utils.Htpasswd("", ""))
	})
}

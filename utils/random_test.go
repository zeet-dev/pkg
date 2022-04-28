package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeet-dev/pkg/utils"
)

func TestRandom(t *testing.T) {
	t.Run("Test random generator", func(t *testing.T) {
		token, err := utils.GenerateRandomString(48)
		fmt.Println("rand", token)
		require.NoError(t, err)
		require.Equal(t, len(token), 64)
	})
}

package test_utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func RequireEqualCmp(t *testing.T, expected, actual interface{}, opts ...cmp.Option) {
	if eq := cmp.Equal(expected, actual, opts...); !eq {
		require.True(t, eq, "diff: %s", cmp.Diff(expected, actual, opts...))
	}
}

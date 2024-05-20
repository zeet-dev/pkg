package test_utils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var (
	MatchAnyContext = mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})
)

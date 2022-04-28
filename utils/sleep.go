package utils

import (
	"context"
	"time"
)

func DoneOrSleep(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(duration):
		return false
	}
}

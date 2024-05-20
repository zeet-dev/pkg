package utils

import (
	"github.com/jinzhu/copier"

	"github.com/zeet-dev/pkg/utils/options"
)

// ensure from and to are the same type
func DeepCopy[T any](to T, from T, opts ...options.MustOption[copier.Option]) error {
	option := options.MustNewWithDefaults(copier.Option{
		DeepCopy: true,
	}, opts...)
	return copier.CopyWithOption(to, from, option)
}

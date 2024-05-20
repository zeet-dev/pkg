package options

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// Option represents the idea of functional param constructor not "optional entities"

type Option[T any] func(o *T) error

func (o Option[T]) ToMustOption() MustOption[T] {
	return func(opts *T) {
		err := o(opts)
		if err != nil {
			panic(err)
		}
	}
}

type MustOption[T any] func(o *T)

func (o MustOption[T]) ToOption() Option[T] {
	return func(opts *T) error {
		o(opts)
		return nil
	}
}

func (o MustOption[T]) MustAssign(t *T) {
	MustAssign(t, o)
}

func CompositeOption[T any](options ...Option[T]) Option[T] {
	return func(o *T) error {
		return Assign[T](o, options...)
	}
}

func LazyNew[T any](options ...Option[T]) func() (T, error) {
	return func() (T, error) {
		return New(options...)
	}
}

func Assign[T any](defaults *T, options ...Option[T]) error {
	if defaults == nil {
		return errors.New("cannot assign to a nil pointer")
	}

	var errs *multierror.Error
	for _, o := range options {
		err := o(defaults)
		errs = multierror.Append(errs, err)
	}

	if errs.ErrorOrNil() != nil {
		return errors.WithStack(errs)
	}

	return nil
}

func MustOptionsToOptions[T any](options ...MustOption[T]) []Option[T] {
	return lo.Map(options, func(option MustOption[T], _ int) Option[T] {
		return func(o *T) error {
			option(o)
			return nil
		}
	})
}

func MustAssign[T any](defaults *T, options ...MustOption[T]) {
	err := Assign(defaults, MustOptionsToOptions(options...)...)
	if err != nil {
		panic(err)
	}
}

func NewWithDefaults[T any](defaults T, options ...Option[T]) (T, error) {
	return defaults, Assign(&defaults, options...)
}

func MustNewWithDefaults[T any](defaults T, options ...MustOption[T]) T {
	MustAssign(&defaults, options...)
	return defaults
}

func New[T any](options ...Option[T]) (T, error) {
	var obj T
	return NewWithDefaults(obj, options...)
}

func MustNew[T any](options ...MustOption[T]) T {
	var obj T
	return MustNewWithDefaults(obj, options...)
}

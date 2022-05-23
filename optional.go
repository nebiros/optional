package optional

import (
	"errors"
)

var (
	ErrNoValuePresent = errors.New("no value present")
)

type Optional[T any] interface {
	Get() (T, error)
	IsPresent() bool
	IsEmpty() bool
	Or(f func(v T) Optional[T]) Optional[T]
	OrElse(other T) T
	OrElseErr() (T, error)
	Filter(f func(v T) bool) Optional[T]
}

type optional[T any] struct {
	value   T
	present *struct{}
}

func New[T any](value T) Optional[T] {
	return &optional[T]{value: value, present: &struct{}{}}
}

func Empty[T any]() Optional[T] {
	var empty T
	return &optional[T]{value: empty, present: nil}
}

func OfNullable[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}

	return New[T](*value)
}

func (o *optional[T]) Get() (T, error) {
	if o.present != nil {
		return o.value, nil
	}

	return *new(T), ErrNoValuePresent
}

func (o *optional[T]) IsPresent() bool {
	return o.present != nil
}

func (o *optional[T]) IsEmpty() bool {
	return o.present == nil
}

func (o *optional[T]) Or(f func(v T) Optional[T]) Optional[T] {
	if o.IsPresent() {
		return o
	}

	return f(o.value)
}

func (o *optional[T]) OrElse(other T) T {
	if o.present == nil {
		return other
	}

	return o.value
}

func (o *optional[T]) OrElseErr() (T, error) {
	if o.present == nil {
		return *new(T), ErrNoValuePresent
	}

	return o.value, nil
}

func (o *optional[T]) Filter(f func(v T) bool) Optional[T] {
	if !o.IsPresent() {
		return Empty[T]()
	}

	if f(o.value) {
		return o
	}

	return Empty[T]()
}

func Map[T, U any](o Optional[T], f func(v T) U) Optional[U] {
	if !o.IsPresent() {
		return Empty[U]()
	}

	ov, err := o.Get()
	if err != nil {
		return Empty[U]()
	}

	return New[U](f(ov))
}

func FlatMap[T, U any](o Optional[T], f func(v T) Optional[U]) Optional[U] {
	if !o.IsPresent() {
		return Empty[U]()
	}

	ov, err := o.Get()
	if err != nil {
		return Empty[U]()
	}

	return f(ov)
}

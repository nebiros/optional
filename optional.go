package optional

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNoValuePresent = errors.New("optional: no value present")
)

type Optional[T any] struct {
	value   T
	present *struct{}
}

func New[T any](value T) *Optional[T] {
	return &Optional[T]{value: value, present: &struct{}{}}
}

func Empty[T any]() *Optional[T] {
	var empty T
	return &Optional[T]{value: empty, present: nil}
}

func OfNillable[T any](value *T) *Optional[T] {
	if value == nil {
		return Empty[T]()
	}

	return New[T](*value)
}

func (o *Optional[T]) Get() (T, error) {
	if o.present != nil {
		return o.value, nil
	}

	return *new(T), ErrNoValuePresent
}

func (o *Optional[T]) MustGet() T {
	if o.present != nil {
		return o.value
	}

	panic(ErrNoValuePresent)
}

func (o *Optional[T]) IsPresent() bool {
	return o.present != nil
}

func (o *Optional[T]) IsEmpty() bool {
	return o.present == nil
}

func (o *Optional[T]) Or(f func(v T) *Optional[T]) *Optional[T] {
	if o.IsPresent() {
		return o
	}

	return f(o.value)
}

func (o *Optional[T]) OrElse(other T) T {
	if o.present == nil {
		return other
	}

	return o.value
}

func (o *Optional[T]) OrElseErr() (T, error) {
	if o.present == nil {
		return *new(T), ErrNoValuePresent
	}

	return o.value, nil
}

func (o *Optional[T]) Filter(f func(v T) bool) *Optional[T] {
	if !o.IsPresent() {
		return Empty[T]()
	}

	if f(o.value) {
		return o
	}

	return Empty[T]()
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.IsPresent() {
		return json.Marshal(nil)
	}

	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	var v *T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	if v != nil {
		o.value = *v
		o.present = &struct{}{}

		return nil
	}

	o.present = nil

	return nil
}

func (o *Optional[T]) Scan(src any) error {
	if src == nil {
		o.present = nil

		return nil
	}

	var ok bool

	o.value, ok = src.(T)
	if !ok {
		return fmt.Errorf("optional: failed to scan a '%v' into an optional", src)
	}

	o.present = &struct{}{}

	return nil
}

func (o *Optional[T]) Value() (driver.Value, error) {
	if !o.IsPresent() {
		return nil, nil
	}

	return o.value, nil
}

func Map[T, U any](o Optional[T], f func(v T) U) *Optional[U] {
	if !o.IsPresent() {
		return Empty[U]()
	}

	ov, err := o.Get()
	if err != nil {
		return Empty[U]()
	}

	return New[U](f(ov))
}

func FlatMap[T, U any](o Optional[T], f func(v T) *Optional[U]) *Optional[U] {
	if !o.IsPresent() {
		return Empty[U]()
	}

	ov, err := o.Get()
	if err != nil {
		return Empty[U]()
	}

	return f(ov)
}

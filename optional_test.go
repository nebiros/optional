package optional

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	want := "a value"

	ov := New(want)

	got, err := ov.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", ov, want)
	}
}

func TestEmpty(t *testing.T) {
	ov := Empty[string]()

	if !ov.IsEmpty() {
		t.Errorf("Empty() = %v, want %v", ov, ov.IsEmpty())
	}
}

func TestOfNullable_Nil(t *testing.T) {
	var want *string

	ov := OfNullable[string](want)

	if !ov.IsEmpty() {
		t.Errorf("OfNullable() = %v, want %v", ov, ov.IsEmpty())
	}
}

func TestOfNullable_Pointer(t *testing.T) {
	tmp := "a value"
	want := &tmp

	ov := OfNullable[string](want)

	if !ov.IsPresent() {
		t.Errorf("OfNullable() = %v, want %v", ov, ov.IsPresent())
	}

	if ov.IsEmpty() {
		t.Errorf("OfNullable() = %v, want %v", ov, ov.IsEmpty())
	}
}

func TestGet(t *testing.T) {
	want := "a value"

	ov := New(want)

	got, err := ov.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", ov, want)
	}
}

func TestGet_ErrNoValuePresent(t *testing.T) {
	ov := Empty[string]()

	_, err := ov.Get()
	if !errors.Is(err, ErrNoValuePresent) {
		t.Errorf("Empty() = %v, want %v", ov, ErrNoValuePresent)
	}
}

func TestOrElse(t *testing.T) {
	var (
		wantNil    *string
		wantNonNil = "a non nil value"
	)

	ov := OfNullable[string](wantNil)

	got := ov.OrElse(wantNonNil)

	if !reflect.DeepEqual(got, wantNonNil) {
		t.Errorf("OfNullable() = %v, want %v", ov, wantNonNil)
	}
}

func TestOrElseErr(t *testing.T) {
	var want *string

	ov := OfNullable[string](want)

	_, err := ov.OrElseErr()
	if !errors.Is(err, ErrNoValuePresent) {
		t.Errorf("OfNullable() = %v, want %v", ov, ErrNoValuePresent)
	}
}

func TestFilter(t *testing.T) {
	want := 2022

	ov := New(want)

	ov = ov.Filter(func(v int) bool {
		return v == want
	})

	if !ov.IsPresent() {
		t.Errorf("New() = %v, want %v", ov, ov.IsPresent())
	}

	if ov.IsEmpty() {
		t.Errorf("New() = %v, want %v", ov, ov.IsEmpty())
	}

	got, err := ov.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", ov, want)
	}
}

func TestMap(t *testing.T) {
	var (
		value = 2022
		want  = "2022"
	)

	ov := New(value)

	m := Map(ov, func(v int) string {
		return fmt.Sprintf("%d", v)
	})

	if !m.IsPresent() {
		t.Errorf("Map() = %v, want %v", m, m.IsPresent())
	}

	if ov.IsEmpty() {
		t.Errorf("Map() = %v, want %v", m, m.IsEmpty())
	}

	got, err := m.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", m, want)
	}
}

func TestFlatMap(t *testing.T) {
	var (
		value = 2022
		want  = "2022"
	)

	ov := New(value)

	m := FlatMap(ov, func(v int) Optional[string] {
		return New(fmt.Sprintf("%d", v))
	})

	if !m.IsPresent() {
		t.Errorf("FlatMap() = %v, want %v", m, m.IsPresent())
	}

	if ov.IsEmpty() {
		t.Errorf("FlatMap() = %v, want %v", m, m.IsEmpty())
	}

	got, err := m.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FlatMap() = %v, want %v", m, want)
	}
}

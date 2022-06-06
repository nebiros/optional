package optional

import (
	"bytes"
	"encoding/json"
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
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestEmpty(t *testing.T) {
	ov := Empty[string]()

	if !ov.IsEmpty() {
		t.Errorf("Empty() = %v, want %v", ov, ov.IsEmpty())
	}
}

func TestOfNillable_Nil(t *testing.T) {
	var want *string

	ov := OfNillable[string](want)

	if !ov.IsEmpty() {
		t.Errorf("OfNillable() = %v, want %v", ov, ov.IsEmpty())
	}
}

func TestOfNillable_Pointer(t *testing.T) {
	tmp := "a value"
	want := &tmp

	ov := OfNillable[string](want)

	if !ov.IsPresent() {
		t.Errorf("OfNillable() = %v, want %v", ov, ov.IsPresent())
	}

	if ov.IsEmpty() {
		t.Errorf("OfNillable() = %v, want %v", ov, ov.IsEmpty())
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
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestGet_ErrNoValuePresent(t *testing.T) {
	ov := Empty[string]()

	_, err := ov.Get()
	if !errors.Is(err, ErrNoValuePresent) {
		t.Errorf("Empty() = %v, want %v", ov, ErrNoValuePresent)
	}
}

func TestMustGet(t *testing.T) {
	want := "a value"

	ov := New(want)

	got := ov.MustGet()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestMustGet_ErrNoValuePresent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustGet() did not panic")
		}
	}()

	ov := Empty[string]()

	_ = ov.MustGet()
}

func TestOrElse(t *testing.T) {
	var (
		wantNil    *string
		wantNonNil = "a non nil value"
	)

	ov := OfNillable[string](wantNil)

	got := ov.OrElse(wantNonNil)

	if !reflect.DeepEqual(got, wantNonNil) {
		t.Errorf("got = %v, want %v", got, wantNonNil)
	}
}

func TestOrElseErr(t *testing.T) {
	var want *string

	ov := OfNillable[string](want)

	_, err := ov.OrElseErr()
	if !errors.Is(err, ErrNoValuePresent) {
		t.Errorf("OfNillable() = %v, want %v", ov, ErrNoValuePresent)
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
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestMarshalJSON(t *testing.T) {
	want := []byte(`{"aField":"a value"}`)

	tmp := "a value"
	aValue := &tmp

	type aTestType struct {
		AField Optional[string] `json:"aField"`
	}

	v := aTestType{
		AField: OfNillable(aValue),
	}

	got, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestMarshalJSON_Nil(t *testing.T) {
	want := []byte(`{"aField":null}`)

	type aTestType struct {
		AField Optional[string] `json:"aField"`
	}

	v := aTestType{}

	got, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	want := "a value"
	data := []byte(`{"aField":"a value"}`)

	type aTestType struct {
		AField Optional[string] `json:"aField"`
	}

	var v *aTestType

	err := json.Unmarshal(data, &v)
	if err != nil {
		t.Error(err)
	}

	if !v.AField.IsPresent() {
		t.Errorf("v.AField = %v, want %v", v.AField, v.AField.IsPresent())
	}

	if v.AField.IsEmpty() {
		t.Errorf("v.AField = %v, want %v", v.AField, v.AField.IsEmpty())
	}

	got, err := v.AField.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
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
		t.Errorf("got = %v, want %v", got, want)
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
		t.Errorf("got = %v, want %v", got, want)
	}
}

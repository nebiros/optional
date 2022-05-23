package optional_test

import (
	"fmt"

	"github.com/nebiros/optional"
)

func Example() {
	ov := optional.New("something")

	sv, err := ov.Get()
	if err != nil {
		panic(err)
	}

	fmt.Println("sv: " + sv)
}

func Example_nullable() {
	var (
		v   *string
		tmp = "something"
	)

	v = &tmp

	ov := optional.OfNullable(v)

	if !ov.IsPresent() {
		panic("v not present")
	}

	sv, err := ov.Get()
	if err != nil {
		panic(err)
	}

	fmt.Println("sv: " + sv)
}

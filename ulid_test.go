package ulidgo_test

import (
	"fmt"

	"github.com/matsuyoshi30/ulidgo"
)

func init() {
	ulidgo.SetNow()
	ulidgo.SetSeed()
}

func ExampleNew() {
	fmt.Println(ulidgo.New())
	// Output: 01G00RPN3GXT7N2G5ZR6AW6TR5
}

func ExampleTime() {
	fmt.Println(ulidgo.New().Time())
	// Output: 2022-04-07 01:02:30 +0000 UTC
}

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

package ulidgo_test

import (
	"fmt"
	"time"

	"github.com/matsuyoshi30/ulidgo"
)

func init() {
	ulidgo.SetSeed()
}

var Now = func() time.Time { return time.Date(2022, time.April, 7, 1, 2, 30, 45, time.UTC) }

func ExampleNew() {
	fmt.Println(ulidgo.New(Now().UnixMilli()))
	// Output: 01G00RPN3GXT7N2G5ZR6AW6TR5
}

func ExampleULID_Time() {
	fmt.Println(ulidgo.New(Now().UnixMilli()).Time())
	// Output: 2022-04-07 01:02:30 +0000 UTC
}

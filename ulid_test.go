package ulidgo_test

import (
	"fmt"
	"time"

	"github.com/matsuyoshi30/ulidgo"
)

func ExampleNew() {
	ulidgo.Now = func() time.Time { return time.Date(2022, time.April, 7, 1, 2, 30, 45, time.UTC) }
	ulidgo.Seed = int64(1234567890)
	fmt.Println(ulidgo.New())
	// Output: 060032TME0XT7N2G5ZR6AW6TR5
}

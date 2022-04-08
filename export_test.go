package ulidgo

import "time"

func SetNow() {
	now = func() time.Time { return time.Date(2022, time.April, 7, 1, 2, 30, 45, time.UTC) }
}

func SetSeed() {
	seed = int64(1234567890)
}

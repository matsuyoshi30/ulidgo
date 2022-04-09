package ulidgo

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      32_bit_uint_time_high                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     16_bit_uint_time_low      |       16_bit_uint_random      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       32_bit_uint_random                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       32_bit_uint_random                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// ULID represents ULID
type ULID [16]byte

// New ...
func New() *ULID {
	var ulid ULID

	ulid.setTimestamp()
	if err := ulid.setRandom(); err != nil {
		panic(err)
	}

	return &ulid
}

var now = time.Now

func (u *ULID) setTimestamp() {
	n := now().UnixMilli()
	u[0] = byte(n >> 40)
	u[1] = byte(n >> 32)
	u[2] = byte(n >> 24)
	u[3] = byte(n >> 16)
	u[4] = byte(n >> 8)
	u[5] = byte(n)
}

var seed = now().UnixMilli()

func (u *ULID) setRandom() error {
	_, err := rand.New(rand.NewSource(seed)).Read(u[6:])
	return err
}

// String implements fmt.Stringer
func (u *ULID) String() string {
	var dst strings.Builder

	// encode timestamp field
	var ts strings.Builder
	for _, b := range u[:6] {
		fmt.Fprintf(&ts, "%08b", b)
	}
	fmt.Fprint(&dst, mustCB32(ts.String()[0:3]))
	for i := 3; i < 48; i += 5 {
		fmt.Fprint(&dst, mustCB32(ts.String()[i:i+5]))
	}

	// encode random field
	var r strings.Builder
	for _, b := range u[6:] {
		fmt.Fprintf(&r, "%08b", b)
	}
	for i := 0; i < 80; i += 5 {
		fmt.Fprint(&dst, mustCB32(r.String()[i:i+5]))
	}

	return dst.String()
}

func mustCB32(s string) string {
	s, err := cb32(s)
	if err != nil {
		panic(err)
	}
	return s
}

const cbs = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func cb32(s string) (string, error) {
	num, err := strconv.ParseInt(s, 2, 8)
	if err != nil {
		return "", err
	}
	if num-1 > int64(len(cbs)) {
		return "", fmt.Errorf("unexpected: num=%d, s=%s", num, s)
	}
	return string(cbs[num]), nil
}

// UnixTime returns ULID unix timestamp value
func (u *ULID) UnixTime() int64 {
	return int64(u[0])<<40 | int64(u[1])<<32 | int64(u[2])<<24 | int64(u[3])<<16 | int64(u[4])<<8 | int64(u[5])
}

// Time returns UTC converted from ULID unix timestamp value
func (u *ULID) Time() time.Time {
	return time.Unix(u.UnixTime()/1000, 0).UTC()
}

// TODO: Monotonicity

// TODO: compare

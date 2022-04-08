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
	// time.Time.UnixMilli returns int64 (8bytes)
	// 11111111 22222222 33333333 44444444 55555555 66666666 77777777 88888888
	// 00000000 00000000 00000000 00000000 00000000 11111111 22222222 33333333 (>> 40) => byte() is 33333333
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
	// TODO: improve
	fmt.Fprint(&dst, mustCB32(ts.String()[0:3]))
	fmt.Fprint(&dst, mustCB32(ts.String()[3:8]))
	fmt.Fprint(&dst, mustCB32(ts.String()[8:13]))
	fmt.Fprint(&dst, mustCB32(ts.String()[13:18]))
	fmt.Fprint(&dst, mustCB32(ts.String()[18:23]))
	fmt.Fprint(&dst, mustCB32(ts.String()[23:28]))
	fmt.Fprint(&dst, mustCB32(ts.String()[28:33]))
	fmt.Fprint(&dst, mustCB32(ts.String()[33:38]))
	fmt.Fprint(&dst, mustCB32(ts.String()[38:43]))
	fmt.Fprint(&dst, mustCB32(ts.String()[43:48]))

	// encode random field
	var r strings.Builder
	for _, b := range u[6:] {
		fmt.Fprintf(&r, "%08b", b)
	}
	// TODO: improve
	fmt.Fprint(&dst, mustCB32(r.String()[0:5]))
	fmt.Fprint(&dst, mustCB32(r.String()[5:10]))
	fmt.Fprint(&dst, mustCB32(r.String()[10:15]))
	fmt.Fprint(&dst, mustCB32(r.String()[15:20]))
	fmt.Fprint(&dst, mustCB32(r.String()[20:25]))
	fmt.Fprint(&dst, mustCB32(r.String()[25:30]))
	fmt.Fprint(&dst, mustCB32(r.String()[30:35]))
	fmt.Fprint(&dst, mustCB32(r.String()[35:40]))
	fmt.Fprint(&dst, mustCB32(r.String()[40:45]))
	fmt.Fprint(&dst, mustCB32(r.String()[45:50]))
	fmt.Fprint(&dst, mustCB32(r.String()[50:55]))
	fmt.Fprint(&dst, mustCB32(r.String()[55:60]))
	fmt.Fprint(&dst, mustCB32(r.String()[60:65]))
	fmt.Fprint(&dst, mustCB32(r.String()[65:70]))
	fmt.Fprint(&dst, mustCB32(r.String()[70:75]))
	fmt.Fprint(&dst, mustCB32(r.String()[75:80]))

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

// TODO: Monotonicity

// TODO: compare

package ulidgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
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
type ULID struct {
	b [16]byte
	e [26]byte
}

// New generates new ULID
func New(ts int64) (*ULID, error) {
	var ulid ULID

	if err := ulid.setTimestamp(ts); err != nil {
		return nil, err
	}
	if err := ulid.setRandom(ts); err != nil {
		return nil, err
	}

	return &ulid, nil
}

func (u *ULID) setTimestamp(ts int64) error {
	if ts > maxTime {
		return fmt.Errorf("invalid timestamp value")
	}

	u.b[0] = byte(ts >> 40)
	u.b[1] = byte(ts >> 32)
	u.b[2] = byte(ts >> 24)
	u.b[3] = byte(ts >> 16)
	u.b[4] = byte(ts >> 8)
	u.b[5] = byte(ts)

	return nil
}

var (
	mu       sync.Mutex
	lasttime int64
	lastrand [10]byte
)

func (u *ULID) setRandom(ts int64) error {
	if lasttime != ts {
		_, err := rand.New(rand.NewSource(ts)).Read(u.b[6:])
		if err != nil {
			return err
		}
		mu.Lock()
		lasttime = ts
		copy(lastrand[:], u.b[6:])
		mu.Unlock()
		return nil
	}

	// increment
	var (
		lo uint64
		hi uint16
	)
	err := binary.Read(bytes.NewReader(lastrand[2:]), binary.BigEndian, &lo)
	if err != nil {
		return err
	}
	err = binary.Read(bytes.NewReader(lastrand[0:2]), binary.BigEndian, &hi)
	if err != nil {
		return err
	}

	lb, hb := new(bytes.Buffer), new(bytes.Buffer)
	if lo+1 != 0 {
		err = binary.Write(lb, binary.BigEndian, lo+1)
		if err != nil {
			return err
		}
		copy(u.b[8:], lb.Bytes())

		err = binary.Write(hb, binary.BigEndian, hi)
		if err != nil {
			return err
		}
		copy(u.b[6:8], hb.Bytes())

		u.copyLastRand()
		return nil
	}

	if hi+1 == 0 {
		return fmt.Errorf("overflow random value")
	}

	err = binary.Write(lb, binary.BigEndian, int64(0)) // move up
	if err != nil {
		return err
	}
	copy(u.b[8:], lb.Bytes())

	err = binary.Write(hb, binary.BigEndian, hi+1)
	if err != nil {
		return err
	}
	copy(u.b[6:8], hb.Bytes())

	u.copyLastRand()
	return nil
}

func (u *ULID) copyLastRand() {
	mu.Lock()
	copy(lastrand[:], u.b[6:])
	mu.Unlock()
}

// String implements fmt.Stringer
func (u *ULID) String() string {
	u.encode()
	return string(u.e[:])
}

// Bytes returns ULID byte slice
func (u *ULID) Bytes() []byte {
	return u.b[:]
}

func (u *ULID) encode() {
	const cbs = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// encode timestamp field by Crockford's base32
	//
	// AAAAAAAA BBBBBBBB CCCCCCCC DDDDDDDD
	// EEEEEEEE FFFFFFFF
	//
	// The exact bit sequence will be as follows (with trailing padding, which differs from the original Unix timestamp)
	//
	// 00AAA AAAAA BBBBB BBBCC CCCCC CDDDD DDDDE EEEEE EEFFF FFFFF
	//
	// (u.b[0] & 11100000) >> 3                                => 00AAA
	// u.b[0] & 00011111                                       => AAAAA
	// (u.b[1] & 11111000) >> 3                                => BBBBB
	// ((u.b[1] & 00000111) << 2) | ((u.b[2] & 11000000) >> 6) => BBBCC
	// (u.b[2] & 00111110) >> 1                                => CCCCC
	// ((u.b[2] & 00000001) << 4) | ((u.b[3] & 11110000) >> 4) => CDDDD
	// ((u.b[3] & 00001111) << 1) | ((u.b[4] & 10000000) >> 7) => DDDDE
	// (u.b[4] & 01111100) >> 2                                => EEEEE
	// ((u.b[4] & 00000011) << 3) | ((u.b[5] & 11100000) >> 5) => EEFFF
	// u.b[5] & 00011111                                       => FFFFF
	u.e[0] = cbs[(u.b[0]&224)>>3]
	u.e[1] = cbs[u.b[0]&31]
	u.e[2] = cbs[(u.b[1]&248)>>3]
	u.e[3] = cbs[((u.b[1]&7)<<2)|((u.b[2]&192)>>6)]
	u.e[4] = cbs[(u.b[2]&62)>>1]
	u.e[5] = cbs[((u.b[2]&1)<<4)|((u.b[3]&240)>>4)]
	u.e[6] = cbs[((u.b[3]&15)<<1)|((u.b[4]&128)>>7)]
	u.e[7] = cbs[(u.b[4]&124)>>2]
	u.e[8] = cbs[((u.b[4]&3)<<3)|((u.b[5]&224)>>5)]
	u.e[9] = cbs[u.b[5]&31]

	// encode random field by Crockford's base32
	//
	//                   AAAAAAAA BBBBBBBB
	// CCCCCCCC DDDDDDDD EEEEEEEE FFFFFFFF
	// GGGGGGGG HHHHHHHH IIIIIIII JJJJJJJJ
	//
	// AAAAA AAABB BBBBB BCCCC CCCCD DDDDD DDEEE EEEEE
	// FFFFF FFFGG GGGGG GHHHH HHHHI IIIII IIJJJ JJJJJ
	//
	// (u.b[6] & 11111000) >> 3                                   => AAAAA
	// ((u.b[6] & 00000111) << 2) | ((u.b[7] & 11000000) >> 6)    => AAABB
	// (u.b[7] & 00111110) >> 1                                   => BBBBB
	// ((u.b[7] & 00000001) << 4) | ((u.b[8] & 11110000) >> 4)    => BCCCC
	// ((u.b[8] & 00001111) << 1) | ((u.b[9] & 10000000) >> 7)    => CCCCD
	// (u.b[9] & 01111100) >> 2                                   => DDDDD
	// ((u.b[9] & 00000011) << 3) | ((u.b[10] & 11100000) >> 5)   => DDEEE
	// u.b[10] & 00011111                                         => EEEEE
	// (u.b[11] & 11111000) >> 3                                  => FFFFF
	// ((u.b[11] & 00000111) << 2) | ((u.b[12] & 11000000) >> 6)  => FFFGG
	// (u.b[12] & 00111110) >> 1                                  => GGGGG
	// ((u.b[12] & 00000001) << 4) | ((u.b[13] & 11110000) >> 4)  => GHHHH
	// ((u.b[13] & 00001111) << 1) | ((u.b[14] & 10000000) >> 7)  => HHHHI
	// (u.b[14] & 01111100) >> 2                                  => IIIII
	// ((u.b[14] & 00000011) << 3) | ((u.b[15] & 11100000) >> 5)  => IIJJJ
	// u.b[15] & 00011111                                         => JJJJJ
	u.e[10] = cbs[(u.b[6]&248)>>3]
	u.e[11] = cbs[((u.b[6]&7)<<2)|((u.b[7]&192)>>6)]
	u.e[12] = cbs[(u.b[7]&62)>>1]
	u.e[13] = cbs[((u.b[7]&1)<<4)|((u.b[8]&240)>>4)]
	u.e[14] = cbs[((u.b[8]&15)<<1)|((u.b[9]&128)>>7)]
	u.e[15] = cbs[(u.b[9]&124)>>2]
	u.e[16] = cbs[((u.b[9]&3)<<3)|((u.b[10]&224)>>5)]
	u.e[17] = cbs[u.b[10]&31]
	u.e[18] = cbs[(u.b[11]&248)>>3]
	u.e[19] = cbs[((u.b[11]&7)<<2)|((u.b[12]&192)>>6)]
	u.e[20] = cbs[(u.b[12]&62)>>1]
	u.e[21] = cbs[((u.b[12]&1)<<4)|((u.b[13]&240)>>4)]
	u.e[22] = cbs[((u.b[13]&15)<<1)|((u.b[14]&128)>>7)]
	u.e[23] = cbs[(u.b[14]&124)>>2]
	u.e[24] = cbs[((u.b[14]&3)<<3)|((u.b[15]&224)>>5)]
	u.e[25] = cbs[u.b[15]&31]
}

const (
	maxTimeVal = int64(255)
	maxTime    = maxTimeVal<<40 | maxTimeVal<<32 | maxTimeVal<<24 | maxTimeVal<<16 | maxTimeVal<<8 | maxTimeVal
)

// UnixTime returns ULID unix timestamp value
func (u *ULID) UnixTime() int64 {
	return int64(u.b[0])<<40 | int64(u.b[1])<<32 | int64(u.b[2])<<24 | int64(u.b[3])<<16 | int64(u.b[4])<<8 | int64(u.b[5])
}

// Time returns UTC converted from ULID unix timestamp value
func (u *ULID) Time() time.Time {
	return time.Unix(u.UnixTime()/1000, 0).UTC()
}

// Compare returns an integer comparing two ULID byte slice
// The result will be 0 if u == target, -1 if u < target and +1 if u > target
func (u *ULID) Compare(target []byte) int {
	return bytes.Compare(u.b[:], target)
}

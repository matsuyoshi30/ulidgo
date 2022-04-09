package ulidgo

import (
	"math/rand"
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
	const cbs = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	var dst [26]byte

	// encode timestamp field by Crockford's base32
	//
	// AAAAAAAA BBBBBBBB CCCCCCCC DDDDDDDD
	// EEEEEEEE FFFFFFFF
	//
	// The exact bit sequence will be as follows (with trailing padding, which differs from the original Unix timestamp)
	//
	// 00AAA AAAAA BBBBB BBBCC CCCCC CDDDD DDDDE EEEEE EEFFF FFFFF
	//
	// (u[0] & 11100000) >> 3                               => 00AAA
	// u[0] & 00011111                                      => AAAAA
	// (u[1] & 11111000) >> 3                               => BBBBB
	// ((u[1] & 00000111) << 2) | ((u[2] & 11000000) >> 6)  => BBBCC
	// (u[2] & 00111110) >> 1                               => CCCCC
	// ((u[2] & 00000001) << 4) | ((u[3] & 11110000) >> 4)  => CDDDD
	// ((u[3] & 00001111) << 1) | ((u[4] & 10000000) >> 7)  => DDDDE
	// (u[4] & 01111100) >> 2                               => EEEEE
	// ((u[4] & 00000011) << 3) | ((u[5] & 11100000) >> 5)  => EEFFF
	// u[5] & 00011111                                      => FFFFF
	dst[0] = cbs[(u[0]&224)>>3]
	dst[1] = cbs[u[0]&31]
	dst[2] = cbs[(u[1]&248)>>3]
	dst[3] = cbs[((u[1]&7)<<2)|((u[2]&192)>>6)]
	dst[4] = cbs[(u[2]&62)>>1]
	dst[5] = cbs[((u[2]&1)<<4)|((u[3]&240)>>4)]
	dst[6] = cbs[((u[3]&15)<<1)|((u[4]&128)>>7)]
	dst[7] = cbs[(u[4]&124)>>2]
	dst[8] = cbs[((u[4]&3)<<3)|((u[5]&224)>>5)]
	dst[9] = cbs[u[5]&31]

	// encode random field by Crockford's base32
	//
	//                   AAAAAAAA BBBBBBBB
	// CCCCCCCC DDDDDDDD EEEEEEEE FFFFFFFF
	// GGGGGGGG HHHHHHHH IIIIIIII JJJJJJJJ
	//
	// AAAAA AAABB BBBBB BCCCC CCCCD DDDDD DDEEE EEEEE
	// FFFFF FFFGG GGGGG GHHHH HHHHI IIIII IIJJJ JJJJJ
	//
	// (u[6] & 11111000) >> 3                                  => AAAAA
	// ((u[6] & 00000111) << 2) | ((u[7] & 11000000) >> 6)     => AAABB
	// (u[7] & 00111110) >> 1                                  => BBBBB
	// ((u[7] & 00000001) << 4) | ((u[8] & 11110000) >> 4)     => BCCCC
	// ((u[8] & 00001111) << 1) | ((u[9] & 10000000) >> 7)     => CCCCD
	// (u[9] & 01111100) >> 2                                  => DDDDD
	// ((u[9] & 00000011) << 3) | ((u[10] & 11100000) >> 5)    => DDEEE
	// u[10] & 00011111                                        => EEEEE
	// (u[11] & 11111000) >> 3                                 => FFFFF
	// ((u[11] & 00000111) << 2) | ((u[12] & 11000000) >> 6)   => FFFGG
	// (u[12] & 00111110) >> 1                                 => GGGGG
	// ((u[12] & 00000001) << 4) | ((u[13] & 11110000) >> 4)   => GHHHH
	// ((u[13] & 00001111) << 1) | ((u[14] & 10000000) >> 7)   => HHHHI
	// (u[14] & 01111100) >> 2                                 => IIIII
	// ((u[14] & 00000011) << 3) | ((u[15] & 11100000) >> 5)   => IIJJJ
	// u[15] & 00011111                                        => JJJJJ
	dst[10] = cbs[(u[6]&248)>>3]
	dst[11] = cbs[((u[6]&7)<<2)|((u[7]&192)>>6)]
	dst[12] = cbs[(u[7]&62)>>1]
	dst[13] = cbs[((u[7]&1)<<4)|((u[8]&240)>>4)]
	dst[14] = cbs[((u[8]&15)<<1)|((u[9]&128)>>7)]
	dst[15] = cbs[(u[9]&124)>>2]
	dst[16] = cbs[((u[9]&3)<<3)|((u[10]&224)>>5)]
	dst[17] = cbs[u[10]&31]
	dst[18] = cbs[(u[11]&248)>>3]
	dst[19] = cbs[((u[11]&7)<<2)|((u[12]&192)>>6)]
	dst[20] = cbs[(u[12]&62)>>1]
	dst[21] = cbs[((u[12]&1)<<4)|((u[13]&240)>>4)]
	dst[22] = cbs[((u[13]&15)<<1)|((u[14]&128)>>7)]
	dst[23] = cbs[(u[14]&124)>>2]
	dst[24] = cbs[((u[14]&3)<<3)|((u[15]&224)>>5)]
	dst[25] = cbs[u[15]&31]

	return string(dst[:])
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

package ulidgo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/matsuyoshi30/ulidgo"
)

var Now = func() time.Time { return time.Date(2022, time.April, 7, 1, 2, 30, 45, time.UTC) }

func ExampleNew() {
	ulid, _ := ulidgo.New(Now().UnixMilli())
	fmt.Println(ulid)
	// Output: 01G00RPN3GNQDPDEA8MJAJS8SJ
}

func ExampleNew_multi() {
	nt := time.Date(2022, time.April, 8, 1, 2, 30, 45, time.UTC)
	ulid1, _ := ulidgo.New(nt.UnixMilli())
	ulid2, _ := ulidgo.New(nt.UnixMilli())
	fmt.Println(ulid1, ulid2)
	// Output: 01G03B3C3GCZVQY8TYZAEPCGQ5 01G03B3C3GCZVQY8TYZAEPCGQ6
}

func ExampleULID_Time() {
	ulid, _ := ulidgo.New(Now().UnixMilli())
	fmt.Println(ulid.Time())
	// Output: 2022-04-07 01:02:30 +0000 UTC
}

func TestNew(t *testing.T) {
	now := time.Date(2022, time.April, 8, 1, 2, 30, 45, time.UTC)

	ulid, err := ulidgo.New(now.UnixMilli())
	if err != nil {
		t.Error(err)
	}
	ulid2, err := ulidgo.New(now.UnixMilli())
	if err != nil {
		t.Error(err)
	}
	if ulid.Compare(ulid2.Bytes()) != -1 {
		t.Errorf("unexpected result of compare")
	}

	// TODO: overflow error
}

func TestULID_Compare(t *testing.T) {
	u1, err := ulidgo.New(time.Date(2022, time.April, 7, 1, 2, 30, 45, time.UTC).UnixMilli())
	if err != nil {
		t.Error(err)
	}
	u2, err := ulidgo.New(time.Date(2022, time.April, 8, 1, 2, 30, 45, time.UTC).UnixMilli())
	if err != nil {
		t.Error(err)
	}
	if u1.Compare(u1.Bytes()) != 0 {
		t.Errorf("want 0 but got %d\n", u1.Compare(u1.Bytes()))
	}
	if u1.Compare(u2.Bytes()) != -1 {
		t.Errorf("want -1 but got %d\n", u1.Compare(u2.Bytes()))
	}
	if u2.Compare(u1.Bytes()) != 1 {
		t.Errorf("want +1 but got %d\n", u2.Compare(u1.Bytes()))
	}
}

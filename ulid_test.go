package ulidgo_test

import (
	"fmt"
	"sync"
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
	ulid3, _ := ulidgo.New(nt.UnixMilli())
	fmt.Printf("%s\n%s\n%s\n", ulid1, ulid2, ulid3)
	// Output:
	// 01G03B3C3GCZVQY8TYZAEPCGQ5
	// 01G03B3C3GCZVQY8TYZAEPCGQ6
	// 01G03B3C3GCZVQY8TYZAEPCGQ7
}

func ExampleULID_Time() {
	ulid, _ := ulidgo.New(Now().UnixMilli())
	fmt.Println(ulid.Time())
	// Output: 2022-04-07 01:02:30 +0000 UTC
}

func TestNew(t *testing.T) {
	ulid, err := ulidgo.New(Now().UnixMilli())
	if err != nil {
		t.Error(err)
	}
	ulid2, err := ulidgo.New(Now().UnixMilli())
	if err != nil {
		t.Error(err)
	}
	if ulid.Compare(ulid2.Bytes()) != -1 {
		t.Errorf("unexpected result of compare: %d\n", ulid.Compare(ulid2.Bytes()))
	}

	ts := ulidgo.GenMaxRandomValULID()
	_, err = ulidgo.New(ts)
	if err != ulidgo.ErrOverflow {
		t.Errorf("want overflow error")
	}
}

func TestNew_Parallel(t *testing.T) {
	var wg sync.WaitGroup
	got := make(map[string]bool, 0)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ulid, err := ulidgo.New(Now().UnixMilli())
			if err != nil {
				t.Error(err)
			}
			if _, ok := got[ulid.String()]; ok {
				t.Errorf("unexpected result")
			} else {
				got[ulid.String()] = true
			}
		}()
	}
	wg.Wait()
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		ulid string
		want time.Time
		err  error
	}{
		{
			name: "normal",
			ulid: "01G00RPN3GNQDPDEA8MJAJS8SJ",
			want: time.Date(2022, time.April, 7, 1, 2, 30, 0, time.UTC),
		},
		{
			name: "invalid length",
			ulid: "01G00RPN3GNQDPDEA8MJA", // timestamp field is '01G00RPN3G'
			err:  ulidgo.ErrInvalidULIDLen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ulidgo.Parse(tt.ulid)
			if tt.err == nil && err != nil {
				t.Errorf("unexpected error: %q", err)
			} else if tt.err != nil && err == nil {
				t.Errorf("expected error %q but got nil", tt.err)
			}
			if !tt.want.Equal(got) {
				t.Errorf("want %q but got %q", tt.want, got)
			}
		})
	}
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

func BenchmarkFactory(b *testing.B) {
	ulid := ulidgo.Factory()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ulid()
	}
}

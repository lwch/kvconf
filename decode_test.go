package kvconf

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestDecodeInterface(t *testing.T) {
	str :=
		`a = b
c = d`
	var intf interface{}
	err := NewDecoder(strings.NewReader(str)).Decode(&intf)
	if err != nil {
		t.Fatalf("decode interface failed: %v", err)
	}
	mp, ok := intf.(map[string]string)
	if !ok {
		t.Fatalf("unexpected type of interface, want map[string]string current %s", reflect.TypeOf(intf).String())
	}
	if mp["a"] != "b" {
		t.Fatalf("unexpected map value of a, want b current %s", mp["a"])
	}
	if mp["c"] != "d" {
		t.Fatalf("unexpected map value of c, want d current %s", mp["c"])
	}
}

func TestDecodeMap(t *testing.T) {
	str :=
		`a = b
c = d`
	var mp map[string]string
	err := NewDecoder(strings.NewReader(str)).Decode(&mp)
	if err != nil {
		t.Fatalf("decode map failed: %v", err)
	}
	if mp["a"] != "b" {
		t.Fatalf("unexpected map value of a, want b current %s", mp["a"])
	}
	if mp["c"] != "d" {
		t.Fatalf("unexpected map value of c, want d current %s", mp["c"])
	}
}

func TestDecodeStruct(t *testing.T) {
	str :=
		`a = 1
b = 2
c = true
d = 3.14
e = abc`
	var st1 struct {
		A int     `kv:"a"`
		B uint    `kv:"b"`
		C bool    `kv:"c"`
		D float64 `kv:"d"`
		E string  `kv:"e"`
		F string  `kv:"f"`
	}
	err := NewDecoder(strings.NewReader(str)).Decode(&st1)
	if err != nil {
		t.Fatalf("decode st1 failed: %v", err)
	}
	if st1.A != 1 {
		t.Fatalf("unexpected st1 value of a, want 1 current %d", st1.A)
	}
	if st1.B != 2 {
		t.Fatalf("unexpected st1 value of b, want 1 current %d", st1.B)
	}
	if !st1.C {
		t.Fatalf("unexpected st1 value of c, want true current %t", st1.C)
	}
	if st1.D != 3.14 {
		t.Fatalf("unexpected st1 value of d, want 3.14 current %f", st1.D)
	}
	if st1.E != "abc" {
		t.Fatalf("unexpected st1 value of e, want abc current %s", st1.E)
	}

	var st2 struct{}
	err = NewDecoder(strings.NewReader(str)).Decode(&st2)
	if err != nil {
		t.Fatalf("decode st2 failed: %v", err)
	}
}

type size uint64

func (s *size) UnmarshalKV(v string) error {
	scale := uint64(1)
	switch {
	case strings.HasSuffix(v, "G"):
		scale = 1024 * 1024 * 1024
		v = strings.TrimSuffix(v, "G")
	case strings.HasSuffix(v, "M"):
		scale = 1024 * 1024
		v = strings.TrimSuffix(v, "M")
	case strings.HasSuffix(v, "K"):
		scale = 1024
		v = strings.TrimSuffix(v, "K")
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	*s = size(uint64(n) * scale)
	return nil
}

func TestDecodeCustom(t *testing.T) {
	str := "size = 30M"
	var st struct {
		Size size `kv:"size"`
	}
	err := NewDecoder(strings.NewReader(str)).Decode(&st)
	if err != nil {
		t.Fatalf("decode custom struct failed: %v", err)
	}
	if st.Size != 30*1024*1024 {
		t.Fatalf("unexpected size value, want %d current %d", 30*1024*1024, st.Size)
	}
}

func TestDecodeEmpty(t *testing.T) {
	str := ""
	var st struct {
		Dummy string `kv:"-"`
	}
	err := NewDecoder(strings.NewReader(str)).Decode(&st)
	if err != nil {
		t.Fatalf("decode empty struct failed: %v", err)
	}
	if len(st.Dummy) > 0 {
		t.Fatal("unexpected decode empty struct")
	}
}

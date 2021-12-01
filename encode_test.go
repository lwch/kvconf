package kvconf

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestEncodeUnsupported(t *testing.T) {
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(1)
	if err.Error() != "kvconf: unsupported type: int" {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestEncodeMap(t *testing.T) {
	mp := map[string]interface{}{
		"a": "b",
		"c": 2,
	}
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(mp)
	if err != nil {
		t.Fatalf("encode map failed: %v", err)
	}
	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}
	if !strings.Contains(buf.String(), "a=b"+newLine) {
		t.Fatalf("encode map failed: unexpected a value\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "c=2"+newLine) {
		t.Fatalf("encode map failed: unexpected c value\n%s", buf.String())
	}
}

func TestEncodeStruct(t *testing.T) {
	var st struct {
		A int     `kv:"a"`
		B uint    `kv:"b"`
		C bool    `kv:"c"`
		D float64 `kv:"d"`
		E string  `kv:"e"`
	}
	st.A = 1
	st.B = 2
	st.C = true
	st.D = 3.14
	st.E = "abc"
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(st)
	if err != nil {
		t.Fatalf("encode struct failed: %v", err)
	}
	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}
	if buf.String() != "a=1"+newLine+
		"b=2"+newLine+
		"c=true"+newLine+
		"d=3.14"+newLine+
		"e=abc"+newLine {
		t.Fatalf("encode struct failed\n%s", buf.String())
	}
}

func (s size) MarshalKV() (string, error) {
	switch {
	case uint64(s) >= 1024*1024*1024:
		return fmt.Sprintf("%dG", s/(1024*1024*1024)), nil
	case uint64(s) >= 1024*1024:
		return fmt.Sprintf("%dM", s/(1024*1024)), nil
	case uint64(s) >= 1024:
		return fmt.Sprintf("%dK", s/1024), nil
	default:
		return fmt.Sprintf("%dB", s), nil
	}
}

func TestEncodeCustom(t *testing.T) {
	var st struct {
		Size size `kv:"size"`
	}
	st.Size = 30 * 1024 * 1024
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(&st)
	if err != nil {
		t.Fatalf("encode custom struct failed: %v", err)
	}
	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}
	if buf.String() != "size=30M"+newLine {
		t.Fatalf("encode custom struct failed\n%s", buf.String())
	}
}

func TestEncodeEmpty(t *testing.T) {
	var st struct {
		Dummy string `kv:"-"`
	}
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(&st)
	if err != nil {
		t.Fatalf("encode empty struct: %v", err)
	}
	if len(buf.String()) > 0 {
		t.Fatalf("encode empty struct failed, is not empty")
	}
}

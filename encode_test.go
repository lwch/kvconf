package kvconf

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func TestEncodeMap(t *testing.T) {
	mp := map[string]interface{}{
		"a": "b",
		"c": 2,
	}
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(&mp)
	if err != nil {
		t.Fatalf("encode map failed: %v", err)
	}
	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}
	if buf.String() != "a=b"+newLine+
		"c=2"+newLine {
		t.Fatalf("encode map failed\n%s", buf.String())
	}
}

func TestEncodeStruct(t *testing.T) {
	var st struct {
		A int    `kv:"a"`
		B string `kv:"b"`
	}
	st.A = 1
	st.B = "b"
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
		"b=b"+newLine {
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

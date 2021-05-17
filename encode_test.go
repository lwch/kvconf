package kvconf

import (
	"bytes"
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

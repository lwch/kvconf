package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/lwch/kvconf"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

type custom uint16

// UnmarshalKV custom unmarshaler interface
func (ct *custom) UnmarshalKV(v string) error {
	n, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		return err
	}
	*ct = custom(n)
	return nil
}

// MarshalKV custom marshaler interface
func (ct custom) MarshalKV() (string, error) {
	return fmt.Sprintf("%d", ct), nil
}

var st struct {
	Listen custom `kv:"listen"` // support custom marshaler and unmarshaler
	LogDir string `kv:"log_dir"`
}

var mp map[string]string

func decode() {
	data := `
	# trim the left space
	listen = 8080
	log_dir = /var/log/kvconf.log # trim this`

	fmt.Println("======= decode =======")

	// decode to struct
	assert(kvconf.NewDecoder(strings.NewReader(data)).Decode(&st))
	fmt.Printf("struct: listen=%d, log_dir=%s\n", st.Listen, st.LogDir)

	// decode to map[string]string
	assert(kvconf.NewDecoder(strings.NewReader(data)).Decode(&mp))
	fmt.Printf("map[string]string: listen=%s, log_dir=%s\n", mp["listen"], mp["log_dir"])
}

func encode() {
	fmt.Println("======= encode =======")

	// encode from struct
	var buf bytes.Buffer
	assert(kvconf.NewEncoder(&buf).Encode(st))
	fmt.Printf("struct:\n%s\n", buf.String())

	// encode from map[string]string
	buf.Reset()
	assert(kvconf.NewEncoder(&buf).Encode(mp))
	fmt.Printf("map[string]string:\n%s\n", buf.String())
}

func main() {
	decode()
	encode()
}

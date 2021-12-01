package kvconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Decoder unmarshal decoder
type Decoder struct {
	r *bufio.Scanner
}

// NewDecoder create decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewScanner(r)}
}

// Decode decode value from io.Reader, supported map[string]string and struct with kv tag
func (d *Decoder) Decode(v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return errors.New("input value is not pointer")
	}

	reg := regexp.MustCompile("#.+$")
	line := 0
	for d.r.Scan() {
		line++
		row := reg.ReplaceAllString(d.r.Text(), "")
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		kv := strings.SplitN(row, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid row format on line %d", line)
		}
		sk, sv := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		err := d.fill(line, vv.Elem(), sk, sv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) fill(line int, value reflect.Value, k, v string) error {
	switch value.Kind() {
	case reflect.Map:
		if value.IsNil() {
			mp := reflect.MakeMap(reflect.TypeOf(map[string]string{}))
			mp.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			value.Set(mp)
			return nil
		}
		value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		return nil
	case reflect.Struct:
		t := value.Type()
		for i := 0; i < t.NumField(); i++ {
			kf := t.Field(i)
			if kf.Tag.Get("kv") == "-" {
				continue
			}
			if kf.Tag.Get("kv") == k {
				return d.set(line, value.Field(i), v)
			}
		}
		return nil
	case reflect.Interface:
		if value.IsNil() {
			mp := reflect.MakeMap(reflect.TypeOf(map[string]string{}))
			mp.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			value.Set(mp)
			return nil
		}
		if mp, ok := value.Interface().(map[string]string); ok {
			mp[k] = v
			return nil
		}
		return fmt.Errorf("unexpected interface type on line %d, want nil", line)
	default:
		return fmt.Errorf("invalid value on line %d, expected map or struct", line)
	}
}

// Unmarshaler custom unmarshaler interface
type Unmarshaler interface {
	UnmarshalKV(string) error
}

func (d *Decoder) set(line int, value reflect.Value, v string) error {
	if value.Type().NumMethod() > 0 && value.CanInterface() {
		if dec, ok := value.Interface().(Unmarshaler); ok {
			err := dec.UnmarshalKV(v)
			if err != nil {
				return fmt.Errorf("unmarshal custom value on line %d failed, err=%v", line, err)
			}
			return nil
		}
	}
	if value.CanAddr() && value.Addr().Type().NumMethod() > 0 && value.Addr().CanInterface() {
		if dec, ok := value.Addr().Interface().(Unmarshaler); ok {
			err := dec.UnmarshalKV(v)
			if err != nil {
				return fmt.Errorf("unmarshal custom value on line %d failed, err=%v", line, err)
			}
			return nil
		}
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("can not parse value %s to int on line %d, err=%v", v, line, err)
		}
		value.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("can not parse value %s to uint on line %d, err=%v", v, line, err)
		}
		value.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("can not parse value %s to float on line %d, err=%v", v, line, err)
		}
		value.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("can not parse value %s to bool on line %d, err=%v", v, line, err)
		}
		value.SetBool(b)
	case reflect.String:
		value.SetString(v)
	default:
		return &UnsupportedTypeError{value.Type()}
	}
	return nil
}

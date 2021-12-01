package kvconf

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
)

// Encoder marshal encoder
type Encoder struct {
	w io.Writer
}

// NewEncoder create marshaler
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encode value to io.Writer, supported map[string]string and struct with kv tag
func (e *Encoder) Encode(v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	switch vv.Kind() {
	case reflect.Map:
		return e.encodeMap(vv)
	case reflect.Struct:
		return e.encodeStruct(vv)
	default:
		return &UnsupportedTypeError{vv.Type()}
	}
}

func toString(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", value.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", value.Float()), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", value.Bool()), nil
	case reflect.String:
		return value.String(), nil
	default:
		return "", &UnsupportedTypeError{value.Type()}
	}
}

func (e *Encoder) encodeMap(value reflect.Value) error {
	keys := value.MapKeys()
	for _, k := range keys {
		v := value.MapIndex(k)
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		str, err := toString(v)
		if err != nil {
			return err
		}
		if runtime.GOOS == "windows" {
			_, err = fmt.Fprintf(e.w, "%s=%s\r\n", k.String(), str)
		} else {
			_, err = fmt.Fprintf(e.w, "%s=%s\n", k.String(), str)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Marshaler custom marshaler interface
type Marshaler interface {
	MarshalKV() (string, error)
}

func (e *Encoder) encodeStruct(value reflect.Value) error {
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		vv := value.Field(i)
		kf := t.Field(i)
		k := kf.Tag.Get("kv")
		if len(k) == 0 {
			continue
		}
		if k == "-" {
			continue
		}
		if vv.Type().NumMethod() > 0 && vv.CanInterface() {
			if enc, ok := vv.Interface().(Marshaler); ok {
				v, err := enc.MarshalKV()
				if err != nil {
					return fmt.Errorf("marshal custom value failed, err=%v", err)
				}
				if runtime.GOOS == "windows" {
					_, err = fmt.Fprintf(e.w, "%s=%s\r\n", k, v)
				} else {
					_, err = fmt.Fprintf(e.w, "%s=%s\n", k, v)
				}
				if err != nil {
					return fmt.Errorf("marshal custom value failed, write error=%v", err)
				}
				continue
			}
		}
		if vv.CanAddr() && vv.Addr().Type().NumMethod() > 0 && vv.Addr().CanInterface() {
			if enc, ok := vv.Addr().Interface().(Marshaler); ok {
				v, err := enc.MarshalKV()
				if err != nil {
					return fmt.Errorf("marshal custom value failed, err=%v", err)
				}
				if runtime.GOOS == "windows" {
					_, err = fmt.Fprintf(e.w, "%s=%s\r\n", k, v)
				} else {
					_, err = fmt.Fprintf(e.w, "%s=%s\n", k, v)
				}
				if err != nil {
					return fmt.Errorf("marshal custom value failed, write error=%v", err)
				}
				continue
			}
		}
		str, err := toString(vv)
		if err != nil {
			return err
		}
		if runtime.GOOS == "windows" {
			_, err = fmt.Fprintf(e.w, "%s=%s\r\n", k, str)
		} else {
			_, err = fmt.Fprintf(e.w, "%s=%s\n", k, str)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

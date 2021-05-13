package kvconf

import "reflect"

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "kvconf: unsupported type: " + e.Type.String()
}

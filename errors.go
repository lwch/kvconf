package kvconf

import "reflect"

// UnsupportedTypeError error for unsupported type
type UnsupportedTypeError struct {
	Type reflect.Type
}

// Error error message
func (e *UnsupportedTypeError) Error() string {
	return "kvconf: unsupported type: " + e.Type.String()
}

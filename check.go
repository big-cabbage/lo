package lo

import "reflect"

// Zeroable is an interface that can be used to check if a value is zero.
type Zeroable interface {
	IsZero() bool
}

// IsZero returns true if the value is zero.
func IsZero(v any) bool {
	if z, ok := v.(Zeroable); ok {
		return z.IsZero()
	}
	return reflect.ValueOf(v).IsZero()
}

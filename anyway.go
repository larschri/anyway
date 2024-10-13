package anyway

import (
	"fmt"
)

// Must panics if err is non-nil. It can be used to wrap results from functions
// that returns a value and an error.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// Lookup traterses the given object to find the value for the given keys. It
// supports the objects created when unmarshaling (yaml og json) into a
// destination object of type 'any'.
func Lookup[T any](obj any, keys ...any) (T, error) {
	var ret T
	for i, k := range keys {
		switch v := k.(type) {
		case int:
			sl, ok := obj.([]any)
			if !ok {
				return ret, fmt.Errorf("not a slice at %v:%v", i, k)
			}
			if len(sl) <= v {
				return ret, fmt.Errorf("out of bounds at %v:%v", i, k)
			}
			obj = sl[v]
		case string:
			mp, ok := obj.(map[string]any)
			if !ok {
				return ret, fmt.Errorf("not a map at %v:%v", i, k)
			}
			o, ok := mp[v]
			if !ok {
				return ret, fmt.Errorf("missing key at %v:%v", i, k)
			}
			obj = o
		default:
			return ret, fmt.Errorf("invalid key type at %v:%T", i, k)
		}
	}
	ret, ok := obj.(T)
	if !ok {
		return ret, fmt.Errorf("invalid return type %T", obj)
	}
	return ret, nil
}

// Skeleton converts the given object by replacing potentially sensitive string
// content. Slices are truncated to size 2, and numbers are zeroed. The idea is
// to produce an object that is compatible with the original value, but can be
// stored and used in tests without leaking any sensitive data.
func Skeleton(obj any, prefix string) any {
	switch v := obj.(type) {
	case map[string]any:
		for k, o := range v {
			v[k] = Skeleton(o, prefix+"."+k)
		}
		return v
	case []any:
		for k, o := range v {
			v[k] = Skeleton(o, fmt.Sprintf("%v[%v]", prefix, k))
			if k > 0 {
				break
			}
		}
		return v
	case string:
		return prefix
	case int:
		return 0
	case float64:
		return 0.0
	case bool:
		return false
	default:
		return fmt.Sprintf("Unsupported type %T", obj)
	}
}

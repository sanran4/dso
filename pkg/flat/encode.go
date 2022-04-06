package flat

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	// ScopeDelimiter is the character used for joining multi-level strings
	// make configurable with a tag?
	ScopeDelimiter = "."
)

var nilTime time.Time
var timeType = reflect.TypeOf(nilTime)
var enumerables = map[reflect.Kind]bool{reflect.Slice: true, reflect.Array: true}

// Flatten returns all keys and corresponding values of a struct in a one-level-deep map
func Flatten(in interface{}) map[string]string {
	if in == nil {
		return nil
	}
	m := map[string]string{}
	appendValue(m, reflect.ValueOf(in), "")
	return m
}

func appendValue(m map[string]string, v reflect.Value, key string) {
	// iterate pointers until the value
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct && v.Type() != timeType {
		// recurse to fields of structs
		for i := 0; i < v.Type().NumField(); i++ {
			sf := v.Type().Field(i) // sub-field
			sv := v.Field(i)        // sub-value
			// skip unexported fields
			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}
			sk := sf.Name // sub-key
			if key != "" {
				sk = key + ScopeDelimiter + sk
			}
			appendValue(m, sv, sk)
		}
		return
	}

	// set empty key for non-struct types to the kind
	if key == "" {
		key = v.Kind().String()
		if v.Type() == timeType {
			key = "time.Time"
		}
	}

	// make 1 vs. 0 indexing configurable with tags?
	if enumerables[v.Kind()] {
		for i := 0; i < v.Len(); i++ {
			sk := key + ScopeDelimiter + strconv.Itoa(i)
			v := valueString(v.Index(i)) // TODO recursive value adding
			if v != "" {
				m[sk] = v
			}
		}
		return
	}
	// write simple values to output map
	vs := valueString(v)
	if vs != "" {
		m[key] = vs
	}
}

// valueString returns the string representation of a value
func valueString(v reflect.Value) string {
	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		if t == nilTime {
			return ""
		}
		return t.Format(time.RFC3339)
	}
	return fmt.Sprint(v.Interface())
}

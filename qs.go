// Package qs encodes structs into URL query strings.
package qs

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Encode returns the HTTP URL Query Values from the given structure v.
//
// It panics if v is not a struct or a struct that contains values other
// than bool, int, float, string, slice/array of those or another slice
// with String(). If v is nil, zero value for url.Values is returned.
//
// Literal values are encoded with their Go string representations.
//
// Pointer fields in v are derefenced until they are no longer pointer.
//
// Each field in v struct gets encoded as an URL Value unless:
//   - the field's tag is "-", or
//   - the field is empty and its tag specifies the "omitempty" option.
//
// Similar to encoding/json package, the default key is the field name
// but can be specified in the struct field's tag value. The "qs" key in
// the struct field's tag value is the key name, followed by an optional comma
// and an optional "omitempty" to discard the key when the value of the
// field is its zero value.
func Encode(v interface{}) url.Values {
	qs := url.Values{}
	t := reflect.TypeOf(v)

	if v == nil {
		return qs
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("qs: passed value (type: %T) is not a struct", v))
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("qs")
		if tag == "-" {
			continue
		}
		p := strings.Split(tag, ",")
		key := p[0]
		if key == "" {
			key = f.Name // Fallback to field name
		}

		omitEmpty := false
		if len(p) > 1 && p[1] == "omitempty" {
			omitEmpty = true
		}

		val := reflect.ValueOf(v).Field(i)

		// Follow pointers
		for val.Kind() == reflect.Ptr && !val.IsNil() {
			val = reflect.Indirect(val)
		}

		// See if the value is default and we can omitempty
		vv := val.Interface()
		def := reflect.Zero(f.Type).Interface()
		if reflect.DeepEqual(vv, def) {
			if omitEmpty {
				continue
			}
			if vv == nil {
				vv = "" // prevent output from being "key=<nil>"
			}
		}

		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			for j := 0; j < val.Len(); j++ {
				qs.Add(key, fmt.Sprintf("%v", val.Index(j)))
			}
		case reflect.Map,
			reflect.Chan,
			reflect.Func,
			reflect.UnsafePointer,
			reflect.Complex64, reflect.Complex128:
			panic(fmt.Sprintf("qs: type cannot be serialized: %s", f.Type))
		default:
			qs.Add(key, fmt.Sprintf("%v", vv))
		}
	}
	return qs
}

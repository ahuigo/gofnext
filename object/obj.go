package object

import (
	"fmt"
	"reflect"
	"strings"
)

type PtrSeen map[uintptr]interface{}

func (ps PtrSeen) Add(rv reflect.Value) uintptr {
	ptr := rv.Pointer()
	if _, ok := ps[ptr]; ok {
		e := fmt.Sprintf("encountered a cycle via %s", rv.Type())
		panic(e)
	}
	ps[ptr] = struct{}{}
	return ptr
}

func ConvertObjectByte2String(o interface{}) interface{} {
	refV := reflect.ValueOf(o)
	// var out interface{}
	ps := PtrSeen{}
	out := convertObjectByte2String(refV, ps)
	return out
}

func convertObjectByte2String(rv reflect.Value, ptrSeen PtrSeen) (out interface{}) {
	outV := reflect.ValueOf(&out).Elem()
	t := rv.Type()

	switch t.Kind() {
	// array slice
	case reflect.Slice:
		ptr := ptrSeen.Add(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Array:
		// 882 go/1.18.1/libexec/src/encoding/json/encode.go
		if t.Elem().Kind() == reflect.Uint8 {
			return string(rv.Bytes())
		} else {
			s := make([]interface{}, rv.Len())
			n := rv.Len()
			for i := 0; i < n; i++ {
				s[i] = convertObjectByte2String(rv.Index(i), ptrSeen)
			}
			return s
		}
		// return rv.Interface()
	// map
	case reflect.Map:
		ptr := ptrSeen.Add(rv)
		defer delete(ptrSeen, ptr)

		// 798 go/1.18.1/libexec/src/encoding/json/encode.go
		m := map[string]interface{}{}
		mi := rv.MapRange()
		for i := 0; mi.Next(); i++ {
			k := mi.Key()
			v := mi.Value()
			m[k.String()] = convertObjectByte2String(v, ptrSeen)
		}
		outV.Set(reflect.ValueOf(m))

	// pointer
	// case reflect.Pointer:
	case reflect.Ptr:
		ptr := ptrSeen.Add(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Interface:
		if rv.IsNil() {
			return out
		}
		return convertObjectByte2String(rv.Elem(), ptrSeen)
	case reflect.Struct:
		m := map[string]interface{}{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			key := f.Name
			isOmitEmpty := false
			if f.Tag != "" {
				key = f.Tag.Get("json")
				keys := strings.Split(key, ",")
				key = keys[0]
				if len(keys) > 1 && keys[1] == "omitempty" {
					isOmitEmpty = true
				}
			}
			fv := rv.Field(i)
			// m[key] = fv.Interface()
			v := convertObjectByte2String(fv, ptrSeen)
			if v == nil && isOmitEmpty {
				continue
			}
			m[key] = v
		}
		outV.Set(reflect.ValueOf(m))

	case reflect.Chan:
	default:
		return rv.Interface()
	}
	return
}

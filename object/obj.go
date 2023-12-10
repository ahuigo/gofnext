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
	case reflect.Slice:
		ptr := ptrSeen.Add(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Array:
		// 882 go/1.18.1/libexec/src/encoding/json/encode.go
		return convertSliceArray(rv, ptrSeen)
	case reflect.Map:
		m := convertMap(rv, ptrSeen)
		outV.Set(reflect.ValueOf(m))
	case reflect.Ptr:
		ptr := ptrSeen.Add(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Interface:
		return convertPtrInterface(rv, ptrSeen)
	case reflect.Struct:
		m := convertStruct(rv, ptrSeen)
		outV.Set(reflect.ValueOf(m))
	case reflect.Chan:
	default:
		return rv.Interface()
	}
	return
}

func convertSliceArray(rv reflect.Value, ptrSeen PtrSeen) interface{} {
	if rv.Type().Elem().Kind() == reflect.Uint8 {
		return string(rv.Bytes())
	}
	s := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		s[i] = convertObjectByte2String(rv.Index(i), ptrSeen)
	}
	return s
}

func convertMap(rv reflect.Value, ptrSeen PtrSeen) interface{} {
	m := make(map[string]interface{})
	mi := rv.MapRange()
	for mi.Next() {
		k := mi.Key()
		v := mi.Value()
		m[k.String()] = convertObjectByte2String(v, ptrSeen)
	}
	return m
}

func convertPtrInterface(rv reflect.Value, ptrSeen PtrSeen) interface{} {
	if rv.IsNil() {
		return nil
	}
	return convertObjectByte2String(rv.Elem(), ptrSeen)
}

func convertStruct(rv reflect.Value, ptrSeen PtrSeen) interface{} {
	m := make(map[string]interface{})
	t := rv.Type()
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
		v := convertObjectByte2String(fv, ptrSeen)
		if v == nil && isOmitEmpty {
			continue
		}
		m[key] = v
	}
	return m
}

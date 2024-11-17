package objectfunc

import (
	"fmt"
	"reflect"
	"strings"
)

type _StructByte2String map[uintptr]interface{}

func ConvertObjectByte2String(o interface{}) interface{} {
	refV := reflect.ValueOf(o)
	// var out interface{}
	sb := &_StructByte2String{}
	out := sb.convertObjectByte2String(refV)
	return out
}

func (sb *_StructByte2String) convertObjectByte2String(rv reflect.Value) (out interface{}) {
	outV := reflect.ValueOf(&out).Elem()
	t := rv.Type()
	ptrSeen := *sb
	switch t.Kind() {
	case reflect.Slice:
		ptr := sb.AddPtrSeen(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Array:
		// 882 go/1.18.1/libexec/src/encoding/json/encode.go
		return sb.convertSliceArray(rv)
	case reflect.Map:
		m := sb.convertMap(rv)
		outV.Set(reflect.ValueOf(m))
	case reflect.Ptr:
		ptr := sb.AddPtrSeen(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Interface:
		return sb.convertPtrInterface(rv)
	case reflect.Struct:
		m := sb.convertStruct(rv)
		outV.Set(reflect.ValueOf(m))
	case reflect.Chan:
	default:
		return rv.Interface()
	}
	return
}

func (sb *_StructByte2String) AddPtrSeen(rv reflect.Value) uintptr {
	ptrSeen := *sb
	ptr := rv.Pointer()
	if _, ok := ptrSeen[ptr]; ok {
		e := fmt.Sprintf("encountered a cycle via %s", rv.Type())
		panic(e)
	}
	ptrSeen[ptr] = struct{}{}
	return ptr
}

func (sb *_StructByte2String) convertSliceArray(rv reflect.Value) interface{} {
	if rv.Type().Elem().Kind() == reflect.Uint8 {
		return string(rv.Bytes())
	}
	s := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		s[i] = sb.convertObjectByte2String(rv.Index(i))
	}
	return s
}

func (sb *_StructByte2String) convertMap(rv reflect.Value) interface{} {
	m := make(map[string]interface{})
	mi := rv.MapRange()
	for mi.Next() {
		k := mi.Key()
		v := mi.Value()
		m[k.String()] = sb.convertObjectByte2String(v)
	}
	return m
}

func (sb *_StructByte2String) convertPtrInterface(rv reflect.Value) interface{} {
	if rv.IsNil() {
		return nil
	}
	return sb.convertObjectByte2String(rv.Elem())
}

func (sb *_StructByte2String) convertStruct(rv reflect.Value) interface{} {
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
		v := sb.convertObjectByte2String(fv)
		if v == nil && isOmitEmpty {
			continue
		}
		m[key] = v
	}
	return m
}

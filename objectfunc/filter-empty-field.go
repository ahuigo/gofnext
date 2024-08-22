package objectfunc

import (
	"fmt"
	"reflect"
	"strings"
)

type _FilterEmptyField map[uintptr]interface{}

func FilterObjectEmptyField(o interface{}) interface{} {
	refV := reflect.ValueOf(o)
	// var out interface{}
	ins := &_FilterEmptyField{}
	out := ins.filterObjectEmptyField(refV)
	return out
}

func (ins *_FilterEmptyField) filterObjectEmptyField(rv reflect.Value) (out interface{}) {
	outV := reflect.ValueOf(&out).Elem()
	t := rv.Type()
	ptrSeen := *ins
	switch t.Kind() {
	case reflect.Slice:
		ptr := ins.AddPtrSeen(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Array:
		// 882 go/1.18.1/libexec/src/encoding/json/encode.go
		// return ins.filterSliceArray(rv)
		return rv.Interface()
	case reflect.Map:
		m := ins.filterMap(rv)
		outV.Set(reflect.ValueOf(m))
	case reflect.Struct:
		m := ins.filterStruct(rv)
		outV.Set(reflect.ValueOf(m))
	case reflect.Ptr:
		ptr := ins.AddPtrSeen(rv)
		defer delete(ptrSeen, ptr)
		fallthrough
	case reflect.Interface:
		return ins.filterPtrInterface(rv)
	case reflect.Chan:
	default:
		return rv.Interface()
	}
	return
}

func (ins *_FilterEmptyField) AddPtrSeen(rv reflect.Value) uintptr {
	ptrSeen := *ins
	ptr := rv.Pointer()
	if _, ok := ptrSeen[ptr]; ok {
		e := fmt.Sprintf("encountered a cycle via %s", rv.Type())
		panic(e)
	}
	ptrSeen[ptr] = struct{}{}
	return ptr
}

// func (ins *_FilterEmptyField) filterSliceArray(rv reflect.Value) interface{} {
// 	if rv.Type().Elem().Kind() == reflect.Uint8 {
// 		return string(rv.Bytes())
// 	}
// 	s := make([]interface{}, rv.Len())
// 	for i := 0; i < rv.Len(); i++ {
// 		s[i] = ins.filterObjectEmptyField(rv.Index(i))
// 	}
// 	return s
// }

func (ins *_FilterEmptyField) filterMap(rv reflect.Value) interface{} {
	m := make(map[string]interface{})
	mr := rv.MapRange()
	for mr.Next() {
		v := mr.Value()
		// key2 := mr.Key().String()
		// _ = key2
		if ins.isEmptyField(v) {
			continue
		}
		key := mr.Key().String()
		m[key] = v.Interface()
	}
	return m
}
func (ins *_FilterEmptyField) isEmptyField(v reflect.Value) bool {
	// v.IsNil() is only for pointer, channel, func, interface, map, or slice
	// !v.IsValid() : if v is from map[notExistKey], v is invalid
	if v.IsZero() {
		return true
	}
	switch v.Kind() {
	// case Chan, Func, Interface, Map, Pointer, Slice, UnsafePointer:
	case reflect.Slice:
		return v.Len() == 0
	default:
		return false
	}

}

func (ins *_FilterEmptyField) filterPtrInterface(rv reflect.Value) interface{} {
	if rv.IsNil() {
		return nil
	}
	return ins.filterObjectEmptyField(rv.Elem())
}

func (ins *_FilterEmptyField) filterStruct(rv reflect.Value) interface{} {
	m := make(map[string]interface{})
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Name
		// 1. is empty
		fv := rv.Field(i)
		if ins.isEmptyField(fv){
			continue
		}
		// 2. set json tag name
		if f.Tag != "" {
			key = f.Tag.Get("json")
			keys := strings.Split(key, ",")
			key = keys[0]
		}
		v := fv.Interface()
		m[key] = v
	}
	return m
}

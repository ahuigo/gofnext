package serial

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

type PtrSeen map[uintptr]struct{}

func (ps PtrSeen) Add(rv reflect.Value) bool {
	ptr := rv.Pointer()
	if _, ok := ps[ptr]; ok {
		// e := fmt.Sprintf("encountered a cycle via %s", rv.Type())
		// panic(e)
		return false
	}
	ps[ptr] = struct{}{}
	return true
}

// Dump any value to string(include private field)
func String(val any, hashPtrAddr bool) string {
	refV := reflect.ValueOf(val)
	ps := PtrSeen{}
	return string(dump(refV, hashPtrAddr, ps))
}

// Dump any value to bytes(include private field)
func Bytes(val any, hashPtrAddr bool) []byte {
	refV := reflect.ValueOf(val)
	ps := PtrSeen{}
	return dump(refV, hashPtrAddr, ps)
}

func dump(refV reflect.Value, hashPtrAddr bool, ps PtrSeen) []byte {
	var buf bytes.Buffer

	switch refV.Kind() {
	case reflect.Invalid:
		buf.WriteString("<invalid>")
	case reflect.String:
		buf.WriteString(`"`)
		buf.WriteString(refV.String())
		buf.WriteString(`"`)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("%d", refV.Int()))
	// refV.CanInt()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteString(fmt.Sprintf("%d", refV.Uint()))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("%f", refV.Float()))
	case reflect.Complex64, reflect.Complex128:
		buf.WriteString(fmt.Sprintf("%f", refV.Complex()))
	case reflect.Ptr, reflect.Interface:
		if refV.IsNil() {
			buf.WriteString("null")
		} else {
			isPtr := refV.Kind() == reflect.Ptr
			if isPtr && !ps.Add(refV) {
				buf.WriteString("<cycle pointer>")
				break
			}
			if hashPtrAddr && isPtr {
				buf.WriteString(fmt.Sprintf("*0x%x", refV.Pointer()))
			} else {
				refV = refV.Elem()
				buf.WriteString(fmt.Sprintf("&%s", dump(refV, hashPtrAddr, ps)))
			}
		}
	case reflect.Slice, reflect.Array:
		buf.WriteString("[")
		for i := 0; i < refV.Len(); i++ {
			buf.Write(dump(refV.Index(i), hashPtrAddr, ps))
			if i != refV.Len()-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
	case reflect.Struct:
		name := refV.Type().Name()
		buf.WriteString(name + "{")
		for i := 0; i < refV.NumField(); i++ {
			buf.WriteString(refV.Type().Field(i).Name)
			buf.WriteString(":")
			buf.Write(dump(refV.Field(i), hashPtrAddr, ps))
			if i != refV.NumField()-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	case reflect.Map:
		sli := make([][]byte, len(refV.MapKeys()))
		for i, key := range refV.MapKeys() {
			keyVal := append(dump(key, hashPtrAddr, ps), ':')
			valbytes := dump(refV.MapIndex(key), hashPtrAddr, ps)
			sli[i] = append(keyVal, valbytes...)
		}
		slices.SortFunc(sli, func(a, b []byte) int {
			return slices.Compare(a, b)
		})
		buf.WriteString("{")
		buf.Write(bytes.Join(sli, []byte{','}))
		buf.WriteString("}")
	case reflect.Func:
		buf.WriteString("<func>")
	case reflect.Chan:
		buf.WriteString("<chan>")
	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(refV.Bool()))
	default:
		msg := fmt.Sprintf("unsupported kind %s", refV.Kind())
		panic(msg)
	}

	return buf.Bytes()
}

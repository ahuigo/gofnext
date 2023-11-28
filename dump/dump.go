package dump

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Dump(val any) string {
	refV := reflect.ValueOf(val)
	return dump(refV)
}
func dump(refV reflect.Value) string {
	switch refV.Kind() {
	case reflect.Invalid:
		return "<invalid>"
	case reflect.String:
		return `"` + refV.String() + `"`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", refV.Int())
	// refV.CanInt()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", refV.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", refV.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", refV.Complex())
	case reflect.Ptr, reflect.Interface:
		if refV.IsNil() {
			return "<nil>"
		}
		refV = refV.Elem()
		return fmt.Sprintf("&%s:%s", refV.Type().Name(), dump(refV))
	case reflect.Slice, reflect.Array:
		ret := "["
		for i := 0; i < refV.Len(); i++ {
			ret += dump(refV.Index(i))
			if i != refV.Len()-1 {
				ret += ","
			}
		}
		ret += "]"
		return ret
	case reflect.Struct:
		ret := "{"
		for i := 0; i < refV.NumField(); i++ {
			ret += refV.Type().Field(i).Name + ":" + dump(refV.Field(i))
			if i != refV.NumField()-1 {
				ret += ","
			}
		}
		ret += "}"
		return ret
	case reflect.Map:
		sli := make([]string, len(refV.MapKeys()))
		for i, key := range refV.MapKeys() {
			sli[i] = dump(key) + ":" + dump(refV.MapIndex(key))
		}
		slices.Sort(sli)
		ret := "{" + strings.Join(sli, ",") + "}"
		return ret
	case reflect.Func:
		return "<func>"
	case reflect.Chan:
		return "<chan>"
	default:
		panic("not supported")
	}
}

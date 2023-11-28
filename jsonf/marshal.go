package jsonf

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// Marshal all data(include private field) to json format
func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	return marshalValue(rv)
}

func marshalValue(refV reflect.Value) ([]byte, error) {
	switch refV.Kind() {
	case reflect.Invalid:
		return nil, nil
	case reflect.String:
		return []byte(`"` + refV.String() + `"`), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%d", refV.Int())), nil
	// refV.CanInt()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return []byte(fmt.Sprintf("%d", refV.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return []byte(fmt.Sprintf("%f", refV.Float())), nil
	case reflect.Complex64, reflect.Complex128:
		return []byte(fmt.Sprintf("%f", refV.Complex())), nil
	case reflect.Ptr, reflect.Interface:
		if refV.IsNil() {
			return []byte("<nil>"), nil
		}
		refV = refV.Elem()
		marshalVal, err := marshalValue(refV)
		if err != nil {
			return nil, err
		}
		return marshalVal,nil
		// return []byte(fmt.Sprintf("&%s:%s", refV.Type().Name(), marshalVal)), nil
	case reflect.Slice, reflect.Array:
		ret := "["
		for i := 0; i < refV.Len(); i++ {
			marshalVal, err := marshalValue(refV.Index(i))
			if err != nil {
				return nil, err
			}
			ret += string(marshalVal)
			if i != refV.Len()-1 {
				ret += ","
			}
		}
		ret += "]"
		return []byte(ret), nil
	case reflect.Struct:
		ret := "{"
		for i := 0; i < refV.NumField(); i++ {
			fieldType := refV.Type().Field(i)
			marshalVal, err := marshalValue(refV.Field(i))
			if err != nil {
				return nil, err
			}
			// parse json tag 
			fieldName := fieldType.Tag.Get("json")
			if fieldName == "" {
				fieldName = fieldType.Name
			}
			ret += fmt.Sprintf("%#v:%s", fieldName, string(marshalVal))
			if i != refV.NumField()-1 {
				ret += ","
			}
		}
		ret += "}"
		return []byte(ret), nil
	case reflect.Map:
		sli := make([]string, len(refV.MapKeys()))
		for i, key := range refV.MapKeys() {
			keyMarshalVal, err := marshalValue(key)
			if err != nil {
				return nil, err
			}
			valMarshalVal, err := marshalValue(refV.MapIndex(key))
			if err != nil {
				return nil, err
			}
			sli[i] = string(keyMarshalVal) + ":" + string(valMarshalVal)
		}
		slices.Sort(sli)
		ret := "{" + strings.Join(sli, ",") + "}"
		return []byte(ret), nil
	case reflect.Func:
		return nil, nil
	case reflect.Chan:
		return nil, nil
	default:
		return nil, nil
	}
}

package jsonf

import (
	"bytes"
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
	var buf bytes.Buffer

	err := marshalValueToBuffer(refV, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func marshalValueToBuffer(refV reflect.Value, buf *bytes.Buffer) error {
	switch refV.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.String:
		buf.WriteString(`"`)
		buf.WriteString(refV.String())
		buf.WriteString(`"`)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("%d", refV.Int()))
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
			refV = refV.Elem()
			err := marshalValueToBuffer(refV, buf)
			if err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		buf.WriteString("[")
		for i := 0; i < refV.Len(); i++ {
			err := marshalValueToBuffer(refV.Index(i), buf)
			if err != nil {
				return err
			}
			if i != refV.Len()-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
	case reflect.Struct:
		buf.WriteString("{")
		for i := 0; i < refV.NumField(); i++ {
			fieldType := refV.Type().Field(i)
			fieldName := fieldType.Tag.Get("json")
			if fieldName == "" {
				fieldName = fieldType.Name
			}
			buf.WriteString(fmt.Sprintf("%#v:", fieldName))
			err := marshalValueToBuffer(refV.Field(i), buf)
			if err != nil {
				return err
			}
			if i != refV.NumField()-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	case reflect.Map:
		keys := refV.MapKeys()
		sli := make([]string, len(keys))
		for i, key := range keys {
			keyBuf := new(bytes.Buffer)
			err := marshalValueToBuffer(key, keyBuf)
			if err != nil {
				return err
			}
			valBuf := new(bytes.Buffer)
			err = marshalValueToBuffer(refV.MapIndex(key), valBuf)
			if err != nil {
				return err
			}
			sli[i] = keyBuf.String() + ":" + valBuf.String()
		}
		slices.Sort(sli)
		buf.WriteString("{")
		buf.WriteString(strings.Join(sli, ","))
		buf.WriteString("}")
	case reflect.Func, reflect.Chan:
		// do nothing
	default:
		return nil
	}

	return nil
}

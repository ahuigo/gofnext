package serial

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Loader struct {
	d   []byte
	pos int
}

func Load(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("non-nil pointer is required to load data into")
	}

	return (&Loader{d: data}).load(rv.Elem())
}

var _buf *bytes.Reader

func Tmpload(b []byte) error {
	_, err := _buf.Read(b)
	return err
}

func (l *Loader) loadFloat(start int, rv reflect.Value) (end int, err error) {
	buf := []byte{}
	if l.d[start] == '-' {
		buf = []byte{'-'}
		start++
	}
	i := start
	for ; i < len(l.d); i++ {
		if (l.d[i] != '.' && l.d[i] < '0') || l.d[i] > '9' {
			break
		}
		buf = append(buf, l.d[i])
	}
	f, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return -1, err
	}
	rv.SetFloat(f)
	l.pos = i
	return i, nil
}

func (l *Loader) loadInt(start int, rv reflect.Value) (end int, err error) {
	buf := []byte{}
	if l.d[start] == '-' {
		buf = []byte{'-'}
		start++
	}
	i := start
	for ; i < len(l.d); i++ {
		if l.d[i] < '0' || l.d[i] > '9' {
			break
		}
		buf = append(buf, l.d[i])
	}
	n, err := strconv.ParseInt(string(buf), 10, 64)
	if err != nil {
		return -1, err
	}
	rv.SetInt(n)
	l.pos = i
	return i, nil
}

func (l *Loader) loadString(start int, rv reflect.Value) (end int, err error) {
	if l.d[start] != '"' {
		return -1, fmt.Errorf("unterminated string: %s", l.d[start:])
	}
	buf := bytes.NewBuffer(nil)
	data := l.d
	start++
	maxIndex := len(data) - 1
	i := start
LOOP:
	for ; i <= maxIndex; i++ {
		switch data[i] {
		case '\\':
			i++
			if i > maxIndex {
				return -1, fmt.Errorf("unterminated string: %s", data[start:])
			}
			switch data[i] {
			case '"', '\\', '/':
				buf.WriteByte(data[i])
			case 'b':
				buf.WriteByte('\b')
			case 't':
				buf.WriteByte(0x09)
			case 'n':
				buf.WriteByte(0x0a)
			default:
				buf.WriteByte(data[i])
			}
		case '"':
			i++
			break LOOP
		default:
			buf.WriteByte(data[i])
		}
	}
	rv.SetString(buf.String())
	l.pos = i
	return i, nil
}

func (l *Loader) load(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		_, err := l.loadString(l.pos, rv)
		if err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, err := l.loadInt(l.pos, rv); err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if _, err := l.loadFloat(l.pos, rv); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported kind %s", rv.Kind())
	}

	return nil
}

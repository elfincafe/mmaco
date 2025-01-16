package mmaco

import (
	"bytes"
	"reflect"
	"regexp"
)

const (
	tagName     = "mmaco"
	helpCmdName = "help"
	trimSpace   = " \t\v\r\n\f"

	Unknown Kind = iota
	String
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Time

	ShortOpt ArgType = iota
	LongOpt
	Value
	LongOptValue
)

type (
	Kind    int
	ArgType int
)

func toSnakeCase(s string) string {
	name := []byte{}
	for _, c := range []byte(s) {
		if c > 64 && c < 91 {
			name = append(name, byte(95), c+32)
		} else {
			name = append(name, c)
		}
	}
	name = regexp.MustCompile(`_{2,}`).ReplaceAll(name, []byte{95})
	return string(bytes.Trim(name, "_"))
}

func isAlphaNumeric(b byte) bool {
	if b >= 48 && b <= 57 { // 0 - 9
		return true
	} else if b >= 65 && b <= 90 { // A - Z
		return true
	} else if b >= 97 && b <= 122 { // a - z
		return true
	} else {
		return false
	}
}

func getFieldKind(field reflect.StructField) Kind {
	switch field.Type.Kind() {
	case reflect.Bool:
		return Bool
	case reflect.Int:
		return Int
	case reflect.Int8:
		return Int8
	case reflect.Int16:
		return Int16
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Uint:
		return Uint
	case reflect.Uint8:
		return Uint8
	case reflect.Uint16:
		return Uint16
	case reflect.Uint32:
		return Uint32
	case reflect.Uint64:
		return Uint64
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	case reflect.String:
		return String
	case reflect.Struct:
		if field.Type.PkgPath() == "time" && field.Type.Name() == "Time" {
			return Time
		} else {
			return Unknown
		}
	default:
		return Unknown
	}
}

package mmaco

import (
	"fmt"
	"reflect"
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
)

type (
	Kind int
)

var debugMode bool

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

func Debug(msg string, args ...any) {
	if !debugMode {
		return
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	println(msg)
}

package mmaco

import (
	"reflect"
	"strings"
)

type ()

func newOption(field reflect.StructField) *option {
	o := new(option)
	o.field = field
	o.short = ""
	o.long = ""
	o.required = false
	o.desc = ""
	o.defaultValue = ""
	o.format = ""
	o.handler = ""
	return o
}

func (o *option) parseTag() {
	tags := strings.Split(o.field.Tag.Get(tagName), ",")
	key := ""
	for _, v := range tags {
		t := strings.TrimLeft(v, " \t\v\r\n\f")
		if strings.HasPrefix(strings.ToLower(t), "short=") {
			short := strings.TrimSpace(t[6:])
			o.short = short
			key = "short"
		} else if strings.HasPrefix(strings.ToLower(t), "long=") {
			long := strings.TrimSpace(t[5:])
			o.long = long
			key = "long"
		} else if strings.HasPrefix(strings.ToLower(t), "desc=") {
			o.desc = t[5:]
			key = "desc"
		} else if strings.HasPrefix(strings.ToLower(t), "default=") {
			o.defaultValue = strings.TrimSpace(t[8:])
			key = "default"
		} else if strings.HasPrefix(strings.ToLower(t), "handler=") {
			o.handler = strings.TrimSpace(t[8:])
			key = "handler"
		} else if strings.HasPrefix(strings.ToLower(t), "format=") {
			o.format = strings.TrimSpace(t[7:])
			key = "format"
		} else if strings.ToLower(strings.TrimSpace(t)) == "required" {
			o.required = true
		} else if key == "desc" {
			o.desc += "," + v // concatinate variable "v" not "t"
		} else if key == "format" {
			o.format += "," + v // concatinate variable "v" not "t"
		} else if key == "default" {
			o.defaultValue += "," + v // concatinate variable "v" not "t"
		}
	}
}

func (o *option) Name() string {
	return o.field.Name
}

func (o *option) Kind() Kind {
	switch o.field.Type.Kind() {
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
		if o.field.Type.PkgPath() == "time" && o.field.Type.Name() == "Time" {
			return Time
		} else {
			return Unknown
		}
	default:
		return Unknown
	}
}

func (o *option) Short() string {
	if o.short == "" {
		return ""
	} else {
		return "-" + o.short
	}
}

func (o *option) Long() string {
	if o.long == "" {
		return ""
	} else {
		return "--" + o.long
	}
}

func (o *option) Required() bool {
	return o.required
}

func (o *option) Desc() string {
	return o.desc
}

func (o *option) Default() string {
	return o.defaultValue
}

func (o *option) Handler() string {
	return o.handler
}

func (o *option) Format() string {
	return o.handler
}

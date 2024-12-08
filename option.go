package mmaco

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func newOption(value reflect.Value, field reflect.StructField) *option {
	if _, ok := field.Tag.Lookup(tagName); !ok {
		return nil
	}
	o := new(option)
	o.value = value
	o.field = field
	o.short = ""
	o.long = ""
	o.required = false
	o.desc = ""
	o.defaultValue = ""
	o.format = ""
	o.handler = ""
	o.specified = false
	return o
}

func (o *option) parse() error {
	tags := strings.Split(o.field.Tag.Get(tagName), ",")
	key := ""
	for _, v := range tags {
		t := strings.TrimLeft(v, trimSpace)
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
			o.defaultValue = t[8:]
			key = "default"
		} else if strings.HasPrefix(strings.ToLower(t), "handler=") {
			o.handler = strings.TrimSpace(t[8:])
			key = "handler"
		} else if strings.HasPrefix(strings.ToLower(t), "format=") {
			o.format = t[7:]
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
	if len(o.short) != 0 && len(o.short) > 1 {
		return fmt.Errorf(`"short" must be 1 character`)
	} else if len(o.short) == 1 && !isAlphaNumeric([]byte(o.short)[0]) {
		return fmt.Errorf(`"short" must be 0-9, a-z, A-Z`)
	} else if len(o.long) == 1 {
		return fmt.Errorf(`"long" must be at least 2 characters`)
	} else if len(o.short) == 0 && len(o.long) == 0 {
		return fmt.Errorf(`neither "short" nor "long" is specified`)
	} else if len(o.format) > 0 && len(o.handler) > 0 {
		return fmt.Errorf(`"format" and "handler" is exclusive`)
	}
	return nil
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

func (o *option) set(value string) error {
	switch o.Kind() {
	case Bool:
		o.value.SetBool(true)
	case Int:
		v, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int type`, o.field.Name)
		}
		o.value.SetInt(v)
		o.specified = true
	case Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int8 type`, o.field.Name)
		}
		o.value.SetInt(v)
		o.specified = true
	case Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int16 type`, o.field.Name)
		}
		o.value.SetInt(v)
		o.specified = true
	case Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int32 type`, o.field.Name)
		}
		o.value.SetInt(v)
		o.specified = true
	case Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int64 type`, o.field.Name)
		}
		o.value.SetInt(v)
		o.specified = true
	case Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint type`, o.field.Name)
		}
		o.value.SetUint(v)
		o.specified = true
	case Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint8 type`, o.field.Name)
		}
		o.value.SetUint(v)
		o.specified = true
	case Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint16 type`, o.field.Name)
		}
		o.value.SetUint(v)
		o.specified = true
	case Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint32 type`, o.field.Name)
		}
		o.value.SetUint(v)
		o.specified = true
	case Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint64 type`, o.field.Name)
		}
		o.value.SetUint(v)
		o.specified = true
	case Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the float32 type`, o.field.Name)
		}
		o.value.SetFloat(v)
		o.specified = true
	case Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the float64 type`, o.field.Name)
		}
		o.value.SetFloat(v)
		o.specified = true
	case String:
		o.value.SetString(value)
		o.specified = true
	case Time:
		t, err := time.Parse(o.format, value)
		if err != nil {
			return fmt.Errorf(`can't parse "%s" for the value of option "%s"`, value, o.field.Name)
		}
		o.value.Set(reflect.ValueOf(t))
		o.specified = true
	default:
		return fmt.Errorf(`the field type of "%s" isn't supported`, o.field.Type.Name())
	}
	return nil

}

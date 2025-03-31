package mmaco

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	option struct {
		Kind      Kind
		Name      string
		Short     string
		Long      string
		Required  bool
		Desc      string
		Default   string
		Format    string
		Handler   string
		value     reflect.Value
		field     reflect.StructField
		ctx       *Context
		specified bool
	}
)

func newOption(value reflect.Value, field reflect.StructField) *option {
	o := new(option)
	o.specified = false
	o.value = value
	o.field = field
	o.ctx = nil
	o.Kind = getFieldKind(o.field)
	o.Name = o.field.Name
	o.Short = ""
	o.Long = ""
	o.Required = false
	o.Desc = ""
	o.Format = ""
	o.Handler = ""

	tags := strings.Split(o.field.Tag.Get(tagName), ",")
	key := ""
	for _, v := range tags {
		t := strings.TrimLeft(v, trimSpace)
		if strings.HasPrefix(strings.ToLower(t), "short=") {
			short := strings.TrimSpace(t[6:])
			o.Short = short
			key = "short"
		} else if strings.HasPrefix(strings.ToLower(t), "long=") {
			long := strings.TrimSpace(t[5:])
			o.Long = long
			key = "long"
		} else if strings.HasPrefix(strings.ToLower(t), "desc=") {
			o.Desc = t[5:]
			key = "desc"
		} else if strings.HasPrefix(strings.ToLower(t), "handler=") {
			o.Handler = strings.TrimSpace(t[8:])
			key = "handler"
		} else if strings.HasPrefix(strings.ToLower(t), "format=") {
			o.Format = t[7:]
			key = "format"
		} else if strings.ToLower(strings.TrimSpace(t)) == "required" {
			o.Required = true
		} else if key == "desc" {
			o.Desc += "," + v // concatinate variable "v" not "t"
		} else if key == "format" {
			o.Format += "," + v // concatinate variable "v" not "t"
		}
	}

	return o
}

func (o *option) validate(sc *SubCommand) error {
	shortRule := regexp.MustCompile(`^[0-9a-zA-Z]$`)
	longRule := regexp.MustCompile(`^[\w_]{2,}$`)
	if o.Short != "" && !shortRule.MatchString(o.Short) {
		return fmt.Errorf(`option "-%s" don't follow the rule (%s)`, o.Short, shortRule.String())
	} else if o.Long != "" && !longRule.MatchString(o.Long) {
		return fmt.Errorf(`option "--%s" don't follow the rule (%s)`, o.Short, longRule.String())
	} else if o.Short == "" && o.Long == "" {
		return fmt.Errorf(`neither "short" nor "long" is specified`)
	} else if o.Format != "" && o.Handler != "" {
		return fmt.Errorf(`"format" and "handler" are exclusive`)
	} else if o.Handler != "" {
		v := reflect.ValueOf(sc.cmd)
		method := v.MethodByName(o.Handler)
		if !method.IsValid() {
			return fmt.Errorf(`"%s" doesn't have the method "%s"`, sc.Name, o.Handler)
		} else if method.Type().NumIn() != 1 || method.Type().In(0).Kind() != reflect.String {
			return fmt.Errorf(`"%s" must have only one argument, which is a string type`, o.Name)
		} else if method.Type().NumOut() != 1 || method.Type().Out(0).Kind() != reflect.Interface {
			return fmt.Errorf(`"%s" must have only one return value, which is a string type`, o.Name)
		}
	}
	return nil
}

func (o *option) isShort(arg string) bool {
	if o.Short != "" && arg == "-"+o.Short {
		return true
	} else {
		return false
	}
}

func (o *option) isLong(arg string) bool {
	if o.Long != "" && arg == "--"+o.Long {
		return true
	} else {
		return false
	}
}

func (o *option) has(arg string) bool {
	if o.Long != "" && strings.HasPrefix(arg, "--"+o.Long+"=") {
		return true
	} else {
		return false
	}
}

func (o *option) set(value string) error {
	if !o.value.CanSet() {
		return fmt.Errorf(`field "%s" is unexported`, o.Name)
	}
	switch o.Kind {
	case Bool:
		o.value.SetBool(true)
	case Int:
		v, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int type`, o.field.Name)
		}
		o.value.SetInt(v)
	case Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int8 type`, o.field.Name)
		}
		o.value.SetInt(v)
	case Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int16 type`, o.field.Name)
		}
		o.value.SetInt(v)
	case Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int32 type`, o.field.Name)
		}
		o.value.SetInt(v)
	case Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the int64 type`, o.field.Name)
		}
		o.value.SetInt(v)
	case Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint type`, o.field.Name)
		}
		o.value.SetUint(v)
	case Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint8 type`, o.field.Name)
		}
		o.value.SetUint(v)
	case Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint16 type`, o.field.Name)
		}
		o.value.SetUint(v)
	case Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint32 type`, o.field.Name)
		}
		o.value.SetUint(v)
	case Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the uint64 type`, o.field.Name)
		}
		o.value.SetUint(v)
	case Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the float32 type`, o.field.Name)
		}
		o.value.SetFloat(v)
	case Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf(`the value of option "%s" should be the float64 type`, o.field.Name)
		}
		o.value.SetFloat(v)
	case String:
		o.value.SetString(value)
	case Time:
		var err error
		var t time.Time
		if o.ctx.loc != nil {
			t, err = time.ParseInLocation(o.Format, value, o.ctx.loc)
		} else {
			t, err = time.Parse(o.Format, value)
		}
		if err != nil {
			return fmt.Errorf(`can't parse "%s" for the value of option "%s"`, value, o.field.Name)
		}
		o.value.Set(reflect.ValueOf(t))
	default:
		return fmt.Errorf(`the field type of "%s" isn't supported`, o.field.Type.Name())
	}
	o.specified = true
	return nil
}

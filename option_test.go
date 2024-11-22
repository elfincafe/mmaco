package mmaco

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type (
	optTest struct {
		field0  string     ``
		field1  string     `mmaco:"testtest"`
		field2  string     `mmaco:"desc=description test"`
		field10 int        `mmaco:"short=a,long=field10,desc=description10,default=10,required"`
		field11 int8       `mmaco:"short=b,long=field11,desc=description11,default=100,required"`
		field12 int16      `mmaco:"short=c,long=field12,desc=description12,default=1000,required"`
		field13 int32      `mmaco:"short=d,long=field13,desc=description13,default=10000,required"`
		field14 int64      `mmaco:"short=e,long=field14,desc=description14,default=100000,required"`
		field15 uint       `mmaco:"short=f,long=field15,desc=description15,default=10,required"`
		field16 uint8      `mmaco:"short=g,long=field16,desc=description16,default=100,required"`
		field17 uint16     `mmaco:"short=h,long=field17,desc=description17,default=1000,required"`
		field18 uint32     `mmaco:"short=i,long=field18,desc=description18,default=10000,required"`
		field19 uint64     `mmaco:"short=j,long=field19,desc=description19,default=100000,required"`
		field20 float32    `mmaco:"short=k,long=field20,desc=description20,default=1.1,required"`
		field21 float64    `mmaco:"short=l,long=field21,desc=description21,default=1.23,required"`
		field22 time.Time  `mmaco:"short=m,long=field22,desc=description22,default=123,required,format=2006/01/02 15:04:05"`
		field99 complex128 `mmaco:""`
	}
)

func TestNewOpts(t *testing.T) {
	// Test Case
	type st struct {
		f1 string
		f2 int
	}
	v := reflect.ValueOf(st{f1: "123", f2: 456})
	cases := []struct {
		field    reflect.StructField
		expected bool
	}{
		{field: v.Type().Field(0)},
		{field: v.Type().Field(1)},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		if fmt.Sprintf("%T", o) != "*mmaco.option" {
			t.Errorf("[%d] %s Result: %v", i, c.field.Name, fmt.Sprintf("%T", o))
		}
	}
}

func TestOptionParse(t *testing.T) {
	// Test Case
	type st struct {
		f0 int `mmaco:"short=f,long=field,required,desc=desc,test,default=default value,format=Mon, 02 Jan 2006 15:04:05 MST,handler=Parse"`
		f1 int `mmaco:""`
		f2 int `mmaco:"test"`
		f3 int `mmaco:"desc=desc,test,short=ff,long=field,required,default=default, value,handler=Parse,format=Mon, 02 Jan 2006 15:04:05 MST"`
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field        reflect.StructField
		short        string
		long         string
		required     bool
		desc         string
		defaultValue string
		format       string
		handler      string
	}{
		{field: v.Type().Field(0), short: "f", long: "field", required: true, desc: "desc,test", defaultValue: "default value", format: "Mon, 02 Jan 2006 15:04:05 MST", handler: "Parse"},
		{field: v.Type().Field(1), short: "", long: "", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{field: v.Type().Field(2), short: "", long: "", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{field: v.Type().Field(3), short: "", long: "field", required: true, desc: "desc,test", defaultValue: "default, value", format: "Mon, 02 Jan 2006 15:04:05 MST", handler: "Parse"},
	}
	// Test
	for i, c := range cases {
		o := newOption((c.field))
		o.parseTag()
		if o.short != c.short {
			t.Errorf("[%d] Short Expected:%v Result(%v)", i, o.short, c.short)
		}
		if o.long != c.long {
			t.Errorf("[%d] Long Expected:%v Result(%v)", i, o.long, c.long)
		}
		if o.required != c.required {
			t.Errorf("[%d] Required Expected:%v Result(%v)", i, o.required, c.required)
		}
		if o.desc != c.desc {
			t.Errorf("[%d] Desc Expected:%v Result(%v)", i, o.desc, c.desc)
		}
		if o.defaultValue != c.defaultValue {
			t.Errorf("[%d] Default Expected:%v Result(%v)", i, o.defaultValue, c.defaultValue)
		}
		if o.format != c.format {
			t.Errorf("[%d] Format Expected:%v Result(%v)", i, o.format, c.format)
		}
		if o.handler != c.handler {
			t.Errorf("[%d] Handler Expected:%v Result(%v)", i, o.handler, c.handler)
		}
	}
}

func TestOptionName(t *testing.T) {
	// Test Case
	type st struct {
		f0 int
		f1 int8
		f2 int16
		f3 int32
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected string
	}{
		{field: v.Type().Field(0), expected: "f0"},
		{field: v.Type().Field(1), expected: "f1"},
		{field: v.Type().Field(2), expected: "f2"},
		{field: v.Type().Field(3), expected: "f3"},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		if o.Name() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Name())
		}
	}
}

func TestOptionKind(t *testing.T) {
	// Test Case
	type st struct {
		f0  int
		f1  int8
		f2  int16
		f3  int32
		f4  int64
		f5  uint
		f6  uint8
		f7  uint16
		f8  uint32
		f9  uint64
		f10 float32
		f11 float64
		f12 bool
		f13 string
		f14 time.Time
		f15 complex64
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected Kind
	}{
		{field: v.Type().Field(0), expected: Int},
		{field: v.Type().Field(1), expected: Int8},
		{field: v.Type().Field(2), expected: Int16},
		{field: v.Type().Field(3), expected: Int32},
		{field: v.Type().Field(4), expected: Int64},
		{field: v.Type().Field(5), expected: Uint},
		{field: v.Type().Field(6), expected: Uint8},
		{field: v.Type().Field(7), expected: Uint16},
		{field: v.Type().Field(8), expected: Uint32},
		{field: v.Type().Field(9), expected: Uint64},
		{field: v.Type().Field(10), expected: Float32},
		{field: v.Type().Field(11), expected: Float64},
		{field: v.Type().Field(12), expected: Bool},
		{field: v.Type().Field(13), expected: String},
		{field: v.Type().Field(14), expected: Time},
		{field: v.Type().Field(15), expected: Unknown},
	}
	// Test
	for i, c := range cases {
		opt := newOption(c.field)
		if opt.Kind() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, opt.Kind())
		}
	}
}

func TestOptionShort(t *testing.T) {
	// Test Case
	type st struct {
		f0 string `mmaco:"short=f"`
		f1 string `mmaco:"short=f,long=field"`
		f2 string `mmaco:"long=field,short=f,required"`
		f3 string `mmaco:"desc=desc,with,comma,short=f,required"`
		f4 string `mmaco:"required,short=f,desc=desc,with,comma"`
		f5 string `mmaco:"short=ff"`
		f6 string `mmaco:"short=ff,long=field"`
		f7 string `mmaco:"long=field,short=ff,required"`
		f8 string `mmaco:"desc=desc,with,comma,short=ff,required"`
		f9 string `mmaco:"required,short=ff,desc=desc,with,comma"`
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected string
	}{
		{field: v.Type().Field(0), expected: "-f"},
		{field: v.Type().Field(1), expected: "-f"},
		{field: v.Type().Field(2), expected: "-f"},
		{field: v.Type().Field(3), expected: "-f"},
		{field: v.Type().Field(4), expected: "-f"},
		{field: v.Type().Field(5), expected: "-ff"},
		{field: v.Type().Field(6), expected: "-ff"},
		{field: v.Type().Field(7), expected: "-ff"},
		{field: v.Type().Field(8), expected: "-ff"},
		{field: v.Type().Field(9), expected: "-ff"},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		o.parseTag()
		if o.Short() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Short())
		}
	}
}

func TestOptionLong(t *testing.T) {
	// Test Case
	type st struct {
		f0 string `mmaco:"long=field"`
		f1 string `mmaco:"long=field,short=f"`
		f2 string `mmaco:"long=field,short=f,required"`
		f3 string `mmaco:"desc=desc,with,comma,long=field,required"`
		f4 string `mmaco:"required,long=field,desc=desc,with,comma"`
		f5 string `mmaco:"long=f"`
		f6 string `mmaco:"long=f,short=field"`
		f7 string `mmaco:"short=ff,long=f,required"`
		f8 string `mmaco:"desc=desc,with,comma,long=f,required"`
		f9 string `mmaco:"required,long=f,desc=desc,with,comma"`
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected string
	}{
		{field: v.Type().Field(0), expected: "--field"},
		{field: v.Type().Field(1), expected: "--field"},
		{field: v.Type().Field(2), expected: "--field"},
		{field: v.Type().Field(3), expected: "--field"},
		{field: v.Type().Field(4), expected: "--field"},
		{field: v.Type().Field(5), expected: "--f"},
		{field: v.Type().Field(6), expected: "--f"},
		{field: v.Type().Field(7), expected: "--f"},
		{field: v.Type().Field(8), expected: "--f"},
		{field: v.Type().Field(9), expected: "--f"},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		o.parseTag()
		if o.Long() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Long())
		}
	}
}

func TestOptionRequired(t *testing.T) {
	// Test Case
	type st struct {
		f0 string `mmaco:"required"`
		f1 string `mmaco:"long=field,required,short=f"`
		f2 string `mmaco:"long=field,short=f"`
		f3 string `mmaco:""`
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected bool
	}{
		{field: v.Type().Field(0), expected: true},
		{field: v.Type().Field(1), expected: true},
		{field: v.Type().Field(2), expected: false},
		{field: v.Type().Field(3), expected: false},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		o.parseTag()
		if o.Required() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Required())
		}
	}
}

func TestOptionDesc(t *testing.T) {
	// Test Case
	type st struct {
		f0 string `mmaco:"desc=description, test"`
		f1 string `mmaco:"short=f,desc=description, test,long=field"`
		f2 string `mmaco:"desc=description, test,long=field,short=f"`
		f3 string `mmaco:"long=field,short=f,desc=description, test"`
		f4 string `mmaco:"short=f,default=default, test,desc=description, test,long=field"`
		f5 string `mmaco:"default=default, test,desc=description, test,long=field,short=f"`
		f6 string `mmaco:"long=field,short=f,desc=description, test,default=default, test"`
	}
	v := reflect.ValueOf(st{})
	cases := []struct {
		field    reflect.StructField
		expected string
	}{
		{field: v.Type().Field(0), expected: "description, test"},
		{field: v.Type().Field(1), expected: "description, test"},
		{field: v.Type().Field(2), expected: "description, test"},
		{field: v.Type().Field(3), expected: "description, test"},
		{field: v.Type().Field(4), expected: "description, test"},
		{field: v.Type().Field(5), expected: "description, test"},
		{field: v.Type().Field(6), expected: "description, test"},
	}
	// Test
	for i, c := range cases {
		o := newOption(c.field)
		o.parseTag()
		if o.Desc() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Desc())
		}
	}
}

func TestOptionDefault(t *testing.T) {
}

func TestOptionHandler(t *testing.T) {
}

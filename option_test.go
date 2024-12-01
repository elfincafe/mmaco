package mmaco

import (
	"reflect"
	"testing"
)

func isSameOption(aSt, bSt *option) bool {
	if aSt == nil && bSt == nil {
		return false
	}
	a := reflect.TypeOf(aSt).Elem()
	b := reflect.TypeOf(bSt).Elem()
	if a.PkgPath() == b.PkgPath() && a.Name() == b.Name() {
		return true
	} else {
		return false
	}
}

func TestNewOption(t *testing.T) {
	// Test Case
	cases := []struct {
		idx int
		st  *option
	}{
		{idx: 0, st: nil},
		{idx: 1, st: nil},
		{idx: 2, st: new(option)},
		{idx: 3, st: new(option)},
		{idx: 4, st: new(option)},
		{idx: 5, st: new(option)},
		{idx: 6, st: new(option)},
		{idx: 7, st: new(option)},
		{idx: 8, st: new(option)},
		{idx: 9, st: new(option)},
		{idx: 10, st: new(option)},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] Expected: *option, Result: "%v"`, i, o.short)
			continue
		}
	}
}

func TestOptionParseTag(t *testing.T) {
	// Test Case
	v := reflect.ValueOf(subCmd0{})
	cases := []struct {
		idx          int
		st           *option
		short        string
		long         string
		required     bool
		desc         string
		defaultValue string
		format       string
		handler      string
	}{
		{idx: 0, st: nil, short: "", long: "", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{idx: 1, st: nil, short: "", long: "", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{idx: 2, st: new(option), short: "2", long: "", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{idx: 3, st: new(option), short: "", long: "field3", required: false, desc: "", defaultValue: "", format: "", handler: ""},
		{idx: 4, st: new(option), short: "", long: "", required: false, desc: "desc4", defaultValue: "", format: "", handler: ""},
		{idx: 5, st: new(option), short: "", long: "", required: true, desc: "", defaultValue: "", format: "", handler: ""},
		{idx: 6, st: new(option), short: "", long: "", required: false, desc: "", defaultValue: "default value", format: "", handler: ""},
		{idx: 7, st: new(option), short: "", long: "", required: false, desc: "", defaultValue: "", format: "2006/01/02 15:04:05", handler: ""},
		{idx: 8, st: new(option), short: "", long: "", required: false, desc: "", defaultValue: "", format: "", handler: "Handler"},
		{idx: 9, st: new(option), short: "9", long: "field9", required: true, desc: " desc9, test ", defaultValue: " default9, Value , ", format: "Mon, 02 Jan 2006 15:04:05 MST", handler: "Handler"},
		{idx: 10, st: new(option), short: "a", long: "field10", required: false, desc: " desc10, test ", defaultValue: " default10, Value , ", format: "", handler: "Handler"},
	}
	// Test
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] Expected: *option, Result: "%v"`, i, o.short)
			continue
		}
		o.parseTag()
		if o.short != c.short {
			t.Errorf(`[%d] Short Expected: "%v", Result: "%v"`, i, o.short, c.short)
		}
		if o.long != c.long {
			t.Errorf(`[%d] Long Expected: "%v" Result: "%v"`, i, o.long, c.long)
		}
		if o.required != c.required {
			t.Errorf(`[%d] Required Expected: "%v", Result: "%v"`, i, o.required, c.required)
		}
		if o.desc != c.desc {
			t.Errorf(`[%d] Desc Expected: "%v", Result: "%v"`, i, o.desc, c.desc)
		}
		if o.defaultValue != c.defaultValue {
			t.Errorf(`[%d] Default Expected: "%v", Result: "%v"`, i, o.defaultValue, c.defaultValue)
		}
		if o.format != c.format {
			t.Errorf(`[%d] Format Expected: "%v", Result: "%v"`, i, o.format, c.format)
		}
		if o.handler != c.handler {
			t.Errorf(`[%d] Handler Expected: "%v", Result: "%v"`, i, o.handler, c.handler)
		}
	}
}

func TestIsValid(t *testing.T) {
	// // Test Case
	// cases := []struct {
	// 	idx      int
	// 	st       bool
	// 	expected bool
	// }{
	// 	{idx: 0, st: false, expected: false},
	// 	{idx: 1, st: false, expected: false},
	// 	{idx: 2, st: true, expected: false},
	// 	{idx: 3, st: true, expected: false},
	// 	{idx: 4, st: true, expected: false},
	// 	{idx: 5, st: true, expected: false},
	// 	{idx: 6, st: true, expected: false},
	// 	{idx: 7, st: true, expected: false},
	// 	{idx: 8, st: true, expected: false},
	// 	{idx: 9, st: true, expected: false},
	// }
	// // Test
	// v := reflect.ValueOf(subCmd0{})
}

func TestOptionName(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: "field2"},
		{idx: 3, st: new(option), expected: "field3"},
		{idx: 4, st: new(option), expected: "field4"},
		{idx: 5, st: new(option), expected: "field5"},
		{idx: 6, st: new(option), expected: "field6"},
		{idx: 7, st: new(option), expected: "field7"},
		{idx: 8, st: new(option), expected: "field8"},
		{idx: 9, st: new(option), expected: "field9"},
		{idx: 10, st: new(option), expected: "field10"},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(i), v.Type().Field(i))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		if o.Name() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Name())
		}
	}
}

func TestOptionKind(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		value    reflect.Value
		field    reflect.StructField
		expected Kind
	}{
		{idx: 0, st: new(option), expected: String},
		{idx: 1, st: new(option), expected: Bool},
		{idx: 2, st: new(option), expected: Int},
		{idx: 3, st: new(option), expected: Int8},
		{idx: 4, st: new(option), expected: Int16},
		{idx: 5, st: new(option), expected: Int32},
		{idx: 6, st: new(option), expected: Int64},
		{idx: 7, st: new(option), expected: Uint},
		{idx: 8, st: new(option), expected: Uint8},
		{idx: 9, st: new(option), expected: Uint16},
		{idx: 10, st: new(option), expected: Uint32},
		{idx: 11, st: new(option), expected: Uint64},
		{idx: 12, st: new(option), expected: Float32},
		{idx: 13, st: new(option), expected: Float64},
		{idx: 14, st: new(option), expected: Time},
		{idx: 15, st: new(option), expected: Unknown},
	}
	// Test
	v := reflect.ValueOf(subCmd1{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		if o.Kind() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Kind())
		}
	}
}

func TestOptionShort(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: "-2"},
		{idx: 3, st: new(option), expected: ""},
		{idx: 4, st: new(option), expected: ""},
		{idx: 5, st: new(option), expected: ""},
		{idx: 6, st: new(option), expected: ""},
		{idx: 7, st: new(option), expected: ""},
		{idx: 8, st: new(option), expected: ""},
		{idx: 9, st: new(option), expected: "-9"},
		{idx: 10, st: new(option), expected: "-a"},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Short() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Short())
		}
	}
}

func TestOptionLong(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: ""},
		{idx: 3, st: new(option), expected: "--field3"},
		{idx: 4, st: new(option), expected: ""},
		{idx: 5, st: new(option), expected: ""},
		{idx: 6, st: new(option), expected: ""},
		{idx: 7, st: new(option), expected: ""},
		{idx: 8, st: new(option), expected: ""},
		{idx: 9, st: new(option), expected: "--field9"},
		{idx: 10, st: new(option), expected: "--field10"},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Long() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Long())
		}
	}
}

func TestOptionRequired(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected bool
	}{
		{idx: 0, st: nil, expected: false},
		{idx: 1, st: nil, expected: false},
		{idx: 2, st: new(option), expected: false},
		{idx: 3, st: new(option), expected: false},
		{idx: 4, st: new(option), expected: false},
		{idx: 5, st: new(option), expected: true},
		{idx: 6, st: new(option), expected: false},
		{idx: 7, st: new(option), expected: false},
		{idx: 8, st: new(option), expected: false},
		{idx: 9, st: new(option), expected: true},
		{idx: 10, st: new(option), expected: false},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Required() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Required())
		}
	}
}

func TestOptionDesc(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: ""},
		{idx: 3, st: new(option), expected: ""},
		{idx: 4, st: new(option), expected: "desc4"},
		{idx: 5, st: new(option), expected: ""},
		{idx: 6, st: new(option), expected: ""},
		{idx: 7, st: new(option), expected: ""},
		{idx: 8, st: new(option), expected: ""},
		{idx: 9, st: new(option), expected: " desc9, test "},
		{idx: 10, st: new(option), expected: " desc10, test "},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Desc() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Desc())
		}
	}
}

func TestOptionDefault(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: ""},
		{idx: 3, st: new(option), expected: ""},
		{idx: 4, st: new(option), expected: ""},
		{idx: 5, st: new(option), expected: ""},
		{idx: 6, st: new(option), expected: "default value"},
		{idx: 7, st: new(option), expected: ""},
		{idx: 8, st: new(option), expected: ""},
		{idx: 9, st: new(option), expected: " default9, Value , "},
		{idx: 10, st: new(option), expected: " default10, Value , "},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Default() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Default())
		}
	}
}

func TestOptionHandler(t *testing.T) {
	// Test Case
	cases := []struct {
		idx      int
		st       *option
		expected string
	}{
		{idx: 0, st: nil, expected: ""},
		{idx: 1, st: nil, expected: ""},
		{idx: 2, st: new(option), expected: ""},
		{idx: 3, st: new(option), expected: ""},
		{idx: 4, st: new(option), expected: ""},
		{idx: 5, st: new(option), expected: ""},
		{idx: 6, st: new(option), expected: ""},
		{idx: 7, st: new(option), expected: ""},
		{idx: 8, st: new(option), expected: "Handler"},
		{idx: 9, st: new(option), expected: "Handler"},
		{idx: 10, st: new(option), expected: "Handler"},
	}
	// Test
	v := reflect.ValueOf(subCmd0{})
	for i, c := range cases {
		o := newOption(v.Field(c.idx), v.Type().Field(c.idx))
		if o == nil && c.st == nil {
			continue
		}
		if !isSameOption(o, c.st) {
			t.Errorf(`[%d] OptionObject Expected: *option, Result: %v`, i, o)
			continue
		}
		o.parseTag()
		if o.Handler() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, o.Handler())
		}
	}
}

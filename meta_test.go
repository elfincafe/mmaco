package mmaco

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewMeta(t *testing.T) {
	m := newMeta()
	if m.short != "" {
		t.Errorf("Short: %v", m.short)
	}
	if m.long != "" {
		t.Errorf("Long: %v", m.long)
	}
	if m.required != false {
		t.Errorf("Required: %v", m.required)
	}
	if m.desc != "" {
		t.Errorf("Desc: %v", m.desc)
	}
	if m.value != "" {
		t.Errorf("Value: %v", m.value)
	}
	if m.handler != "" {
		t.Errorf("Handler: %v", m.handler)
	}
}

func TestGetMetas(t *testing.T) {
	type (
		st1 struct {
			p1 int    `mmaco:"short=p,long=prop1,required,value=123,desc=p1 description"`
			p2 string `mmaco:"short=P,long=prop2,desc=p2 description"`
		}
		st2 struct { // Wrong tag
			p1 uint      `maco:"short=p,long=prop1"`
			p2 time.Time `mmaco:"short=P,handler=OptionHandler"`
		}
		st3 struct { // No tags
			p1 int64
			p2 float64
		}
	)
	// Cases
	cases := []struct {
		rf  reflect.Type
		exp map[string]meta
	}{
		{
			rf: reflect.TypeOf(st1{}),
			exp: map[string]meta{
				"p1": meta{short: "p", long: "prop1", required: true, value: "123", desc: "p1 description"},
				"p2": meta{short: "P", long: "prop2", desc: "p1 description"},
			},
		},
		{
			rf: reflect.TypeOf(st2{}),
			exp: map[string]meta{
				"p2": meta{short: "P", handler: "OptionHandler"},
			},
		},
		{
			rf:  reflect.TypeOf(st3{}),
			exp: map[string]meta{},
		},
	}
	// Test
	for i, c := range cases {
		v := getMetas(c.rf)
		if len(v) != len(c.exp) {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.exp, v)
		}
	}
}

func TestGetMeta(t *testing.T) {
	// Cases
	cases := []struct {
		tag string
		exp meta
	}{
		{
			tag: "",
			exp: meta{},
		},
		{
			tag: "short=v",
			exp: meta{short: "v"},
		},
		{
			tag: "long=verbose",
			exp: meta{long: "verbose"},
		},
		{
			tag: "required",
			exp: meta{required: true},
		},
		{
			tag: "desc=test description",
			exp: meta{desc: "test description"},
		},
		{
			tag: "value=123",
			exp: meta{value: "123"},
		},
		{
			tag: "handler=OptionHandler",
			exp: meta{handler: "OptionHandler"},
		},
		{
			tag: "short=v,long=verbose,required,desc=test description,value=123,handler=OptionHandler",
			exp: meta{short: "v", long: "verbose", required: true, desc: "test description", value: "123", handler: "OptionHandler"},
		},
		{ // Space
			tag: "short=v ,  long=verbose",
			exp: meta{short: "v", long: "verbose"},
		},
		{ // Multiple
			tag: "short=v ,  long=verbose, short=t",
			exp: meta{short: "t", long: "verbose"},
		},
		{ // including a comma in value.
			tag: "desc= test, description",
			exp: meta{desc: " test, description"},
		},
		{ // including some commas in value.
			tag: "desc= test, descr,iption",
			exp: meta{desc: " test, descr,iption"},
		},
	}
	// Test
	for i, c := range cases {
		m := getMeta(c.tag)
		if fmt.Sprintf("%T", m) != "mmaco.meta" {
			t.Errorf("[%d Short] Expected: mmaco, Returned: %T", i, m)
		}
		if c.exp.short != m.short {
			t.Errorf("[%d Short] Expected: '%v', Returned: '%v'", i, c.exp.short, m.short)
		}
		if c.exp.long != m.long {
			t.Errorf("[%d Long] Expected: '%v', Returned: '%v'", i, c.exp.long, m.long)
		}
		if c.exp.required != m.required {
			t.Errorf("[%d Required] Expected: %v, Returned: %v", i, c.exp.required, m.required)
		}
		if c.exp.desc != m.desc {
			t.Errorf("[%d Desc] Expected: '%v', Returned: '%v'", i, c.exp.desc, m.desc)
		}
		if c.exp.value != m.value {
			t.Errorf("[%d Value] Expected: '%v', Returned: '%v'", i, c.exp.value, m.value)
		}
		if c.exp.handler != m.handler {
			t.Errorf("[%d Handler] Expected: '%v', Returned: '%v'", i, c.exp.handler, m.handler)
		}
	}
}

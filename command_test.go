package mmaco

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Test Case
	cmd := New(tagName)
	typ := reflect.TypeOf(cmd).Elem().Name()
	pkg := reflect.TypeOf(cmd).Elem().PkgPath()
	fullType := fmt.Sprintf("*%s.%s", pkg, typ)
	expected := fmt.Sprintf("*%s.Command", tagName)
	// Test
	if fullType != expected {
		t.Errorf("Expected: %v, Result: %v", expected, fullType)
	}
}

func TestCommandParse(t *testing.T) {
	cases := []struct {
		name  string
		short string
		long  string
		kind  Kind
	}{
		{name: "help", short: "h", long: "help", kind: Bool},
		{name: "verbose", short: "v", long: "verbose", kind: Bool},
	}
	cmd := New(tagName)
	cmd.parse()
	var opt *option
	for i, c := range cases {
		opt = nil
		for _, o := range cmd.opts {
			o.parseTag()
			if o.Name() == c.name {
				opt = o
				break
			}
		}
		if opt != nil {
			if opt.short != c.short || opt.long != c.long || opt.Kind() != c.kind {
				t.Errorf("[%d] Short:%v(%v), Long:%v(%v), Kind:%v(%v), ", i, opt.short, c.short, opt.long, c.long, opt.Kind(), c.kind)
			}
		} else {
			t.Errorf("[%d] Can't find the field '%s'", i, c.name)
		}
	}
}

func TestCommandAdd(t *testing.T) {
	// Test Case
	cases := []struct {
		sc   SubCommandInterface
		name string
	}{
		{
			sc:   subCmdTest1{},
			name: "sub_cmd_test1",
		},
		{
			sc:   subCmdTest2{},
			name: "sub_cmd_test2",
		},
		{
			sc:   subCmdTest3{},
			name: "sub_cmd_test3",
		},
	}
	// Test
	for i, c := range cases {
		cmd := New("cmd")
		cmd.Add(c.sc)
		if cmd.subCmds[c.name].Name() != c.name {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.sc, cmd.subCmds[c.name])
		}
	}
}

func TestCommandRoute(t *testing.T) {
	cases := []struct {
		args []string
	}{
		{
			[]string{"-v", "--help", "download", "--date=202410"},
		},
	}
	cmd := New("test")
	for _, c := range cases {
		cmd.scOrder = append(cmd.scOrder, "create")
		cmd.scOrder = append(cmd.scOrder, "download")
		cmd.route(c.args)
	}
}

func TestCommandRun(t *testing.T) {

}

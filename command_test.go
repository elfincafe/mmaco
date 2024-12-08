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
			o.parse()
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
			sc:   subCmd0{},
			name: "sub_cmd0",
		},
		{
			sc:   subCmd1{},
			name: "sub_cmd1",
		},
		{
			sc:   subCmd2{},
			name: "sub_cmd2",
		},
	}
	// Test
	for i, c := range cases {
		cmd := New("cmd")
		cmd.parse()
		cmd.Add(c.sc)
		if cmd.subCmds[c.name].Name != c.name {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.sc, cmd.subCmds[c.name])
		}
	}
}

func TestCommandRoute(t *testing.T) {
	cases := []struct {
		args    []string
		help    bool
		verbose bool
		subcmd  string
	}{
		{
			[]string{"-v", "-h"},
			true,
			true,
			"",
		},
		{
			[]string{"--help", "--verbose"},
			true,
			true,
			"",
		},
		{
			[]string{"-v", "create"},
			false,
			true,
			"create",
		},
	}
	// Test
	cmd := New(tagName)
	cmd.parse()
	fmt.Println(cmd.opts)
	for _, c := range cases {
		cmd.scOrder = append(cmd.scOrder, "create")
		cmd.scOrder = append(cmd.scOrder, "download")
		cmd.route(c.args)
		fmt.Println(cmd.help, cmd.verbose, cmd.subCmd)
	}
}

func TestCommandGetSubCmdIndex(t *testing.T) {
	// Test Case
	cases := []struct {
		args     []string
		expected int
	}{
		{
			[]string{},
			-1,
		},
		{
			[]string{"-v"},
			-1,
		},
		{
			[]string{"-v", "create"},
			1,
		},
		{
			[]string{"-v", "-h", "create"},
			2,
		},
		{
			[]string{"-v", "-h", "help"},
			-1,
		},
	}
	// Test
	cmd := New(tagName)
	for i, c := range cases {
		cmd.scOrder = append(cmd.scOrder, "create")
		cmd.scOrder = append(cmd.scOrder, "download")
		idx := cmd.getSubCmdIndex(c.args)
		if c.expected != idx {
			t.Errorf(`[%d] Expected: %v, Result: %v`, i, c.expected, idx)
		}
	}
}

func TestCommandRun(t *testing.T) {

}

func TestCommandHelpCommand(t *testing.T) {
	cmd := New(tagName)
	cmd.Add(subCmd0{Desc: "Sub Command 0 for Test"})
	cmd.Add(subCmd1{})
	cmd.Add(subCmd2{})
	cmd.parse()
	cmd.helpCommand()
}

func TestCommandHelpSubCommand(t *testing.T) {
	cmd := New(tagName)
	cmd.Add(subCmd0{Desc: "Sub Command 0 for Test"})
	cmd.parse()
	cmd.route([]string{"-h", "sub_cmd0"})
	cmd.helpSubCommand()
}

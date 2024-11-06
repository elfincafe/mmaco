package mmaco

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

}

type (
	cmd1 struct {
	}
	cmd2 struct {
	}
	cmd3 struct {
	}
)

func (c cmd1) Exec() error {
	return nil
}
func (c cmd2) Run() error {
	return nil
}
func (c cmd3) Init() error {
	return nil
}
func (c cmd3) Validate() error {
	return nil
}
func (c cmd3) Run() error {
	return nil
}

func TestAddSubCommand(t *testing.T) {
	// Cases
	cases := []struct {
		name string
		st   any
	}{
		{
			name: "cmd",
			st: subCommand{
				cmd:  reflect.ValueOf(cmd1{}),
				meta: meta{},
			},
		},
		{
			name: "cmd",
			st: subCommand{
				cmd:  reflect.ValueOf(cmd2{}),
				meta: meta{},
			},
		},
		{
			name: "cmd",
			st: subCommand{
				cmd:  reflect.ValueOf(cmd3{}),
				meta: meta{},
			},
		},
	}
	// Test
	for i, c := range cases {
		cmd := New("cmd")
		cmd.AddSubCommand(c.name, c.st)
		if cmd.subCmds[c.name] == c.st {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.st, cmd.subCmds[c.name])
		}
	}
}

func TestConvSubCommandName(t *testing.T) {
	// Cases
	cases := []struct {
		name string
		exp  string
	}{
		{name: "abc", exp: "abc"},
		{name: "ABC", exp: "abc"},
		{name: "A0B1C2", exp: "a0b1c2"},
		{name: "aB_C", exp: "ab_c"},
		{name: "AbC_0", exp: "abc_0"},
		{name: "0ABc_", exp: "abc"},
		{name: "_AbC_9", exp: "abc_9"},
		{name: "_AあbC_9", exp: "abc_9"},
		{name: "AbC-", exp: "abc"},
	}
	// Test
	cmd := New("test")
	for i, c := range cases {
		v := cmd.convSubCommandName(c.name)
		if v != c.exp {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.exp, v)
		}
	}
}

func TestRoute(t *testing.T) {
}

func TestRun(t *testing.T) {

}

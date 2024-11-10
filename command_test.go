package mmaco

import (
	"fmt"
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

func TestCommandAddSubCommand(t *testing.T) {
	// Cases
	cases := []struct {
		st any
	}{
		{
			st: &subCommand{
				cmd:  reflect.ValueOf(cmd1{}),
				meta: newMeta(),
			},
		},
		{
			st: &subCommand{
				cmd:  reflect.ValueOf(cmd2{}),
				meta: newMeta(),
			},
		},
		{
			st: &subCommand{
				cmd:  reflect.ValueOf(cmd3{}),
				meta: newMeta(),
			},
		},
	}
	fmt.Println(cases)
	// Test
	// for i, c := range cases {
	// 	cmd := New("cmd")
	// 	cmd.AddSubCommand(c.st)
	// 	if cmd.subCmds[c.st.getName()] == c.st {
	// 	t.Errorf("[%d] Expected: %v, Returned: %v", i, c.st, cmd.subCmds[c.name])
	// 	}
	// }
}

func TestCommandRoute(t *testing.T) {
	cases := []struct {
		args []string
	}{
		{
			[]string{"-v", "--help", "download"},
		},
	}
	cmd := New("test")
	for _, c := range cases {
		cmd.route(c.args)
	}
}

func TestCommandRun(t *testing.T) {

}

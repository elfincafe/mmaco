package mmaco

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {

}

type (
	sc struct {
	}
)

func (c sc) Init() error {
	return nil
}
func (c sc) Validate() error {
	return nil
}
func (c sc) Run(args []string) error {
	return nil
}

func TestCommandAdd(t *testing.T) {
	// Cases
	cases := []struct {
		st any
	}{
		{
			st: newSubCommand(subCmdTest{}),
		},
		{
			st: newSubCommand(subCmdTest{}),
		},
		{
			st: newSubCommand(subCmdTest{}),
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

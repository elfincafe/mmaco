package mmaco

import "testing"

type (
	subCmdTest1 struct {
	}
	subCmdTest2 struct {
	}
	subCmdTest3 struct {
	}
)

func (sc subCmdTest1) Init() error {
	return nil
}
func (sc subCmdTest1) Validate() error {
	return nil
}
func (sc subCmdTest1) Run() error {
	return nil
}
func (sc subCmdTest2) Init() error {
	return nil
}
func (sc subCmdTest2) Validate() error {
	return nil
}
func (sc subCmdTest2) Run() error {
	return nil
}
func (sc subCmdTest3) Init() error {
	return nil
}
func (sc subCmdTest3) Validate() error {
	return nil
}
func (sc subCmdTest3) Run() error {
	return nil
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	cleanup()
}

func setup() {

}

func cleanup() {

}

func TestToSnakeCase(t *testing.T) {
	// Test Case
	cases := []struct {
		str      string
		expected string
	}{
		{str: "TestSubCommand", expected: "test_sub_command"},
		{str: "TestSubCommand2", expected: "test_sub_command2"},
		{str: "Test_Sub_Command_3", expected: "test_sub_command_3"},
		{str: "4Test_Sub_Command_", expected: "4_test_sub_command"},
	}
	// Test
	for i, c := range cases {
		name := toSnakeCase(c.str)
		if name != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, name)
		}
	}
}

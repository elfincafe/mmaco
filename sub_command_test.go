package mmaco

import (
	"fmt"
	"testing"
)

type st struct {
	f0 string
	f1 string `mmaco:"short=s,long=option,desc=description,"`
	f2 string `mmaco:"short=s,long=option,desc=description,"`
	// Init     func() error
	// Validate func() error
	// Run      func() error
}

func (sc st) Init() error {
	return nil
}
func (sc st) Validate() error {
	return nil
}
func (sc st) Run() error {
	return nil
}

func TestNewSubCommand(t *testing.T) {
	// Test Case
	v := newSubCommand(st{})
	// Test
	result := fmt.Sprintf("%T", v)
	expected := fmt.Sprintf("*%s.subCommand", tagName)
	if result != expected {
		t.Errorf("Expected: %v, Result: %v", expected, result)
	}
	// fmt.Printf("%T", (*v).cmd)
	// if fmt.Sprintf("%T", (*v).cmd) != "reflect.Type" {
	// 	t.Errorf("Expected: reflect.Type, Result: %v", fmt.Sprintf("%T", (*v).cmd))
	// }
}

func TestSubCommandParse(t *testing.T) {
	// Test Case
	v := newSubCommand(st{})
	expectedCount := 0
	for i := 0; i < v.cmd.NumField(); i++ {
		_, ok := v.cmd.Field(i).Tag.Lookup(tagName)
		if ok {
			expectedCount += 1
		}
	}
	v.parse()
	// Test
	if len(v.opts) != expectedCount {
		t.Errorf("Expected: %v, Result: %v", expectedCount, len(v.opts))
	}
}

func TestSubCommandName(t *testing.T) {
	// Test Case
	cases := []struct {
		sc       SubCommandInterface
		expected string
	}{
		{sc: subCmdTest1{}, expected: "sub_cmd_test1"},
		{sc: subCmdTest2{}, expected: "sub_cmd_test2"},
		{sc: subCmdTest3{}, expected: "sub_cmd_test3"},
	}
	// Test
	for i, c := range cases {
		cmd := newSubCommand(c.sc)
		if cmd.Name() != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, cmd.Name())
		}
	}
}

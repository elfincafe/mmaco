package mmaco

import (
	"fmt"
	"reflect"
	"testing"
)

type (
	sc1 struct {
		name string
	}
)

func (s sc1) Init() error {
	return nil
}
func (s sc1) Validate() error {
	return nil
}
func (s sc1) Run() error {
	return nil
}
func (s sc1) Name() string {
	return s.name
}

func TestNewSubCommand(t *testing.T) {
	s := sc1{}
	sc := newSubCommand(s)
	scType := reflect.TypeOf(sc).Name()
	if scType != "subCommand" {
		t.Errorf("sc == %s", scType)
	}
	cmdType := reflect.TypeOf(sc.cmd).Name()
	if cmdType != "Value" {
		t.Errorf("SubCommand Reflector is not '%s'", cmdType)
	}
	metaType := reflect.TypeOf(sc.meta).Name()
	if metaType != "meta" {
		t.Errorf("meta is not '%s'", metaType)
	}
}

func TestSubCommandValidate(t *testing.T) {
	sc := newSubCommand(sc1{})
	sc.validate()
}

func TestSubCommandGetName(t *testing.T) {
	cases := []struct {
		st   SubCommandInterface
		name string
	}{
		{st: sc1{name: "123"}, name: ""},
	}
	for i, c := range cases {
		sc := newSubCommand(c.st)
		fmt.Println("**", sc.Name(), "**")
		if sc.Name() != c.name {
			t.Errorf("[%d] Expected: %v, Returned: %v", i, c.name, sc.Name())
		}
	}
}

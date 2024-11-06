package mmaco

import (
	"reflect"
	"testing"
)

type (
	sc1 struct {
		opt1 bool `mmaco:"short:"`
	}
)

func (s *sc1) Init() error {
	return nil
}
func (s *sc1) Validate() error {
	return nil
}
func (s *sc1) Run() error {
	return nil
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

func TestValidate(t *testing.T) {

}

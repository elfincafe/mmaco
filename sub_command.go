package mmaco

import (
	"fmt"
	"reflect"
)

type (
	subCommand struct {
		cmd  reflect.Value
		meta meta
	}
)

func newSubCommand(s any) subCommand {
	sc := subCommand{}
	sc.cmd = reflect.ValueOf(s)
	sc.meta = meta{}
	return sc
}

func (sc *subCommand) validate() error {
	// Init, Validate
	for _, method := range []string{"Init", "Validate"} {
		m := sc.cmd.MethodByName(method)
		if m.IsValid() {
			if m.Kind() != reflect.Func || m.Type().NumOut() != 1 || m.Type().Out(0).String() != "error" {
				return fmt.Errorf(`%s.Run field is NOT a method which must returns only an error`, sc.cmd.Type())
			}
		}
	}
	// Run
	run := sc.cmd.MethodByName("Run")
	if !run.IsValid() {
		return fmt.Errorf(`%s has no "Run" method`, sc.cmd.Type())
	}
	if run.Kind() != reflect.Func || run.Type().NumOut() != 1 || run.Type().Out(0).String() != "error" {
		return fmt.Errorf(`%s.Run field is NOT a method which must returns only an error`, sc.cmd.Type())
	}
	return nil
}

package mmaco

import (
	"reflect"
)

func newSubCommand(s SubCommandInterface) *subCommand {
	sc := new(subCommand)
	sc.cmd = reflect.ValueOf(&s).Elem()
	sc.opts = []*option{}
	return sc
}

func (sc *subCommand) parse() {
	for i := 0; i < sc.cmd.NumField(); i++ {
		value := sc.cmd.Field(i)
		field := sc.cmd.Type().Field(i)
		opt := newOption(value, field)
		if opt != nil {
			sc.opts = append(sc.opts, opt)
		}
	}
}

func (sc *subCommand) Name() string {
	return toSnakeCase(sc.cmd.Type().Name())
}

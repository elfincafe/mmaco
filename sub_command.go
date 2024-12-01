package mmaco

import (
	"reflect"
)

func newSubCommand(s SubCommandInterface) *subCommand {
	sc := new(subCommand)
	sc.cmd = reflect.ValueOf(s)
	sc.opts = []*option{}
	return sc
}

func (sc *subCommand) parse() {
	for i := 0; i < sc.cmd.NumField(); i++ {
		o := newOption(sc.cmd.Field(i), sc.cmd.Type().Field(i))
		if o != nil {
			sc.opts = append(sc.opts, o)
		}
	}
}

func (sc *subCommand) Name() string {
	return toSnakeCase(sc.cmd.Type().Name())
}

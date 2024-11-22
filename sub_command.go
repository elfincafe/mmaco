package mmaco

import (
	"reflect"
)

type (
	subCommand struct {
		cmd  reflect.Value
		opts []*option
	}
)

func newSubCommand(s SubCommandInterface) *subCommand {
	sc := new(subCommand)
	sc.cmd = reflect.ValueOf(s)
	sc.opts = []*option{}
	return sc
}

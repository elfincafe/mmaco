package mmaco

import (
	"fmt"
	"reflect"
	"regexp"
)

type (
	subCommand struct {
		cmd  reflect.Value
		meta *meta
	}
)

func newSubCommand(s SubCommandInterface) *subCommand {
	sc := new(subCommand)
	sc.cmd = reflect.ValueOf(s)
	sc.meta = newMeta()
	return sc
}

func (sc *subCommand) validate() error {
	// Name rules
	name := sc.Name()
	re := regexp.MustCompile(`^[a-zA-Z][\w_]*[0-9a-zA-Z]$`)
	if !re.MatchString(name) {
		return fmt.Errorf("SubCommand name '%s' is't doesn't follow the rule", name)
	}
	return nil
}

func (sc *subCommand) Name() string {
	ret := sc.cmd.MethodByName("Name").Call([]reflect.Value{})
	if len(ret) > 0 {
		return ret[0].String()
	}
	return ""
}

package mmaco

import (
	"fmt"
	"reflect"
)

func newSubCommand(s SubCommandInterface) *subCommand {
	sc := new(subCommand)
	sc.cmd = reflect.ValueOf(s)
	sc.Name = toSnakeCase(sc.cmd.Type().Name())
	sc.Desc = fmt.Sprintf("%s command", sc.cmd.Type().Name())
	sc.opts = []*option{}
	return sc
}

func (sc *subCommand) parse() error {
	var err error
	// Desc
	_, ok := sc.cmd.Type().FieldByName("Desc")
	if ok && sc.cmd.FieldByName("Desc").Kind() == reflect.String {
		sc.Desc = sc.cmd.FieldByName("Desc").String()
	}
	// Field
	for i := 0; i < sc.cmd.NumField(); i++ {
		o := newOption(sc.cmd.Field(i), sc.cmd.Type().Field(i))
		if o == nil {
			continue
		}
		err = o.parse()
		if err != nil {
			return err
		}
		sc.opts = append(sc.opts, o)
	}
	return nil
}

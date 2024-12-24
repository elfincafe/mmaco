package mmaco

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func newSubCommand(s any, loc *time.Location) *subCommand {
	t := reflect.TypeOf(s)
	if !isSubCommand(t) {
		return nil
	}

	sc := new(subCommand)
	sc.loc = loc
	sc.cmd = reflect.ValueOf(s)
	sc.opts = []*option{}
	sc.Name = toSnakeCase(sc.cmd.Elem().Type().Name())
	sc.hasValidate = hasValidateMethod(t)

	if hasInitMethod(t) {
		sc.cmd.MethodByName("Init").Call(nil)
	}
	if hasDescField(t) {
		sc.Desc = sc.cmd.Elem().FieldByName("Desc").String()
	}

	return sc
}

func (sc *subCommand) parse() error {
	var err error

	// Field
	v := sc.cmd.Elem()
	t := sc.cmd.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		o := newOption(v.Field(i), t.Field(i), sc.loc)
		if o == nil {
			continue
		}
		err = o.validate(sc)
		if err != nil {
			return err
		}
		sc.opts = append(sc.opts, o)
	}

	return nil
}

func (sc *subCommand) parseArgs(args []string) ([]string, error) {
	var err error
	in, out := []reflect.Value{}, []reflect.Value{}
	params := []string{}
	maxIdx := len(args) - 1
	skip := false
	setFlg := false
	for i, arg := range args {
		setFlg = false
		if skip {
			skip = false
			continue
		}
		for _, o := range sc.opts {
			err = nil
			if (o.isShort(arg) || o.isLong(arg)) && o.Kind == Bool {
				if o.Handler == "" {
					err = o.set("true")
					setFlg = true
				} else {
					in = []reflect.Value{reflect.ValueOf("true")}
					out = sc.cmd.MethodByName(o.Handler).Call(in)
					setFlg = true
					if !out[0].IsNil() {
						err = out[0].Interface().(error)
					}
				}
				if err != nil {
					return nil, err
				} else {
					break
				}
			} else if o.isShort(arg) && o.Kind != Bool {
				argVal := ""
				if maxIdx > i {
					argVal = args[i+1]
					if strings.HasPrefix(argVal, "-") {
						return nil, fmt.Errorf(`the option "%s" needs a value`, arg)
					}
				} else {
					return nil, fmt.Errorf(`the option "%s" needs a value`, arg)
				}
				if o.Handler == "" {
					err = o.set(argVal)
					setFlg = true
					skip = true
				} else {
					in = []reflect.Value{reflect.ValueOf(argVal)}
					out = sc.cmd.MethodByName(o.Handler).Call(in)
					setFlg = true
					skip = true
					if !out[0].IsNil() {
						err = out[0].Interface().(error)
					}
				}
				if err != nil {
					return nil, err
				}
			} else if o.has(arg) {
				length := len("--" + o.Long + "=")
				argVal := arg[length:]
				if o.Handler == "" {
					err = o.set(argVal)
					setFlg = true
				} else {
					in = []reflect.Value{reflect.ValueOf(argVal)}
					out = sc.cmd.MethodByName(o.Handler).Call(in)
					setFlg = true
					if !out[0].IsNil() {
						err = out[0].Interface().(error)
					}
				}
				if err != nil {
					return nil, err
				}
			}
		}
		if !setFlg {
			params = append(params, arg)
		}
	}

	return params, nil
}

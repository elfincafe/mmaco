package mmaco

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	SubCommandInterface interface {
		Init()
		Validate() error
		Run(*Context) error
	}
	SubCommand struct {
		Name string
		Desc string
		cmd  reflect.Value
		ctx  *Context
		opts []*option
	}
)

func newSubCommand(s SubCommandInterface) *SubCommand {

	s.Init()
	t := reflect.TypeOf(s)
	sc := new(SubCommand)
	sc.Name = toSnakeCase(t.Elem().Name())
	sc.cmd = reflect.ValueOf(s)
	sc.Desc = ""
	sc.opts = []*option{}

	// description
	field := sc.cmd.Elem().FieldByName("Desc")
	if field.IsValid() {
		sc.Desc = field.String()
	}

	return sc
}

func (sc *SubCommand) parse() error {
	var err error

	// Field
	t := sc.cmd
	for i := 0; i < t.Elem().NumField(); i++ {
		o := newOption(sc.cmd.Elem().Field(i), t.Elem().Type().Field(i))
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

func (sc *SubCommand) parseArgs(args []string) ([]string, error) {
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
					if err != nil {
						return nil, err
					}
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
					if err != nil {
						return nil, err
					}
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
					if err != nil {
						return nil, err
					}
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

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
		cmd  SubCommandInterface
		ctx  *Context
		opts []*option
	}
)

func newSubCommand(s SubCommandInterface, name, desc string) *SubCommand {
	sc := new(SubCommand)
	sc.Name = name
	sc.Desc = desc
	sc.cmd = s
	sc.opts = []*option{}
	return sc
}

func (sc *SubCommand) parse() {
	// Field
	v := reflect.ValueOf(sc.cmd).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		if f.Kind() == reflect.Pointer {
			continue
		}
		tag := ft.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		opt := newOption(f, ft)
		sc.ctx.subCmd.opts = append(sc.ctx.subCmd.opts, opt)
	}
}

func (sc *SubCommand) parseArgs(args []string) ([]string, error) {
	var err error
	v := reflect.ValueOf(sc.cmd)
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
		for _, o := range sc.ctx.subCmd.opts {
			err = nil
			if (o.isShort(arg) || o.isLong(arg)) && o.Kind == Bool {
				if o.Handler == "" {
					err = o.set("true")
					setFlg = true
				} else {
					err = o.setByHandler(v, "true")
					setFlg = true
				}
				if err != nil {
					return nil, err
				}
				break
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
					err = o.setByHandler(v, argVal)
					setFlg = true
					skip = true
				}
				if err != nil {
					return nil, err
				}
				break
			} else if o.has(arg) {
				length := len("--" + o.Long + "=")
				argVal := arg[length:]
				if o.Handler == "" {
					err = o.set(argVal)
					setFlg = true
				} else {
					err = o.setByHandler(v, argVal)
					setFlg = true
				}
				if err != nil {
					return nil, err
				}
				break
			}
		}
		if !setFlg {
			params = append(params, arg)
		}
	}
	// check required
	for _, o := range sc.opts {
		if o.Required && !o.specified {
			optName := ""
			if o.Short != "" && o.Long == "" {
				optName = "-" + o.Short
			} else if o.Short == "" && o.Long != "" {
				optName = "--" + o.Long
			} else if o.Short != "" && o.Long != "" {
				optName = fmt.Sprintf(`-%s, --%s`, o.Short, o.Long)
			}
			return nil, fmt.Errorf(`option "%s" is required`, optName)
		}
	}

	return params, nil
}

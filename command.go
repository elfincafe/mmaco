package mmaco

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func New(name string) *Command {
	cmd := new(Command)
	cmd.start = time.Now()
	cmd.Name = name
	cmd.subCmds = map[string]*subCommand{}
	cmd.subCmd = ""
	cmd.scOrder = []string{}
	cmd.opts = []*option{}
	return cmd
}

func (cmd *Command) parse() {
	ref := reflect.TypeOf(cmd).Elem()
	for i := 0; i < ref.NumField(); i++ {
		opt := newOption(ref.Field(i))
		if opt != nil {
			cmd.opts = append(cmd.opts, opt)
		}
	}
}

func (cmd *Command) Add(subCmd SubCommandInterface) {
	sc := newSubCommand(subCmd)
	sc.parse()
	name := toSnakeCase(sc.Name())
	cmd.subCmds[name] = sc
	exists := false
	for _, v := range cmd.scOrder {
		if name == v {
			exists = true
			break
		}
	}
	if !exists {
		cmd.scOrder = append(cmd.scOrder, name)
	}
}

func (cmd *Command) route(args []string) error {
	var err error
	idx := 0
	skip := false
	for i, arg := range args {
		ok := false
		if skip {
			continue
		}
		for _, opt := range cmd.opts {
			if arg == opt.short {
				if opt.Kind() == Bool {
					// setArg(&opt.field, opt.Short(), "true")
				} else if opt.Kind() != Unknown {
					// setArg(&opt.field, opt.Short(), args[i+1])
				} else {

				}
				break
			} else if arg == opt.long {
				if opt.Kind() == Bool {
					// err = setArg()
				}
				break
			} else if strings.HasPrefix(arg, opt.long+"=") {
				length := len(opt.long + "=")
				// err = setArg(opt.field, arg[length:])
				break
			}
		}
		if ok {
			idx = i + 1
		} else {
			break
		}
	}
	fmt.Println(args[:idx], idx, len(args[idx:]))
	os.Exit(123)

	// SubCommand
	if len(args[idx:]) > 0 {
		ok := false
		for _, subcmd := range cmd.scOrder {
			if args[idx] == subcmd {
				cmd.subCmd = subcmd
				ok = true
				break
			}
			if ok {
				break
			}
		}
	} else {
		return fmt.Errorf("SubCommand isn't passed")
	}
	idx += 1
	fmt.Println(cmd.subCmd, args[idx:])

	// skip := false
	// for i, arg := range args {
	// 	if skip {
	// 		continue
	// 	}
	// 	for name, meta := range metas {
	// 		field := c.Elem().FieldByName(name)
	// 		kind := field.Kind()
	// 		short := "-" + meta.short
	// 		long := "--" + meta.long
	// 		if arg == long {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if kind == reflect.Bool {
	// 				cmd.setArg(&field, long, "true")
	// 				break
	// 			} else {
	// 				err = fmt.Errorf("needs value for the '%s' (e.g. --%s=something)", name, meta.long)
	// 			}
	// 		} else if arg == short {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if kind == reflect.Bool {
	// 				cmd.setArg(&field, short, "true")
	// 				break
	// 			}
	// 			if len(args) >= i && !strings.HasPrefix(args[i+1], "-") {
	// 				skip = true
	// 				break
	// 			}
	// 		} else if strings.HasPrefix(arg, long) {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if !strings.HasPrefix(arg, long+"=") {
	// 				err = fmt.Errorf("needs value for the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			n := len("--" + meta.long + "=")
	// 			cmd.setArg(&field, long, arg[n:])
	// 			break
	// 		}
	// 	}
	// 	idx = i - 1
	// }

	// fmt.Println(metas["help"], metas["verbose"])
	return nil
	// ERROR:
	// 	return err
}

func (cmd *Command) Run() error {
	// Routing
	subCmdPos := cmd.route(os.Args[1:])
	fmt.Println(subCmdPos, cmd.subCmd)

	// Analizing
	sc := reflect.ValueOf(cmd.subCmd)

	// Intialize
	init := sc.MethodByName("Init")
	if init.IsValid() {
		init.Call([]reflect.Value{})
	}

	// Parsing Arguments

	// Run
	sc.MethodByName("Run").Call([]reflect.Value{})

	return nil
}

func (cmd *Command) Report() {
	if !cmd.verbose {
		return
	}
}

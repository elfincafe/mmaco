package mmaco

import (
	"fmt"
	"math"
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

func (cmd *Command) parse() error {
	var err error
	v := reflect.ValueOf(cmd)
	for i := 0; i < v.Type().Elem().NumField(); i++ {
		opt := newOption(v.Elem().Field(i), v.Type().Elem().Field(i))
		if opt == nil {
			continue
		}
		err = opt.parse()
		if err != nil {
			return err
		}
		cmd.opts = append(cmd.opts, opt)
	}
	return nil
}

func (cmd *Command) Add(subCmd SubCommandInterface) {
	sc := newSubCommand(subCmd)
	sc.parse()
	name := sc.Name
	if name == helpCmdName {
		return
	}
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

	subCmdIdx := cmd.getSubCmdIndex(args)
	if subCmdIdx < 0 {
		cmd.help = true
		return nil
	} else {
		cmd.subCmd = args[subCmdIdx]
	}

	// parse root options
	skip := false
	for i, arg := range args[:subCmdIdx] {
		if skip {
			skip = false
			continue
		}
		for _, opt := range cmd.opts {
			if arg == opt.Short {
				if opt.Kind() == Bool {
					opt.set("true")
				} else if opt.Kind() != Unknown {
					if i+1 < subCmdIdx {
						opt.set(args[i+1])
						skip = true
					} else {
						return fmt.Errorf(`option "%s" needs a value`, opt.Name())
					}
				} else {
					return fmt.Errorf(`option "%s" needs a value`, opt.Name())
				}
			} else if arg == opt.Long {
				if opt.Kind() == Bool {
					opt.set("true")
				} else {
					return fmt.Errorf(`option "%s" needs a value`, opt.Name())
				}
			} else if strings.HasPrefix(arg, opt.Long+"=") {
				length := len(opt.Long + "=")
				opt.set(arg[length:])
				break
			}
		}

	}

	// parse sub command options
	if len(args[subCmdIdx:]) > 0 {
	}

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

func (cmd *Command) getSubCmdIndex(args []string) int {
	idx := -1
	for i, arg := range args {
		for _, scName := range cmd.scOrder {
			if arg == scName {
				idx = i
				break
			}
		}
		if idx > -1 {
			break
		}
	}
	return idx
}

func (cmd *Command) Run() error {
	in, out := []reflect.Value{}, []reflect.Value{}

	// Routing
	subCmdPos := cmd.route(os.Args[1:])
	fmt.Println(subCmdPos, cmd.subCmd)

	// Analizing
	sc := reflect.ValueOf(cmd.subCmds[cmd.subCmd])

	// Intialize
	init := sc.MethodByName("Init")
	if init.IsValid() && init.Type().NumIn() == 0 && init.Type().NumOut() == 0 {
		out = init.Call(in)
	}

	// Validate
	vali := sc.MethodByName("Validate")
	if vali.IsValid() && init.Type().NumIn() == 0 && vali.Type().NumOut() == 1 {
		out = vali.Call(in)
	}

	// Run
	out = sc.MethodByName("Run").Call(in)
	cmd.report()

	return out[0].Interface().(error)
}

func (cmd *Command) report() {
	if !cmd.verbose {
		return
	}
}

func (cmd *Command) helpCommand() error {
	// Sub Command Help
	if len(cmd.subCmd) > 0 {
		return cmd.helpSubCommand()
	}
	// Command Help
	sb := strings.Builder{}
	sb.WriteString("Usage:\n")
	sb.WriteString("    " + cmd.Name + " [options] <sub command> [sub command options] [arg] ...\n")
	sb.WriteString("\nOptions:\n")
	// Option
	max := 0
	for _, o := range cmd.opts {
		max = int(math.Max(float64(max), float64(len(o.Long))))
	}
	format := fmt.Sprintf("    %%-2s, %%-%ds   %%s\n", max)
	for _, o := range cmd.opts {
		sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
	}
	// Sub Command
	max += 4
	sb.WriteString("\nSub Commands:\n")
	for _, sc := range cmd.scOrder {
		max = int(math.Max(float64(max), float64(len(sc))))
	}
	format = fmt.Sprintf("    %%-%ds   %%s\n", max)
	for _, sc := range cmd.scOrder {
		sb.WriteString(fmt.Sprintf(format, sc, cmd.subCmds[sc].Desc))
	}
	// Sub Command Options
	sb.WriteString("\nSub Command Options:\n")
	sb.WriteString("    execute the following command\n")
	sb.WriteString(fmt.Sprintf("\n    %s -h <SubCommand>\n", cmd.Name))

	println(sb.String())
	return nil
}

func (cmd *Command) helpSubCommand() error {
	sc := cmd.subCmds[cmd.subCmd]

	sb := strings.Builder{}
	sb.WriteString("Usage:\n")
	sb.WriteString("    " + cmd.Name + " " + cmd.subCmd + " [options] [arg] ...\n")
	if len(sc.opts) > 0 {
		sb.WriteString("Options:\n")
		max := 0
		for _, o := range sc.opts {
			max = int(math.Max(float64(max), float64(len(o.Name()))))
		}
		format := fmt.Sprintf("    %%-2s, %%-%ds   %%s\n", max)
		for _, o := range sc.opts {
			sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
		}
	}

	println(sb.String())
	return nil
}

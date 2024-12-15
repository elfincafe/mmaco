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
	if !exists && name != helpCmdName {
		cmd.scOrder = append(cmd.scOrder, name)
	}
}

func (cmd *Command) route(args []string) ([]string, error) {
	subArgs := []string{}
	// parse root options
	lastIdx := len(args) - 1
	subCmdIdx := lastIdx
	skip := false
	subCmdChain := "@" + strings.Join(cmd.scOrder, "@") + "@"
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		// Sub Command
		if strings.Contains(subCmdChain, "@"+arg+"@") {
			cmd.subCmd = arg
			subCmdIdx = i
			break
		}
		// Root Options
		for _, opt := range cmd.opts {
			if arg == "-"+opt.Short && opt.Kind() == Bool {
				opt.set("true")
				break
			} else if arg == "-"+opt.Short && opt.Kind() != Unknown {
				if lastIdx > i {
					skip = true
					break
				} else {
					return subArgs, fmt.Errorf(`option "%s" needs a value`, opt.Name())
				}
			} else if arg == "--"+opt.Long && opt.Kind() == Bool {
				opt.set("true")
				break
			} else if strings.HasPrefix(arg, "--"+opt.Long+"=") {
				length := len("--" + opt.Long + "=")
				opt.set(arg[length:])
				break
			}
		}
		return subArgs, fmt.Errorf(`unkown options or sub command: %s`, arg)
	}
	// Sub Command Options
	if len(args)-1 > subCmdIdx {
		for i, arg := range args[subCmdIdx+1:] {
			subCmd := cmd.subCmds[cmd.subCmd]
			for _, opt := range subCmd.opts {
				if arg == "-"+opt.Short && opt.Kind() == Bool {
					opt.set("true")
					break
				} else if arg == "-"+opt.Short && opt.Kind() != Unknown {
					if lastIdx > i {
						skip = true
						break
					} else {
						return subArgs, fmt.Errorf(`option "%s" needs a value`, opt.Name())
					}
				} else if arg == "--"+opt.Long && opt.Kind() == Bool {
					opt.set("true")
					break
				} else if strings.HasPrefix(arg, "--"+opt.Long+"=") {
					length := len("--" + opt.Long + "=")
					opt.set(arg[length:])
					break
				}
			}
			subArgs = append(subArgs, arg)
		}
	}

	return subArgs, nil
}

func (cmd *Command) Run() error {
	in, out := []reflect.Value{}, []reflect.Value{}

	// Routing
	subArgs, err := cmd.route(os.Args[1:])
	if err != nil {
		return err
	}
	fmt.Println(subArgs, cmd.subCmd)

	// Analizing
	sc := reflect.ValueOf(cmd.subCmds[cmd.subCmd])

	// Initialize
	init := sc.MethodByName("Init")
	if init.IsValid() && init.Type().NumIn() == 0 && init.Type().NumOut() == 0 {
		init.Call([]reflect.Value{})
	}

	// Argument Handler
	for _, opt := range cmd.opts {
		if opt.Handler == "" {
			continue
		}
		h := sc.MethodByName(opt.Handler)
		if h.IsValid() && h.Type().NumIn() == 1 && h.Type().Kind() == reflect.String && init.Type().NumOut() == 1 {
			out = h.Call(in)
			if out[0].CanInterface() {
				return out[0].Interface().(error)
			} else {
				return fmt.Errorf(`"%s" method should return error`, opt.Handler)
			}
		} else {
			return fmt.Errorf(`"%s" method should return error`, opt.Handler)
		}
	}

	// Validate
	vali := sc.MethodByName("Validate")
	if vali.IsValid() && vali.Type().NumIn() == 0 && vali.Type().NumOut() == 1 {
		out = vali.Call(in)
		if out[0].CanInterface() {
			return out[0].Interface().(error)
		} else {
			return fmt.Errorf(`%s.Validate method should return error`)
		}
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

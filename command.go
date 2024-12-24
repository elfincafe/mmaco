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
	cmd.loc = nil
	cmd.start = time.Now()
	cmd.Name = name
	cmd.subCmds = map[string]*subCommand{}
	cmd.scOrder = []string{}
	cmd.opts = []*option{}
	return cmd
}

func (cmd *Command) SetLocation(loc *time.Location) {
	cmd.loc = loc
}

func (cmd *Command) parse() error {
	v := reflect.ValueOf(cmd).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		opt := newOption(v.Field(i), t.Field(i), cmd.loc)
		if opt == nil {
			continue
		}
		cmd.opts = append(cmd.opts, opt)
	}
	return nil
}

func (cmd *Command) Add(subCmd SubCommandInterface) error {
	return cmd.addSubCmd(subCmd, false)
}

func (cmd *Command) addSubCmd(subCmd SubCommandInterface, force bool) error {
	sc := newSubCommand(subCmd, cmd.loc)
	if sc == nil {
		return fmt.Errorf("an invalid command")
	}
	err := sc.parse()
	if err != nil {
		return err
	}
	name := sc.Name
	if name == helpCmdName && !force {
		return fmt.Errorf(`reserved command: help`)
	}
	cmd.subCmds[name] = sc
	cmd.scOrder = append(cmd.scOrder, name)
	return nil
}

func (cmd *Command) route(args []string) (int, error) {
	idxs := []int{}
	for i, arg := range args {
		for _, sc := range cmd.scOrder {
			if arg == sc {
				idxs = append(idxs, i)
			}
		}
	}
	for _, idx := range idxs {
		if idx == 0 {
			return idx, nil
		}
		preArg := args[idx-1]
		for _, o := range cmd.opts {
			if preArg == "--"+o.Long && o.Kind == Bool {
				return idx, nil
			} else if preArg == "-"+o.Short && o.Kind == Bool {
				return idx, nil
			} else if strings.HasPrefix(preArg, "--"+o.Long+"=") {
				return idx, nil
			}
		}
	}
	if len(idxs) > 0 {
		return idxs[0], nil
	} else {
		return -1, fmt.Errorf("sub command isn't passed")
	}
}

func (cmd *Command) parseArgs(args []string) error {
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			cmd.help = true
		case "-r", "--report":
			cmd.report = true
		default:
			return fmt.Errorf(`unknown option: %s`, arg)
		}
	}
	return nil
}

func (cmd *Command) Run() error {
	var err error
	in, out := []reflect.Value{}, []reflect.Value{}

	// Add Help Command
	err = cmd.addSubCmd(new(help), true)
	if err != nil {
		return err
	}

	// Parse Options
	err = cmd.parse()
	if err != nil {
		return err
	}

	// Routing
	args := os.Args[1:]
	subCmdIdx, err := cmd.route(args)
	if err != nil {
		return err
	}
	subCmd := args[subCmdIdx]
	sc := cmd.subCmds[subCmd]

	// Parse Argument for Root
	err = cmd.parseArgs(args[:subCmdIdx])
	if err != nil {
		return err
	}
	// Parse Argument for Sub Command
	params, err := sc.parseArgs(args[subCmdIdx+1:])
	if err != nil {
		return err
	}

	// Validate
	if sc.hasValidate {
		out = sc.cmd.MethodByName("Validate").Call(nil)
		if !out[0].IsNil() {
			return out[0].Interface().(error)
		}
	}

	// Run
	cmd.startSubCmd = time.Now()
	out = sc.cmd.MethodByName("Run").Call(append(in, reflect.ValueOf(params)))
	cmd.endSubCmd = time.Now()

	fmt.Println(time.Now().Sub(cmd.start))
	fmt.Println(cmd.endSubCmd.Sub(cmd.startSubCmd))

	if out[0].IsNil() {
		return nil
	} else {
		return out[0].Interface().(error)
	}
}

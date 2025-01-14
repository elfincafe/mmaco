package mmaco

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type (
	Command struct {
		Name   string
		ctx    *Context
		opts   []*option
		report bool `mmaco:"short=r,long=report,desc=report verbosely."`
		help   bool `mmaco:"short=h,long=help,desc=this help."`
	}
)

func New(name string) *Command {
	ctx := newContext(os.Args[1:])
	ctx.cmd = new(Command)
	ctx.cmd.ctx = ctx
	ctx.cmd.Name = name
	ctx.cmd.opts = []*option{}
	return ctx.cmd
}

func (cmd *Command) SetLocation(loc *time.Location) {
	cmd.ctx.loc = loc
}

func (cmd *Command) parse() error {
	v := reflect.ValueOf(cmd).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		opt := newOption(v.Field(i), t.Field(i), cmd.ctx.loc)
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
	sc := newSubCommand(subCmd, cmd.ctx.loc)
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
	cmd.ctx.subCmds[name] = sc
	cmd.ctx.scOrder = append(cmd.ctx.scOrder, name)
	return nil
}

func (cmd *Command) route(args []string) int {
	idx := -1
	for i, arg := range args {
		if arg == "-r" || arg == "--report" {
			cmd.report = true
		} else if arg == "-h" || arg == "--help" {
			cmd.help = true
		} else {
			idx = i
			break
		}
	}
	return idx
}

// func (cmd *Command) parseArgs(args []string) error {
// 	for _, arg := range args {
// 		switch arg {
// 		case "-h", "--help":
// 			cmd.help = true
// 		case "-r", "--report":
// 			cmd.report = true
// 		default:
// 			return fmt.Errorf(`unknown option: %s`, arg)
// 		}
// 	}
// 	return nil
// }

func (cmd *Command) showReport(ctx *Context) {
	buf := strings.Builder{}
	buf.WriteString("--------------------------------------------------\n")
	buf.WriteString(" MMaco CLI Framework \n")
	buf.WriteString("--------------------------------------------------\n")
	fmt.Println(buf.String())
}

func (cmd *Command) Run() error {
	var err error
	in, out := []reflect.Value{}, []reflect.Value{}

	// Add Help Command
	err = cmd.addSubCmd(help{}, true)
	if err != nil {
		return err
	}

	// Parse Options
	err = cmd.parse()
	if err != nil {
		return err
	}

	// Routing
	rowArgs := cmd.ctx.RowArgs()
	subCmdIdx := cmd.route(rowArgs)
	subCmd := ""
	if cmd.help { // passed -h or --help option.
		println(1)
		subCmd = helpCmdName
	} else if subCmdIdx < 0 { // passed no sub command.
		println(2)
		subCmd = helpCmdName
	} else {
		println(3)
		subCmd = cmd.ctx.RowArg(subCmdIdx)
	}
	sc := cmd.ctx.subCmds[subCmd]

	// Parse Argument for Sub Command
	cmd.ctx.args, err = sc.parseArgs(rowArgs[subCmdIdx+1:])
	if err != nil {
		return err
	}

	// Validate
	out = sc.cmd.MethodByName("Validate").Call(nil)
	if !out[0].IsNil() {
		return out[0].Interface().(error)
	}

	// Run
	cmd.ctx.subCmdStart = time.Now().UnixMicro()
	out = sc.cmd.MethodByName("Run").Call(append(in, reflect.ValueOf(cmd.ctx)))
	cmd.ctx.subCmdFinish = time.Now().UnixMicro()

	// Report
	if cmd.report && subCmd != helpCmdName {
		cmd.showReport(cmd.ctx)
	}

	if out[0].IsNil() {
		return nil
	} else {
		return out[0].Interface().(error)
	}
}

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
	for i := 0; i < v.Type().NumField(); i++ {
		opt := newOption(v.Field(i), v.Type().Field(i))
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
	sc := newSubCommand(subCmd)
	if sc.Name == helpCmdName && !force {
		return fmt.Errorf(`reserved command: help`)
	}
	sc.ctx = cmd.ctx
	cmd.ctx.subCmds[sc.Name] = sc
	cmd.ctx.scOrder = append(cmd.ctx.scOrder, sc.Name)
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

func (cmd *Command) showReport(ctx *Context) {
	subCmdTime := time.UnixMicro(cmd.ctx.subCmdFinish).Sub(time.UnixMicro(cmd.ctx.subCmdStart))
	cmdTime := time.Now().Sub(time.UnixMicro(cmd.ctx.cmdStart))
	buf := strings.Builder{}
	buf.WriteString("\n")
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(" MMaco CLI Framework \n")
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf(" DateTime:   %v\n", time.UnixMicro(cmd.ctx.cmdStart).In(cmd.ctx.loc)))
	buf.WriteString(fmt.Sprintf(" Command:    %v\n", cmdTime))
	buf.WriteString(fmt.Sprintf(" SubCommand: %v\n", subCmdTime))
	buf.WriteString("------------------------------------------------------------\n")
	println(buf.String())
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
	rowArgs := cmd.ctx.RowArgs()
	subCmdIdx := cmd.route(rowArgs)
	subCmdName := ""
	if cmd.help { // passed -h or --help option.
		subCmdName = helpCmdName
	} else if subCmdIdx < 0 { // passed no sub command.
		subCmdName = helpCmdName
	} else {
		subCmdName = cmd.ctx.RowArg(subCmdIdx)
	}

	// Copy Sub Command
	sc := new(SubCommand)
	sc.Name = cmd.ctx.subCmds[subCmdName].Name
	sc.Desc = cmd.ctx.subCmds[subCmdName].Desc
	sc.cmd = cmd.ctx.subCmds[subCmdName].cmd

	err = sc.parse()
	if err != nil {
		return err
	}

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
	if cmd.report && subCmdName != helpCmdName {
		cmd.showReport(cmd.ctx)
	}

	if out[0].IsNil() {
		return nil
	} else {
		return out[0].Interface().(error)
	}
}

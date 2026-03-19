package mmaco

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type (
	Command struct {
		ctx        *Context
		name       string
		opts       []*option
		subcmdRule *regexp.Regexp
		debug      bool `mmaco:"long=debug,desc=run as debug mode"`
		report     bool `mmaco:"long=report,desc=report when command is finished without error"`
		help       bool `mmaco:"short=h,long=help,desc=this help"`
	}
)

func New(name string) *Command {
	// Rules (defined in mmaco.go)
	ruleShortOpt = regexp.MustCompile(`^[0-9a-zA-Z]$`)
	ruleLongOpt = regexp.MustCompile(`^[0-9a-zA-Z\-]{2,10}$`)

	ctx := newContext(name, os.Args[1:])
	cmd := new(Command)
	ctx.cmd = cmd
	cmd.ctx = ctx
	cmd.name = name
	cmd.opts = []*option{}
	cmd.subcmdRule = regexp.MustCompile(`^[a-z][\da-z_\-:]*[\da-z]$`)
	cmd.debug = false
	cmd.report = false
	cmd.help = false
	return cmd
}

func (cmd *Command) SetLocation(loc *time.Location) {
	cmd.ctx.loc = loc
}

func (cmd *Command) parse() {
	v := reflect.ValueOf(cmd).Elem()
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
		opt := newOption(f, ft, cmd.ctx)
		cmd.opts = append(cmd.opts, opt)
	}
}

func (cmd *Command) Add(subCmd SubCommand, name, desc string) error {
	if !cmd.subcmdRule.MatchString(name) {
		return fmt.Errorf("sub command name is wrong (%s)", cmd.subcmdRule.String())
	}
	t := reflect.TypeOf(subCmd)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf(`pass a pointer type that implements the SubCommand interface`)
	}
	p := t.Elem()
	if p.Kind() == reflect.Pointer {
		return fmt.Errorf("double pointers are not allowed")
	}
	return cmd.addSubCmd(subCmd, name, desc, false)
}

func (cmd *Command) addSubCmd(subCmd SubCommand, name, desc string, force bool) error {
	bucket := newSubCommand(subCmd, name, desc)
	if bucket.Name == helpCmdName && !force {
		return fmt.Errorf(`"help" is reserved`)
	}
	bucket.ctx = cmd.ctx
	if _, ok := cmd.ctx.subCmds[bucket.Name]; !ok {
		cmd.ctx.scOrder = append(cmd.ctx.scOrder, bucket.Name)
	}
	cmd.ctx.subCmds[bucket.Name] = bucket
	return nil
}

func (cmd *Command) route(args []string) (int, error) {
	idx := -1
	arg := ""
	for idx, arg = range args {
		flg := false
		for _, o := range cmd.opts {
			if o.isShort(arg) || o.isLong(arg) {
				switch o.Name {
				case "debug":
					flg = true
					cmd.debug = true
				case "report":
					flg = true
					cmd.report = true
				case "help":
					flg = true
					cmd.help = true
				}
			}
		}
		if !flg && strings.HasPrefix(arg, "-") {
			return -1, fmt.Errorf(`unknown option "%s"`, arg)
		} else if !flg {
			return idx, nil
		}
	}
	return -1, nil
}

func (cmd *Command) showReport(ctx *Context) {
	subCmdTime := time.UnixMicro(cmd.ctx.subCmdFinish).Sub(time.UnixMicro(cmd.ctx.subCmdStart))
	cmdTime := time.Since(time.UnixMicro(cmd.ctx.cmdStart))
	idx := -1
	for k, v := range ctx.rawArgs {
		if v == ctx.subCmd.Name {
			idx = k
			break
		}
	}
	gOpts := ctx.rawArgs[:idx]
	sOpts := ctx.rawArgs[idx+1:]
	buf := strings.Builder{}
	buf.WriteString("\n")
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf(" %v Command\n", ctx.cmd.name))
	buf.WriteString("------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf(" Options:    %v\n", strings.Join(gOpts, " ")))
	buf.WriteString(fmt.Sprintf(" SubCommand: %v\n", ctx.subCmd.Name))
	buf.WriteString(fmt.Sprintf(" SubOptions: %v\n", strings.Join(sOpts, " ")))
	buf.WriteString(fmt.Sprintf(" DateTime:   %v\n", time.UnixMicro(cmd.ctx.cmdStart).In(cmd.ctx.loc)))
	buf.WriteString(fmt.Sprintf(" ExecTime:   %v\n", cmdTime))
	buf.WriteString(fmt.Sprintf(" SubTime:    %v\n", subCmdTime))
	buf.WriteString("------------------------------------------------------------\n")
	println(buf.String())
}

func (cmd *Command) Run() error {
	var err error

	// Add Help Command
	err = cmd.addSubCmd(new(help), "help", "this help", true)
	if err != nil {
		return err
	}

	// Parse Options
	cmd.parse()

	// Routing
	rowArgs := cmd.ctx.RawArgs()
	subCmdIdx, err := cmd.route(rowArgs)
	if err != nil {
		return err
	}
	subCmdName := ""
	if cmd.help { // passed -h or --help option.
		subCmdName = helpCmdName
	} else if subCmdIdx < 0 { // passed no sub command.
		subCmdName = helpCmdName
	} else {
		subCmdName = cmd.ctx.RawArg(subCmdIdx)
	}

	// Debug Mode
	debugMode = cmd.debug

	// Copy Sub Command
	if _, ok := cmd.ctx.subCmds[subCmdName]; !ok {
		return fmt.Errorf(`unknown command "%s" is specified`, subCmdName)
	}
	cmd.ctx.subCmd = cmd.ctx.subCmds[subCmdName]

	// Init
	cmd.ctx.subCmd.cmd.Init()

	// parse arguments for global options
	cmd.ctx.subCmd.parse()

	// Parse Argument for Sub Command options
	cmd.ctx.args, err = cmd.ctx.subCmd.parseArgs(rowArgs[subCmdIdx+1:])
	if err != nil {
		return err
	}

	// Run
	cmd.ctx.subCmdStart = time.Now().UnixMicro()
	err = cmd.ctx.subCmd.cmd.Run(cmd.ctx)
	cmd.ctx.subCmdFinish = time.Now().UnixMicro()

	// Report
	if cmd.report && subCmdName != helpCmdName {
		cmd.showReport(cmd.ctx)
	}

	return err
}

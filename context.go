package mmaco

import "time"

type (
	Context struct {
		cmd          *Command
		subCmd       *SubCommand
		subCmds      map[string]*SubCommand
		scOrder      []string
		loc          *time.Location
		cmdStart     int64
		subCmdStart  int64
		subCmdFinish int64
		rawArgs      []string
		args         []string
	}
)

func newContext(cmdName string, rawArgs []string) *Context {
	ctx := new(Context)
	ctx.cmd = nil
	ctx.subCmd = nil
	ctx.subCmds = map[string]*SubCommand{}
	ctx.scOrder = []string{}
	ctx.loc, _ = time.LoadLocation("")
	ctx.cmdStart = time.Now().UnixMicro()
	ctx.subCmdStart = int64(0)
	ctx.subCmdFinish = int64(0)
	ctx.rawArgs = rawArgs
	ctx.args = []string{}
	return ctx
}

func (ctx *Context) Location() *time.Location {
	return ctx.loc
}

func (ctx *Context) StartTime(command bool) time.Time {
	if command {
		return time.UnixMicro(ctx.cmdStart)
	} else {
		return time.UnixMicro(ctx.subCmdStart)
	}
}

func (ctx *Context) Arg(i int) string {
	if i >= 0 && i < len(ctx.args) {
		return ctx.args[i]
	}
	return ""
}

func (ctx *Context) Args() []string {
	return ctx.args
}

func (ctx *Context) NumArg() int {
	return len(ctx.args)
}

func (ctx *Context) RawArg(i int) string {
	if i >= 0 && i < len(ctx.rawArgs) {
		return ctx.rawArgs[i]
	}
	return ""
}

func (ctx *Context) RawArgs() []string {
	return ctx.rawArgs
}

func (ctx *Context) RawNumArg() int {
	return len(ctx.rawArgs)
}

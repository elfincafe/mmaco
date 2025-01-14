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
		rowArgs      []string
		args         []string
	}
)

func newContext(rowArgs []string) *Context {
	ctx := new(Context)
	ctx.cmd = nil
	ctx.subCmd = nil
	ctx.subCmds = map[string]*SubCommand{}
	ctx.scOrder = []string{}
	ctx.loc, _ = time.LoadLocation("")
	ctx.cmdStart = time.Now().UnixMicro()
	ctx.subCmdStart = int64(0)
	ctx.subCmdFinish = int64(0)
	ctx.rowArgs = rowArgs
	ctx.args = []string{}
	return ctx
}

func (ctx *Context) Command() *Command {
	return ctx.cmd
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

func (ctx *Context) RowArg(i int) string {
	if i >= 0 && i < len(ctx.rowArgs) {
		return ctx.rowArgs[i]
	}
	return ""
}

func (ctx *Context) RowArgs() []string {
	return ctx.rowArgs
}

func (ctx *Context) RowNumArg() int {
	return len(ctx.rowArgs)
}

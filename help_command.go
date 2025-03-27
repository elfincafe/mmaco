package mmaco

import (
	"fmt"
	"math"
	"strings"
)

type (
	help struct {
	}
)

func (cmd *help) Init() {
}

func (cmd *help) Validate() error {
	return nil
}

func (cmd *help) Run(ctx *Context) error {
	var sb strings.Builder
	subCmdName := ctx.Arg(0)
	if _, ok := ctx.subCmds[subCmdName]; ok {
		ctx.subCmd = ctx.subCmds[subCmdName]
		ctx.subCmd.parse()
		sb = cmd.getSubCommandHelp(ctx)
	} else {
		sb = cmd.getCommandHelp(ctx)
	}
	println(sb.String())
	return nil
}

func (cmd *help) getCommandHelp(ctx *Context) strings.Builder {
	line := ""
	sb := strings.Builder{}

	sb.WriteString("\nUsage:\n")
	line = fmt.Sprintf("    %s [options] <sub command> [sub command options] [arg] ...\n", ctx.cmd.name)
	sb.WriteString(line)
	sb.WriteString("\nOptions:\n")

	// option width
	width := 0
	for _, o := range ctx.cmd.opts {
		width = int(math.Max(float64(len(o.Long)+6), float64(width)))
	}
	for _, sc := range ctx.scOrder {
		width = int(math.Max(float64(len(sc)), float64(width)))
	}

	// Global Options
	spaces := ""
	for _, o := range ctx.cmd.opts {
		if o.Short != "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-6-len(o.Long))
			line = fmt.Sprintf(`    -%s, --%s%s   %s`, o.Short, o.Long, spaces, o.Desc)
		} else if o.Short != "" && o.Long == "" {
			spaces = strings.Repeat(" ", width-2)
			line = fmt.Sprintf(`    -%s%s   %s`, o.Short, spaces, o.Desc)
		} else if o.Short == "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-(len(o.Long)+6))
			line = fmt.Sprintf(`        --%s%s   %s`, o.Long, spaces, o.Desc)
		}
		sb.WriteString(line + "\n")
	}

	// Sub Commands
	sb.WriteString("\nSub Commands:\n")
	for _, sc := range ctx.scOrder {
		spaces = strings.Repeat(" ", width-len(sc))
		line = fmt.Sprintf(`    %s%s   %s`, sc, spaces, ctx.subCmds[sc].Desc)
		sb.WriteString(line + "\n")
	}

	return sb
}

func (cmd *help) getSubCommandHelp(ctx *Context) strings.Builder {
	line := ""
	sb := strings.Builder{}

	sb.WriteString("\nUsage:\n")
	line = fmt.Sprintf("    %s [options] %s [sub command options] [arg] ...\n", ctx.cmd.name, ctx.subCmd.Name)
	sb.WriteString(line)
	sb.WriteString("\nOptions:\n")

	// option width
	width := 0
	for _, o := range ctx.cmd.opts {
		width = int(math.Max(float64(len(o.Long)+6), float64(width)))
	}
	for _, o := range ctx.subCmd.opts {
		width = int(math.Max(float64(len(o.Long)+6), float64(width)))
	}

	// Global Options
	spaces := ""
	for _, o := range ctx.cmd.opts {
		if o.Short != "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-6-len(o.Long))
			line = fmt.Sprintf(`    -%s, --%s%s   %s`, o.Short, o.Long, spaces, o.Desc)
		} else if o.Short != "" && o.Long == "" {
			spaces = strings.Repeat(" ", width-2)
			line = fmt.Sprintf(`    -%s%s   %s`, o.Short, spaces, o.Desc)
		} else if o.Short == "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-(len(o.Long)+6))
			line = fmt.Sprintf(`        --%s%s   %s`, o.Long, spaces, o.Desc)
		}
		sb.WriteString(line + "\n")
	}
	if len(ctx.subCmd.opts) == 0 {
		return sb
	}

	// Sub Command Options
	sb.WriteString("\nSub Command Options:\n")
	for _, o := range ctx.subCmd.opts {
		if o.Short != "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-6-len(o.Long))
			line = fmt.Sprintf(`    -%s, --%s%s   %s`, o.Short, o.Long, spaces, o.Desc)
		} else if o.Short != "" && o.Long == "" {
			spaces = strings.Repeat(" ", width-2)
			line = fmt.Sprintf(`    -%s%s   %s`, o.Short, spaces, o.Desc)
		} else if o.Short == "" && o.Long != "" {
			spaces = strings.Repeat(" ", width-(len(o.Long)+6))
			line = fmt.Sprintf(`        --%s%s   %s`, o.Long, spaces, o.Desc)
		}
		sb.WriteString(line + "\n")
	}

	return sb
}

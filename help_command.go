package mmaco

type (
	help struct {
		Desc string
	}
)

func (h help) Init() {
	println("Exec: Help.Init")
}

func (h help) Validate() error {
	println("Exec: Help.Validate")
	return nil
}

func (h help) Run(ctx *Context) error {
	println("Exec: Help.Run")
	// Sub Command Help
	// if len(cmd.subCmd) > 0 {
	// 	return cmd.helpSubCommand()
	// }
	// // Command Help
	// sb := strings.Builder{}
	// sb.WriteString("Usage:\n")
	// sb.WriteString("    " + cmd.Name + " [options] <sub command> [sub command options] [arg] ...\n")
	// sb.WriteString("\nOptions:\n")
	// // Option
	// max := 0
	// for _, o := range cmd.opts {
	// 	max = int(math.Max(float64(max), float64(len(o.Long))))
	// }
	// format := fmt.Sprintf("    %%-2s, %%-%ds   %%s\n", max)
	// for _, o := range cmd.opts {
	// 	sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
	// }
	// // Sub Command
	// max += 4
	// sb.WriteString("\nSub Commands:\n")
	// for _, sc := range cmd.scOrder {
	// 	max = int(math.Max(float64(max), float64(len(sc))))
	// }
	// format = fmt.Sprintf("    %%-%ds   %%s\n", max)
	// for _, sc := range cmd.scOrder {
	// 	sb.WriteString(fmt.Sprintf(format, sc, cmd.subCmds[sc].Desc))
	// }
	// // Sub Command Options
	// sb.WriteString("\nSub Command Options:\n")
	// sb.WriteString("    execute the following command\n")
	// sb.WriteString(fmt.Sprintf("\n    %s -h <SubCommand>\n", cmd.Name))

	// println(sb.String())
	return nil
}

func (h help) helpSubCommand() error {
	// sc := cmd.subCmds[cmd.subCmd]

	// sb := strings.Builder{}
	// sb.WriteString("Usage:\n")
	// sb.WriteString("    " + cmd.Name + " " + cmd.subCmd + " [options] [arg] ...\n")
	// if len(sc.opts) > 0 {
	// 	sb.WriteString("Options:\n")
	// 	max := 0
	// 	for _, o := range sc.opts {
	// 		max = int(math.Max(float64(max), float64(len(o.Name))))
	// 	}
	// 	format := fmt.Sprintf("    %%-2s, %%-%ds   %%s\n", max)
	// 	for _, o := range sc.opts {
	// 		sb.WriteString(fmt.Sprintf(format, o.Short, o.Long, o.Desc))
	// 	}
	// }

	// println(sb.String())
	return nil
}

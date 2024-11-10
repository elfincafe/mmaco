package mmaco

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

type (
	Command struct {
		Name    string
		subCmds map[string]*subCommand
		subCmd  string
		scOrder []string
		start   time.Time
		help    bool `mmaco:"short=h,long=help"`
		verbose bool `mmaco:"short=v,long=verbose"`
	}
)

func New(name string) Command {
	cmd := Command{
		start:   time.Now(),
		Name:    name,
		subCmds: map[string]*subCommand{},
		subCmd:  "",
		scOrder: []string{},
	}
	return cmd
}

func (cmd Command) AddSubCommand(subCmd SubCommandInterface) error {
	sc := newSubCommand(subCmd)
	err := sc.validate()
	if err != nil {
		return err
	}
	cmd.subCmds[sc.Name()] = sc
	return nil
}

func (cmd Command) route(args []string) int {
	metas := getMetas(reflect.TypeOf(cmd))
	fmt.Println(metas["help"], metas["verbose"])
	return 0
}

func (cmd Command) Run() error {
	// Routing
	subCmdPos := cmd.route(os.Args[1:])
	fmt.Println(subCmdPos, cmd.subCmd)

	// Analizing
	sc := reflect.ValueOf(cmd.subCmd)

	// Intialize
	init := sc.MethodByName("Init")
	if init.IsValid() {
		init.Call([]reflect.Value{})
	}

	// Parsing Arguments

	// Run
	sc.MethodByName("Run").Call([]reflect.Value{})

	return nil
}

func (cmd Command) Report() {
	if !cmd.verbose {
		return
	}
}

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
		Name    string
		subCmds map[string]subCommand
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
		subCmds: map[string]subCommand{},
		subCmd:  "",
		scOrder: []string{},
	}
	return cmd
}

func (cmd Command) AddSubCommand(name string, subCommand any) {
	sc := newSubCommand(subCommand)
	name = cmd.convSubCommandName(name)
	cmd.subCmds[name] = sc
}

func (cmd Command) convSubCommandName(name string) string {
	buf := []byte{}
	name = strings.TrimSpace(strings.ToLower(name))
	length := len(name)
	lastIndex := length - 1
	for i, b := range []byte(name) {
		if (b < 97 || b > 122) && (b < 48 || b > 57) && b != 95 { // a-z, 0-9, -, _
			continue
		}
		if i == 0 && (b < 97 || b > 122) {
			// Initial character allows only a to z
			continue
		} else if i == lastIndex && b == 95 {
			// Last character allows a to z, 0 to 9
			continue
		}
		buf = append(buf, b)
	}
	return string(buf)
}

func (cmd Command) route() int {
	pos := len(os.Args[1:])

	props := getMetas(reflect.TypeOf(cmd))
	fmt.Println(props)

	return pos
}

func (cmd Command) Run() error {
	// Routing
	subCmdPos := cmd.route()
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

package mmaco

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type (
	Command struct {
		Name    string
		subCmds map[string]*subCommand
		subCmd  string
		scOrder []string
		start   time.Time
		Help    bool `mmaco:"short=h,long=help"`
		Verbose bool `mmaco:"short=v,long=verbose"`
	}
)

func New(name string) *Command {
	cmd := new(Command)
	cmd.start = time.Now()
	cmd.Name = name
	cmd.subCmds = map[string]*subCommand{}
	cmd.subCmd = ""
	cmd.scOrder = []string{}
	return cmd
}

func (cmd *Command) Add(name string, subCmd SubCommandInterface) error {
	sc := newSubCommand(subCmd)
	re := regexp.MustCompile(`^[a-z][0-9a-z_\-]*[0-9a-z]?$`)
	if !re.MatchString(name) {
		return fmt.Errorf("the sub command does not follow the format (^[a-z][0-9a-z_-]*[0-9a-z]?$)")
	}
	cmd.subCmds[name] = sc
	return nil
}

func (cmd *Command) route(args []string) error {
	// var err error

	// c := reflect.ValueOf(cmd)
	// metas := getMetas(c.Type())

	// Root Options
	opts := []string{"-h", "--help", "-v", "--verbose"}
	idx := 0
	for i, arg := range args {
		ok := false
		for _, opt := range opts {
			if arg == opt {
				ok = true
				break
			}
		}
		if ok {
			idx = i + 1
		}
	}
	fmt.Println(args[:idx], idx, len(args[idx:]))

	// SubCommand
	if len(args[idx:]) > 0 {
		ok := false
		for _, subcmd := range cmd.scOrder {
			if args[idx] == subcmd {
				cmd.subCmd = subcmd
				ok = true
				break
			}
			if ok {
				break
			}
		}
	} else {
		return fmt.Errorf("SubCommand isn't passed")
	}
	idx += 1
	fmt.Println(cmd.subCmd, args[idx:])

	// skip := false
	// for i, arg := range args {
	// 	if skip {
	// 		continue
	// 	}
	// 	for name, meta := range metas {
	// 		field := c.Elem().FieldByName(name)
	// 		kind := field.Kind()
	// 		short := "-" + meta.short
	// 		long := "--" + meta.long
	// 		if arg == long {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if kind == reflect.Bool {
	// 				cmd.setArg(&field, long, "true")
	// 				break
	// 			} else {
	// 				err = fmt.Errorf("needs value for the '%s' (e.g. --%s=something)", name, meta.long)
	// 			}
	// 		} else if arg == short {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if kind == reflect.Bool {
	// 				cmd.setArg(&field, short, "true")
	// 				break
	// 			}
	// 			if len(args) >= i && !strings.HasPrefix(args[i+1], "-") {
	// 				skip = true
	// 				break
	// 			}
	// 		} else if strings.HasPrefix(arg, long) {
	// 			if !field.CanSet() {
	// 				err = fmt.Errorf("can't set to the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			if !strings.HasPrefix(arg, long+"=") {
	// 				err = fmt.Errorf("needs value for the field '%s'", name)
	// 				goto ERROR
	// 			}
	// 			n := len("--" + meta.long + "=")
	// 			cmd.setArg(&field, long, arg[n:])
	// 			break
	// 		}
	// 	}
	// 	idx = i - 1
	// }

	// fmt.Println(metas["help"], metas["verbose"])
	return nil
	// ERROR:
	// 	return err
}

func (cmd *Command) setArg(field *reflect.Value, opt, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(true)
	case reflect.Int:
		v, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the int type", opt)
		}
		field.SetInt(v)
	case reflect.Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the int8 type", opt)
		}
		field.SetInt(v)
	case reflect.Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the int16 type", opt)
		}
		field.SetInt(v)
	case reflect.Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the int32 type", opt)
		}
		field.SetInt(v)
	case reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the int64 type", opt)
		}
		field.SetInt(v)
	case reflect.Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the uint type", opt)
		}
		field.SetUint(v)
	case reflect.Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the uint8 type", opt)
		}
		field.SetUint(v)
	case reflect.Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the uint16 type", opt)
		}
		field.SetUint(v)
	case reflect.Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the uint32 type", opt)
		}
		field.SetUint(v)
	case reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the uint64 type", opt)
		}
		field.SetUint(v)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the float32 type", opt)
		}
		field.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("The value of option '%s' should be the float64 type", opt)
		}
		field.SetFloat(v)
	case reflect.String:
		field.SetString(value)
	default:
		return fmt.Errorf("The field type of '%s' isn't supported", field.Type().Name())
	}
	return nil
}

func (cmd *Command) Run() error {
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

func (cmd *Command) Report() {
	if !cmd.Verbose {
		return
	}
}

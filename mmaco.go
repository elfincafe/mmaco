package mmaco

import (
	"bytes"
	"reflect"
	"regexp"
	"time"
)

const (
	tagName   = "mmaco"
	trimSpace = " \t\v\r\n\f"

	Unknown Kind = iota
	String
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Time
)

type (
	Kind int

	Command struct {
		Name    string
		subCmds map[string]*subCommand
		subCmd  string
		scOrder []string
		opts    []*option
		start   time.Time
		help    bool `mmaco:"short=h,long=help,default=false"`
		verbose bool `mmaco:"short=v,long=verbose,default=false"`
	}

	SubCommandInterface interface {
		Init() error
		Run(...[]string) error
	}

	subCommand struct {
		cmd  reflect.Value
		opts []*option
	}

	option struct {
		value        reflect.Value
		field        reflect.StructField
		short        string
		long         string
		required     bool
		desc         string
		defaultValue string
		format       string
		handler      string
		specified    bool
	}
)

func toSnakeCase(s string) string {
	name := []byte{}
	for _, c := range []byte(s) {
		if c > 64 && c < 91 {
			name = append(name, byte(95), c+32)
		} else {
			name = append(name, c)
		}
	}
	name = regexp.MustCompile(`_{2,}`).ReplaceAll(name, []byte{95})
	return string(bytes.Trim(name, "_"))
}

func isAlphaNumeric(b byte) bool {
	if b >= 48 && b <= 57 { // 0 - 9
		return true
	} else if b >= 65 && b <= 90 { // A - Z
		return true
	} else if b >= 97 && b <= 122 { // a - z
		return true
	} else {
		return false
	}
}

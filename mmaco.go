package mmaco

import "reflect"

const (
	tagName = "mmaco"

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

	SubCommandInterface interface {
		Init() error
		Validate([]string) error
		Run([]string) error
	}

	option struct {
		field        reflect.StructField
		short        string
		long         string
		required     bool
		desc         string
		defaultValue string
		format       string
		handler      string
	}
)

package mmaco

const (
	tagName           = "mmaco"
	helpCommandName   = "help"
	nameFieldLabel    = "Name"
	initFieldLabel    = "Init"
	isValidFieldLabel = "IsValid"
	runFieldLabel     = "Run"
)

type (
	SubCommandInterface interface {
		Name() string
		Init() error
		Validate() error
		Run() error
	}
)

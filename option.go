package mmaco

type (
	option struct {
		name string
	}
)

func newOpts(name string, metas []meta) option {

	return option{name: name}
}

func (o option) Name() string {
	return o.name
}

func (o option) Parse(args []string) error {
	return nil
}

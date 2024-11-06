package mmaco

type (
	opts struct {
	}
)

func newOpts(metas []meta) opts {
	return opts{}
}

func (o opts) parse(args []string) error {
	return nil
}

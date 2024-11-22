package mmaco

import "testing"

type (
	subCmdTest struct {
	}
)

func (sc subCmdTest) Init() error {
	return nil
}
func (sc subCmdTest) Validate(args []string) error {
	return nil
}
func (sc subCmdTest) Run(args []string) error {
	return nil
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	cleanup()
}

func setup() {

}

func cleanup() {

}

package mmaco

import "testing"

func TestMain(m *testing.M) {
	println("前処理")
	m.Run()
	println("前処理")
}

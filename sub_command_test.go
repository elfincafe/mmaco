package mmaco

import (
	"reflect"
	"testing"
	"time"
)

func isSameSubCmd(aSt, bSt *SubCommand) bool {
	a := reflect.TypeOf(aSt)
	b := reflect.TypeOf(bSt)
	if a.PkgPath() == b.PkgPath() && a.Name() == b.Name() {
		return true
	} else {
		return false
	}
}

func TestNewSubCommand(t *testing.T) {
	// Test Case
	cases := []struct {
		cmd SubCommandInterface
		st  *SubCommand
	}{
		{cmd: subCmd0{}, st: new(SubCommand)},
		{cmd: subCmd1{}, st: new(SubCommand)},
		{cmd: subCmd2{}, st: new(SubCommand)},
	}

	// Test
	loc := time.FixedZone("JST", 9*60*60)
	for i, c := range cases {
		s := newSubCommand(c.cmd, loc)
		if !isSameSubCmd(s, c.st) {
			t.Errorf(`[%d] Expected: *subCommand, Result: %v`, i, s)
		}
	}
}

func TestSubCommandParse(t *testing.T) {
	// Test Case
	cases := []struct {
		sc       SubCommandInterface
		expected int
	}{
		{sc: subCmd0{}, expected: 9},
		{sc: subCmd1{}, expected: 16},
		{sc: subCmd2{}, expected: 0},
	}
	// Test
	var err error
	loc := time.FixedZone("JST", 9*60*60)
	for i, c := range cases {
		v := newSubCommand(c.sc, loc)
		err = v.parse()
		if err != nil {
			continue
		}
		if c.expected != len(v.opts) {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, len(v.opts))
		}
	}
}

func TestSubCommandName(t *testing.T) {
	// Test Case
	cases := []struct {
		sc       SubCommandInterface
		expected string
	}{
		{sc: subCmd0{}, expected: "sub_cmd0"},
		{sc: subCmd1{}, expected: "sub_cmd1"},
		{sc: subCmd2{}, expected: "sub_cmd2"},
	}
	// Test
	loc := time.FixedZone("JST", 9*60*60)
	for i, c := range cases {
		s := newSubCommand(c.sc, loc)
		if s.Name != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, s.Name)
		}
	}
}

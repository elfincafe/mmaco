package mmaco

import (
	"testing"
	"time"
)

type (
	subCmd0 struct {
		Desc    string
		field1  string `mmaco:""`
		field2  string `mmaco:"short=2"`
		field3  string `mmaco:"long=field3"`
		field4  string `mmaco:"desc=desc4"`
		field5  string `mmaco:"required"`
		field6  string `mmaco:"default=default value"`
		field7  string `mmaco:"format=2006/01/02 15:04:05"`
		field8  string `mmaco:"handler=Handler"`
		field9  string `mmaco:"short=9,long=field9,desc= desc9, test , default= default9, Value , ,required,format=Mon, 02 Jan 2006 15:04:05 MST,handler=Handler"`
		field10 string `mmaco:"short=a,long=field10,desc= desc10, test , default= default10, Value , ,handler=Handler"`
	}
	subCmd1 struct {
		string_0  string    `mmaco:"short=s,long=string,desc=string desc"`
		bool_0    bool      `mmaco:"short=b,long=bool,desc=bool desc"`
		int_0     int       `mmaco:"short=i,long=int,desc=int desc"`
		int8_0    int8      `mmaco:"short=I,long=int8,desc=int8 desc"`
		int16_0   int16     `mmaco:"short=J,long=int16,desc=int16 desc"`
		int32_0   int32     `mmaco:"short=K,long=int32,desc=int32 desc"`
		int64_0   int64     `mmaco:"short=L,long=int64,desc=int64 desc"`
		uint_0    uint      `mmaco:"short=M,long=uint,desc=uint desc"`
		uint8_0   uint8     `mmaco:"short=U,long=uint8,desc=uint8 desc"`
		uint16_0  uint16    `mmaco:"short=V,long=uint16,desc=uint16 desc"`
		uint32_0  uint32    `mmaco:"short=W,long=uint32,desc=uint32 desc"`
		uint64_0  uint64    `mmaco:"short=X,long=uint64,desc=uint64 desc"`
		float32_0 float32   `mmaco:"short=f,long=float32,desc=float32 desc"`
		float64_0 float64   `mmaco:"short=F,long=float64,desc=float64 desc"`
		time_0    time.Time `mmaco:"short=t,long=time,desc=time desc"`
		unknown_0 complex64 `mmaco:"short=c,long=complex64,desc=complex64 desc"`
	}
	subCmd2 struct {
	}
)

func (cmd subCmd0) Run(...[]string) error {
	return nil
}

func (cmd subCmd1) Run(...[]string) error {
	return nil
}

func (cmd subCmd2) Run(...[]string) error {
	return nil
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	cleanup()
}

func setup() {
	println("* Setting up")
}

func cleanup() {
	println("* Cleaning up")
}

func TestToSnakeCase(t *testing.T) {
	// Test Case
	cases := []struct {
		s        string
		expected string
	}{
		{s: "TestSubCommand", expected: "test_sub_command"},
		{s: "TestSubCommand2", expected: "test_sub_command2"},
		{s: "Test_Sub_Command_3", expected: "test_sub_command_3"},
		{s: "4Test_Sub_Command_", expected: "4_test_sub_command"},
	}
	// Test
	for i, c := range cases {
		name := toSnakeCase(c.s)
		if name != c.expected {
			t.Errorf("[%d] Expected: %v, Result: %v", i, c.expected, name)
		}
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	// Test Case
	cases := []struct {
		s        string
		expected bool
	}{
		{s: "/", expected: false},
		{s: "0", expected: true},
		{s: "9", expected: true},
		{s: ":", expected: false},
		{s: "@", expected: false},
		{s: "A", expected: true},
		{s: "Z", expected: true},
		{s: "[", expected: false},
		{s: "`", expected: false},
		{s: "a", expected: true},
		{s: "z", expected: true},
		{s: "{", expected: false},
	}
	// Test
	for i, c := range cases {
		b := []byte(c.s)
		result := isAlphaNumeric(b[0])
		if result != c.expected {
			t.Errorf("[%d] %s Expected: %v, Result: %v", i, c.s, c.expected, result)
		}
	}
}

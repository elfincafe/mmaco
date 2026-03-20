# mmaCo

mmaCo is a framework for command-line interfaces. If you're building systems in Go, we believe it's a very useful option. By using mmaCo, you can easily create a subcommand-based CLI.

mmaCo has the following features:
* Simple: With limited functionality, it requires minimal learning effort.
* Fast: Processing is fast and does not slow down execution.
* Help: Automatically generates help sub command.

We hope you have a wonderful development experience with mmaCo.

# Quick Start
You just add structures that implement the SubCommand interface to the mmaCo structure. You can automatically retrieve subcommand options by declaring them as tagged fields in the structure.
In the following example, we are creating the `build` command as a subcommand. We have configured it so that the optimization level can be specified as an option for the subcommand.
```go
package main

import (
	"fmt"

	"github.com/elfincafe/mmaco"
)

type (
	Build struct {
		Optimize int `mmaco:"short:o,long=optimize,desc=optimize level,required"`
	}
)

func (b *Build) Init() {
	println("Initialize")
}

func (b *Build) Run(ctx *mmaco.Context) error {
	fmt.Printf("Optimize level: %d", b.Optimize)
	return nil
}

func main() {
	cmd := mmaco.New("sample")
	cmd.Add(&Build{}, "build", "build the binary from the source files.")
	cmd.Run()
}
```
When you run the code above, the following output is displayed.
```bash
$ go run main.go build -o 2
Initialize
Optimize level: 2
```
# Usage
## Sub Command
The requirements for a structure to be registered as a subcommand are simple: it must implement the mmaco.SubCommand interface. That’s all there is to it.
The mmaco.SubCommand interface requires the implementation of the following two functions:
* Init()
* Run(*mmaco.Context) error
## Options
Subcommands allow you to declaratively define options by using tags on fields of structures that implement the SubCommand interface.
The types of fields that can be used as options are as follows:
* string
* bool
* int
* int8
* int16
* int32
* int64
* uint
* uint8
* uint16
* uint32
* uint64
* float32
* float64
* time.Time

The following items can be configured using tags.
|Item|Description|
|-|-|
|short|Short option name. 1 character from numerics or alphabets(0-9 a-z A-Z). Either `short` or `long` is required.|
|long|Long option name. 2 to 15 characters from numerics, alphabets, hyphen or undersocre(0-9 a-z A-Z - _). Either `short` or `long` is required.|
|required|required option.|
|desc|Description. Displayed within the help command.|
|format|format for time.Time. [time package format](https://pkg.go.dev/time#pkg-constants) or custom layouts (e.g. 2006/01/02)|
```go
Build struct {
  Label string `mmaco:"short=l,long=label,desc=display name,required"`
  Today time.Time `mmaco:"short=t,desc=today,format=RFC3339"`
}
```
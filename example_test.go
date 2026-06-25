package res_test

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/res"
)

type loggerWire struct{}

func (loggerWire) NewResource() (any, error) {
	return &struct{ Level string }{Level: "DEBUG"}, nil
}

type Logger struct {
	Level string
}

func ExampleNew() {
	reg := res.New()
	_ = reg.Add(loggerWire{})

	entries := reg.GetByType(reflect.TypeFor[loggerWire]())
	if len(entries) > 0 {
		fmt.Println("wired")
	}
	// Output: wired
}

func ExampleMustAddToGlobalWithTags() {
	res.MustAddToGlobalWithTags(loggerWire{}, res.TagReplaceable)
	fmt.Println("ok")
	// Output: ok
}

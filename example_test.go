package res_test

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/res"
)

type Logger struct {
	Level string
}

func ExampleNew() {
	reg := res.New()
	_ = reg.Add(&Logger{Level: "DEBUG"})

	entries := reg.GetByType(reflect.TypeFor[*Logger]())
	if len(entries) > 0 {
		fmt.Println(entries[0].Value.(*Logger).Level)
	}
	// Output: DEBUG
}

func ExampleNew_addWithTags() {
	reg := res.New()
	_ = reg.AddWithTags("default-key", res.TagReplaceable)
	_ = reg.Add("my-secret-key")

	entries := reg.GetByType(reflect.TypeFor[string]())
	fmt.Println("count:", len(entries))
	// Output: count: 2
}

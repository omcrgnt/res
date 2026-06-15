package res_test

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/res"
)

type Logger struct {
	Level string
}

func ExampleAdd() {
	res.Add(&Logger{Level: "DEBUG"})

	entries := res.GetByType(reflect.TypeFor[*Logger]())
	if len(entries) > 0 {
		fmt.Println(entries[0].Value.(*Logger).Level)
	}
	// Output: DEBUG
}

func ExampleAddWithTags() {
	res.AddWithTags("default-key", res.TagReplaceable)
	res.Add("my-secret-key")

	entries := res.GetByType(reflect.TypeFor[string]())
	fmt.Println("count:", len(entries))
	// Output: count: 2
}

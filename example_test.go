package res_test

import (
	"fmt"

	"github.com/omcrgnt/res"
)

type Logger struct {
	Level string
}

func ExampleAddAll() {
	res.AddAll(&Logger{Level: "DEBUG"})

	logger, ok := res.Get[*Logger]()
	if ok {
		fmt.Println(logger.Level)
	}
	// Output: DEBUG
}

func ExampleGet() {
	res.Add("my-secret-key")

	val, ok := res.Get[string]()
	if ok {
		fmt.Println("Found:", val)
	}
	// Output: Found: my-secret-key
}

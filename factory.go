package res

import (
	"fmt"
	"maps"
	"reflect"
	"sync"

	"github.com/mcrgnt/extractor"
)

var gf = newFactory()

func Register(b Builder) {
	gf.register(b)
}

func Build(source any) error {
	return gf.withSource(source).run()
}

type factory struct {
	mu                sync.RWMutex
	systemBuilderList map[reflect.Type]Builder
	userSource        any
}

func newFactory() *factory {
	return &factory{
		systemBuilderList: make(map[reflect.Type]Builder),
	}
}

func (t *factory) withSource(source any) *factory {
	t.userSource = source
	return t
}

func (t *factory) register(b Builder) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.systemBuilderList[reflect.TypeOf(b)] = b
}

func (t *factory) run() error {
	var (
		finalBuilders   = make(map[reflect.Type]Builder)
		userBuilderList = extractor.New[Builder](t.userSource).Extract()
	)

	t.mu.Lock()
	maps.Copy(finalBuilders, t.systemBuilderList)
	t.mu.Unlock()

	for _, b := range userBuilderList {
		bType := reflect.TypeOf(b)
		finalBuilders[bType] = b
	}

	var reg = newRegistry()
	for _, b := range finalBuilders {
		resource, err := b.Build()
		if err != nil {
			return fmt.Errorf("build resource failed for builder %T: %w", b, err)
		}

		if err := addAny(reg, resource); err != nil {
			return err
		}
	}

	t.mu.Lock()
	globalRegistry = reg
	t.mu.Unlock()

	return nil
}

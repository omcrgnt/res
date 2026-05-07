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

func GlobalFactory() *factory {
	return gf
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

func (t *factory) WithSource(source any) *factory {
	t.userSource = source
	return t
}

func (t *factory) register(b Builder) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.systemBuilderList[reflect.TypeOf(b)] = b
}

func (t *factory) Build() error {
	reg, err := t.build()
	if err != nil {
		return err
	}

	// Обновляем глобальную переменную
	t.mu.Lock()
	globalRegistry = reg
	t.mu.Unlock()

	return nil
}

func (t *factory) build() (*registry, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var (
		reg             = newRegistry()
		userBuilderList = extractor.New[Builder](t.userSource).Extract()
		finalBuilders   = make(map[reflect.Type]Builder)
	)

	maps.Copy(finalBuilders, t.systemBuilderList)

	for _, b := range userBuilderList {
		bType := reflect.TypeOf(b)
		finalBuilders[bType] = b
	}

	for _, b := range finalBuilders {
		resource, err := b.Build()
		if err != nil {
			return nil, fmt.Errorf("build resource failed for builder %T: %w", b, err)
		}

		if err := addAny(reg, resource); err != nil {
			return nil, err
		}
	}

	return reg, nil
}

package res

import (
	"errors"
	"fmt"
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
	systemBuilderList map[string]Builder

	userSource      any
	userBuilderList map[string]Builder

	err error
}

func newFactory() *factory {
	return &factory{
		systemBuilderList: make(map[string]Builder),
		userBuilderList:   make(map[string]Builder),
	}
}

func (t *factory) WithSource(source any) *factory {
	t.userSource = source
	return t
}

func (t *factory) Build() (*registry, error) {
	return t.build()
}

func (t *factory) register(b Builder) {
	t.mu.Lock()
	defer t.mu.Unlock()

	label := b.Label()
	if _, ok := t.systemBuilderList[label]; ok {
		t.err = errors.Join(t.err, fmt.Errorf("resource already registered with label: %q", label))
	}
	t.systemBuilderList[label] = b
}

func (t *factory) build() (*registry, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.err != nil {
		return nil, t.err
	}

	var (
		builderList     []Builder
		registry        = newRegistry()
		userBuilderList = extractor.New[Builder](t.userSource).Extract()
	)

	for _, builder := range userBuilderList {
		switch t.unsafeCheckLabelExistsInSystemBuilderList(builder) {
		case true:
			if t.unsafeCompareBuilderWithExistingSystemBuilder(builder) {
				t.systemBuilderList[builder.Label()] = builder
			}
		default:
			t.userBuilderList[builder.Label()] = builder
		}
	}

	for _, builder := range t.systemBuilderList {
		builderList = append(builderList, builder)
	}

	for _, builder := range t.userBuilderList {
		builderList = append(builderList, builder)
	}

	for _, builder := range builderList {
		resource, err := builder.Build()
		if err != nil {
			return nil, err
		}
		addAny(registry, resource)
	}

	return registry, nil
}

func (t *factory) unsafeCheckLabelExistsInSystemBuilderList(builder Builder) bool {
	_, ok := t.systemBuilderList[builder.Label()]
	return ok
}

func (t *factory) unsafeCompareBuilderWithExistingSystemBuilder(builder Builder) bool {
	return reflect.TypeOf(t.systemBuilderList[builder.Label()]) == reflect.TypeOf(builder)
}

package res

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/mcrgnt/extractor"
)

var ErrResourceRegisteredAlready = errors.New("resource already registered")

var (
	gf = newFactory()
)

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

func (t *factory) register(b Builder) {
	t.mu.Lock()
	defer t.mu.Unlock()

	label := b.Label()
	if _, ok := t.systemBuilderList[label]; ok {
		t.err = errors.Join(t.err, fmt.Errorf("%w: label %q", ErrResourceRegisteredAlready, label))
	}
	t.systemBuilderList[label] = b
}

func (t *factory) WithSource(source any) *factory {
	t.userSource = source
	return t
}

func (t *factory) Build() ([]any, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.err != nil {
		return nil, t.err
	}

	var (
		resourceList    []any
		userBuilderList = extractor.New[Builder](t.userSource).Extract()
	)

	for _, builder := range userBuilderList {
		t.userBuilderList[builder.Label()] = builder
	}

	allBuilders := make(map[string]Builder, len(t.systemBuilderList)+len(userBuilderList))
	maps.Copy(allBuilders, t.systemBuilderList)

	for _, b := range userBuilderList {
		allBuilders[b.Label()] = b
	}

	for _, builder := range allBuilders {
		resource, err := builder.Build()
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, resource)
	}

	return resourceList, nil
}

// func (t *factory) CloneSystem() *factory {
// 	t.mu.RLock()
// 	defer t.mu.RUnlock()

// 	newFactory := newFactory()
// 	maps.Copy(newFactory.systemBuilderList, t.systemBuilderList)
// 	return newFactory
// }

package res

import (
	"fmt"
	"maps"
	"reflect"
	"sync"

	"github.com/mcrgnt/extractor"
)

var gf = newFactory()

// Register регистрирует системный билдер в глобальной фабрике.
//
// Обычно вызывается в блоке init() пакетов инфраструктуры или сторонних библиотек.
// Если при вызове Build() в источнике (source) обнаружится билдер того же типа,
// системный билдер будет заменен пользовательским.
func Register(b Builder) {
	gf.register(b)
}

// Build собирает глобальный реестр ресурсов, используя предоставленный источник (source).
//
// В качестве источника может выступать структура (struct), срез (slice) или карта (map).
// Функция извлекает из источника все объекты, реализующие интерфейс Builder,
// инициализирует их и сохраняет результат в реестре.
//
// Если типы билдеров в источнике совпадают с системными (зарегистрированными через Register),
// пользовательские билдеры переопределяют системные.
func Build(source any) error {
	var localFactory = newFactory()

	gf.mu.RLock()
	maps.Copy(localFactory.systemBuilderList, gf.systemBuilderList)
	gf.mu.RUnlock()

	reg, err := localFactory.withSource(source).run()
	if err != nil {
		return err
	}

	gf.mu.Lock()
	globalRegistry = reg
	gf.mu.Unlock()

	return nil
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

func (t *factory) run() (*registry, error) {
	var finalBuilders = make(map[reflect.Type]Builder)
	var userBuilderList []Builder

	if t.userSource != nil {
		userBuilderList = extractor.New[Builder](t.userSource).Extract()
	}

	maps.Copy(finalBuilders, t.systemBuilderList)

	for _, b := range userBuilderList {
		bType := reflect.TypeOf(b)
		finalBuilders[bType] = b
	}

	var reg = newRegistry()
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

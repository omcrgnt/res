package unique

import (
	"reflect"

	"github.com/omcrgnt/res"
)

type uniqueEntry struct {
	reg   *Registry
	inner res.Entry
}

func (e *uniqueEntry) Type() reflect.Type {
	return e.inner.Type()
}

func (e *uniqueEntry) Value() any {
	return e.inner.Value()
}

func (e *uniqueEntry) Has(tag res.Tag) bool {
	return e.inner.Has(tag)
}

func (e *uniqueEntry) Regular() bool {
	return e.inner.Regular()
}

func (e *uniqueEntry) Replaceable() bool {
	return e.inner.Replaceable()
}

func (e *uniqueEntry) Fixed() bool {
	return e.inner.Fixed()
}

func (e *uniqueEntry) Tags() []res.Tag {
	return e.inner.Tags()
}

func (e *uniqueEntry) GetCustomTag(key string) (any, bool) {
	return e.inner.GetCustomTag(key)
}

func (e *uniqueEntry) ChangeValue(new any) error {
	if new == nil {
		return errNilValue
	}

	newTyp := reflect.TypeOf(new)
	oldTyp := e.inner.Type()
	if newTyp != oldTyp {
		entries := e.reg.reg.GetByType(newTyp)
		for _, other := range entries {
			if other.Value() != e.inner.Value() {
				return ErrRegularExists
			}
		}
	}

	return e.inner.ChangeValue(new)
}

func wrapEntry(reg *Registry, e res.Entry) res.Entry {
	return &uniqueEntry{reg: reg, inner: e}
}

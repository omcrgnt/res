package unique

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/res"
)

var errNilValue = fmt.Errorf("unique: nil value")

func (r *Registry) WalkEntries(fn func(res.Entry) bool) {
	r.reg.WalkEntries(func(e res.Entry) bool {
		return fn(wrapEntry(r, e))
	})
}

func (r *Registry) GetByType(t reflect.Type) []res.Entry {
	entries := r.reg.GetByType(t)
	out := make([]res.Entry, len(entries))
	for i, e := range entries {
		out[i] = wrapEntry(r, e)
	}
	return out
}

func (r *Registry) GetByInterface(iface reflect.Type) []res.Entry {
	entries := r.reg.GetByInterface(iface)
	out := make([]res.Entry, len(entries))
	for i, e := range entries {
		out[i] = wrapEntry(r, e)
	}
	return out
}

func (r *Registry) GetOneByType(t reflect.Type) (any, error) {
	return r.reg.GetOneByType(t)
}

func (r *Registry) GetOneByInterface(iface reflect.Type) (any, error) {
	return r.reg.GetOneByInterface(iface)
}

func (r *Registry) Transform(fns ...res.TransformFunc) error {
	return r.reg.Transform(fns...)
}

func (r *Registry) Remove(v any) error {
	return r.reg.Remove(v)
}

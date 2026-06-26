package unique

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/res"
)

var errNilValue = fmt.Errorf("unique: nil value")

func (r *Registry) WalkEntries(fn func(res.Entry) bool) {
	r.reg.WalkEntries(fn)
}

func (r *Registry) GetByType(t reflect.Type) []res.Entry {
	return r.reg.GetByType(t)
}

func (r *Registry) GetByInterface(iface reflect.Type) []res.Entry {
	return r.reg.GetByInterface(iface)
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

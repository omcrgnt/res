package respolicy

import (
	"reflect"

	"github.com/omcrgnt/res"
)

// Registry wraps [res.Registry] and applies [AddPolicy] on Add and AddWithTags.
type Registry struct {
	reg      res.Registry
	policies []AddPolicy
}

// Wrap returns a policy-aware registry over reg.
// When policies is empty, Add behaves like reg.Add (no transformation).
func Wrap(reg res.Registry, policies ...AddPolicy) *Registry {
	return &Registry{reg: reg, policies: policies}
}

// Registry returns the underlying [res.Registry].
func (w *Registry) Underlying() res.Registry {
	return w.reg
}

var _ res.Registry = (*Registry)(nil)

func (w *Registry) Add(v any) error {
	stored, err := applyPolicies(v, w.policies)
	if err != nil {
		return err
	}
	return w.reg.Add(stored)
}

func (w *Registry) AddWithTags(v any, tags ...res.Tag) error {
	stored, err := applyPolicies(v, w.policies)
	if err != nil {
		return err
	}
	return w.reg.AddWithTags(stored, tags...)
}

func (w *Registry) WalkEntries(fn func(res.Entry) bool) {
	w.reg.WalkEntries(fn)
}

func (w *Registry) GetByType(t reflect.Type) []res.Entry {
	return w.reg.GetByType(t)
}

func (w *Registry) GetByInterface(iface reflect.Type) []res.Entry {
	return w.reg.GetByInterface(iface)
}

func (w *Registry) GetOneByType(t reflect.Type) (any, error) {
	return w.reg.GetOneByType(t)
}

func (w *Registry) GetOneByInterface(iface reflect.Type) (any, error) {
	return w.reg.GetOneByInterface(iface)
}

func (w *Registry) Transform(fns ...res.TransformFunc) error {
	return w.reg.Transform(fns...)
}

func (w *Registry) Remove(v any) error {
	return w.reg.Remove(v)
}

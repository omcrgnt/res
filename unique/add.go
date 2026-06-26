package unique

import (
	"reflect"

	"github.com/omcrgnt/res"
)

// Add registers v with [res.TagRegular], replacing an existing [res.TagReplaceable] entry.
func (r *Registry) Add(v any) error {
	if v == nil {
		return errNilValue
	}

	typ := reflect.TypeOf(v)
	entries := r.reg.GetByType(typ)
	switch len(entries) {
	case 0:
		return r.reg.AddWithTags(v, res.TagRegular)
	case 1:
		e := entries[0]
		switch {
		case e.Replaceable():
			return res.ReplaceAtType(r.reg, v, res.TagRegular)
		case isRegular(e):
			return ErrRegularExists
		case e.Fixed():
			return ErrFixed
		default:
			return ErrRegularExists
		}
	default:
		return ErrDuplicateType
	}
}

// AddReplaceable registers v with [res.TagReplaceable].
func (r *Registry) AddReplaceable(v any) error {
	if v == nil {
		return errNilValue
	}

	typ := reflect.TypeOf(v)
	entries := r.reg.GetByType(typ)
	switch len(entries) {
	case 0:
		return r.reg.AddWithTags(v, res.TagReplaceable)
	case 1:
		e := entries[0]
		switch {
		case e.Replaceable():
			return res.ReplaceAtType(r.reg, v, res.TagReplaceable)
		case isRegular(e), e.Fixed():
			return ErrTypeOccupied
		default:
			return ErrTypeOccupied
		}
	default:
		return ErrDuplicateType
	}
}

// AddFixed registers v with [res.TagFixed] when no entry exists for the type.
func (r *Registry) AddFixed(v any) error {
	if v == nil {
		return errNilValue
	}

	typ := reflect.TypeOf(v)
	entries := r.reg.GetByType(typ)
	switch len(entries) {
	case 0:
		return r.reg.AddWithTags(v, res.TagFixed)
	case 1:
		return ErrTypeOccupied
	default:
		return ErrDuplicateType
	}
}

func isRegular(e res.Entry) bool {
	if e.Regular() {
		return true
	}
	return !e.Replaceable() && !e.Fixed()
}

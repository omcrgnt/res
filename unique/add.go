package unique

import (
	"fmt"
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

// AddWithCustomTag registers v with [res.TagRegular] and one custom tag key/value.
func (r *Registry) AddWithCustomTag(v any, key string, val any) error {
	if v == nil {
		return errNilValue
	}
	if key == "" {
		return fmt.Errorf("unique: empty custom tag key")
	}
	if val == nil {
		return fmt.Errorf("unique: nil custom tag value")
	}

	typ := reflect.TypeOf(v)
	customTags := map[string]any{key: val}
	entries := r.reg.GetByType(typ)
	switch len(entries) {
	case 0:
		return res.AddWithTagsAndCustomTags(r.reg, v, customTags, res.TagRegular)
	case 1:
		e := entries[0]
		switch {
		case e.Replaceable():
			return res.ReplaceAtTypeWithCustomTags(r.reg, v, customTags, res.TagRegular)
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

func (r *Registry) addReplaceable(v any) error {
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

func (r *Registry) addFixed(v any) error {
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

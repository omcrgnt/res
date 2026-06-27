package unique

import (
	"errors"
	"reflect"

	"github.com/omcrgnt/res"
)

var (
	errNilRegistry  = errors.New("unique: nil registry")
	errSameRegistry = errors.New("unique: cannot merge registry into itself")
)

// Merge copies resources from src into dst by calling [Registry.Add] or [Registry.AddWithCustomTag].
func Merge(dst, src *Registry) error {
	if dst == nil || src == nil {
		return errNilRegistry
	}
	if dst == src {
		return errSameRegistry
	}

	var err error
	src.reg.WalkEntries(func(e res.Entry) bool {
		customTags := res.EntryCustomTags(e)
		switch len(customTags) {
		case 0:
			err = dst.Add(e.Value())
		case 1:
			for key, val := range customTags {
				err = dst.AddWithCustomTag(e.Value(), key, val)
				break
			}
		default:
			err = mergeEntryWithCustomTags(dst, e.Value(), customTags)
		}
		return err == nil
	})
	return err
}

func mergeEntryWithCustomTags(dst *Registry, v any, customTags map[string]any) error {
	if v == nil {
		return errNilValue
	}

	typ := reflect.TypeOf(v)
	entries := dst.reg.GetByType(typ)
	switch len(entries) {
	case 0:
		return res.AddWithTagsAndCustomTags(dst.reg, v, customTags, res.TagRegular)
	case 1:
		e := entries[0]
		switch {
		case e.Replaceable():
			return res.ReplaceAtTypeWithCustomTags(dst.reg, v, customTags, res.TagRegular)
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

package unique

import (
	"errors"

	"github.com/omcrgnt/res"
)

var (
	errNilRegistry  = errors.New("unique: nil registry")
	errSameRegistry = errors.New("unique: cannot merge registry into itself")
)

// Merge copies resources from src into dst by calling [Registry.Add] for each entry in src.
func Merge(dst, src *Registry) error {
	if dst == nil || src == nil {
		return errNilRegistry
	}
	if dst == src {
		return errSameRegistry
	}

	var err error
	src.WalkEntries(func(e res.Entry) bool {
		if addErr := dst.Add(e.Value); addErr != nil {
			err = addErr
			return false
		}
		return true
	})
	return err
}

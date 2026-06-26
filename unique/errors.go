package unique

import "errors"

var (
	ErrFixed         = errors.New("unique: cannot override entry with TagFixed")
	ErrRegularExists = errors.New("unique: entry for type already has TagRegular")
	ErrTypeOccupied  = errors.New("unique: entry for type has incompatible tag")
	ErrDuplicateType = errors.New("unique: multiple entries for same concrete type")
)

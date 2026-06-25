package res

import "fmt"

func validateRegistrant(v any) error {
	_, isNew := v.(NewResourceer)
	_, isBuild := v.(BuildConfiger)
	switch {
	case isNew && isBuild:
		return fmt.Errorf("%w", ErrInvalidRegistrant)
	case isNew || isBuild:
		return nil
	default:
		return fmt.Errorf("%w: %T", ErrInvalidRegistrant, v)
	}
}

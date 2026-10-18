package complaint

import "errors"

// Sentinel errors. Callers check them with errors.Is, because the functions in
// this package wrap them with extra context using fmt.Errorf("%w: ...").
var (
	// ErrInvalidInput is returned when the data of a complaint breaks a validation rule.
	ErrInvalidInput = errors.New("invalid complaint input")

	// ErrNotFound is returned when no complaint exists for the requested ID.
	ErrNotFound = errors.New("complaint not found")

	// ErrInvalidTransition is returned when an operation is not allowed in the
	// current status of the complaint (for example, processing it twice).
	ErrInvalidTransition = errors.New("invalid status transition")
)

package assign

import (
	"github.com/pkg/errors"
)

// ErrNotSupported is returned when a type conversion is not supported.
// This error indicates that the source type cannot be converted to the target type
// using the available conversion functions.
var ErrNotSupported = errors.New("not supported")

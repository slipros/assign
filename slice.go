package assign

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// SliceOptionFunc is a functional option for configuring SliceString behavior
type SliceOptionFunc func(*sliceOptions)

// sliceOptions contains configuration for SliceString
type sliceOptions struct {
	separator string
}

// defaultSliceOptions returns the default slice options
func defaultSliceOptions() sliceOptions {
	return sliceOptions{
		separator: ",",
	}
}

// WithSeparator sets a custom separator for joining string slices.
// Used as an option in SliceString.
//
// Example: WithSeparator("|") // Join with pipe character
func WithSeparator(sep string) SliceOptionFunc {
	return func(o *sliceOptions) {
		o.separator = sep
	}
}

// SliceString converts a string slice to appropriate types for target fields.
// Handles conversions to strings (joins elements), numeric slices, boolean slices,
// and other compatible types.
//
// Parameters:
//   - field: Target field to set (reflect.Value).
//   - arr: String slice to convert and set.
//   - options: Optional settings like custom separator (default: ",").
//
// Returns:
//   - error: If conversion or assignment fails.
func SliceString(to reflect.Value, from []string, options ...SliceOptionFunc) error {
	// Check if the to is settable
	if !to.CanSet() {
		return errors.Errorf("field of type %s is not settable", to.Type())
	}

	// Default options
	opts := defaultSliceOptions()
	for _, opt := range options {
		opt(&opts)
	}

	// Handle nil pointer initialization
	if to.Kind() == reflect.Pointer {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}

		return SliceString(to.Elem(), from, options...)
	}

	fieldType := to.Type()
	switch to.Kind() {
	case reflect.String:
		// Join the string slice into a single string with specified separator
		to.SetString(strings.Join(from, opts.separator))
		return nil

	case reflect.Slice:
		elemType := fieldType.Elem()

		switch elemType.Kind() {
		case reflect.String:
			// Create a new slice of the correct type and convert each element
			slice := reflect.MakeSlice(fieldType, len(from), len(from))
			for i, v := range from {
				slice.Index(i).Set(reflect.ValueOf(v).Convert(elemType))
			}

			to.Set(slice)

			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool:

			slice := reflect.MakeSlice(fieldType, len(from), len(from))
			for i, strVal := range from {
				elemValue := reflect.New(elemType).Elem()
				if err := String(elemValue, strVal); err != nil {
					return errors.Wrapf(err, "failed to convert string '%s' to %s", strVal, elemType.String())
				}
				slice.Index(i).Set(elemValue)
			}

			to.Set(slice)

			return nil
		case reflect.Interface:
			slice := reflect.MakeSlice(fieldType, len(from), len(from))
			for i, v := range from {
				slice.Index(i).Set(reflect.ValueOf(v))
			}

			to.Set(slice)

			return nil
		}

	case reflect.Interface:
		// For any fields, prefer setting as []string directly
		to.Set(reflect.ValueOf(from))

		return nil
	}

	// If the to doesn't match any of the above types, return an error
	return errors.Wrapf(ErrNotSupported, "cannot convert []string to field of type %s", fieldType.String())
}

package assign

import (
	"math"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// FloatValue is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type FloatValue interface {
	~float32 | ~float64
}

// Float converts a floating-point value to the appropriate type for a target field.
// Handles conversion to strings, booleans, numeric types, complex types, and interfaces.
// Special values (NaN, Infinity) are handled appropriately for each target type.
//
// Parameters:
//   - field: Target field to set (reflect.Value).
//   - number: Floating-point value to convert and set.
//
// Returns:
//   - error: If conversion or assignment fails.
func Float[F FloatValue](to reflect.Value, from F) error {
	// Check if the to is settable
	if !to.CanSet() {
		return errors.Errorf("field of type %s is not settable", to.Type())
	}

	// Check for special float values (NaN and Infinity)
	floatVal := float64(from)
	isSpecial := math.IsNaN(floatVal) || math.IsInf(floatVal, 0)

	switch to.Kind() {
	case reflect.String:
		// Format float based on value:
		// - Use normal decimal notation for most numbers
		// - Use scientific notation only for very large or small numbers
		var str string
		// Handle special cases explicitly
		switch {
		case math.IsNaN(floatVal):
			str = "NaN"
		case math.IsInf(floatVal, 1):
			str = "+Inf"
		case math.IsInf(floatVal, -1):
			str = "-Inf"
		default:
			abs := math.Abs(floatVal)
			if (abs >= 0.0001 && abs < 10000000) || floatVal == 0 {
				// Use regular decimal format for normal range numbers
				str = strconv.FormatFloat(floatVal, 'f', -1, 64)
			} else {
				// Use scientific notation for very large or small numbers
				str = strconv.FormatFloat(floatVal, 'E', -1, 64)
			}
		}
		to.SetString(str)

		return nil
	case reflect.Bool:
		// Convert float to bool (true if > 0, false otherwise)
		to.SetBool(from > 0 || from < 0)

		return nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		// Handle special cases for int types
		if isSpecial {
			if math.IsNaN(floatVal) {
				return errors.Errorf("cannot convert NaN to integer type %s", to.Type())
			}
			if math.IsInf(floatVal, 1) {
				return errors.Errorf("cannot convert +Infinity to integer type %s", to.Type())
			}
			if math.IsInf(floatVal, -1) {
				return errors.Errorf("cannot convert -Infinity to integer type %s", to.Type())
			}
		}

		// Check if value is too small and will be truncated to zero
		if math.Abs(floatVal) < 1.0 && floatVal != 0.0 { //nolint:staticcheck // for future use
			// This is not an error, but a warning that could be logged in a real-world app
			// For now, just proceed with truncation
		}

		// Check if the value is within int64 range
		if floatVal > float64(math.MaxInt64) || floatVal < float64(math.MinInt64) {
			return errors.Errorf("value %v is outside the range of any signed integer type", floatVal)
		}

		// Check for specific int sizes
		int64Val := int64(floatVal)
		if err := checkSignedIntegerRange(to, int64Val); err != nil {
			return err
		}

		to.SetInt(int64Val)
		return nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		// Handle special cases for uint types
		if isSpecial {
			if math.IsNaN(floatVal) {
				return errors.Errorf("cannot convert NaN to unsigned integer type %s", to.Type())
			}
			if math.IsInf(floatVal, 1) {
				return errors.Errorf("cannot convert +Infinity to unsigned integer type %s", to.Type())
			}
			if math.IsInf(floatVal, -1) {
				return errors.Errorf("cannot convert -Infinity to unsigned integer type %s", to.Type())
			}
		}

		// Check for negative float values
		if from < 0 {
			return errors.Errorf("cannot set negative value %v to unsigned type %s", from, to.Type())
		}

		// Check if value is too small and will be truncated to zero
		if floatVal > 0 && floatVal < 1.0 { //nolint:staticcheck // for future use
			// This is not an error, but a warning that could be logged in a real-world app
			// For now, just proceed with truncation
		}

		// Check if the value is within uint64 range
		if floatVal > float64(math.MaxUint64) {
			return errors.Errorf("value %v is outside the range of any unsigned integer type", floatVal)
		}

		// Convert to uint64 and check range for specific uint sizes
		uint64Val := uint64(floatVal)
		if err := checkUnsignedIntegerRange(to, uint64Val); err != nil {
			return err
		}

		to.SetUint(uint64Val)

		return nil
	case reflect.Float32:
		// Handle special float values separately
		if isSpecial {
			if math.IsNaN(floatVal) {
				to.SetFloat(math.NaN())
				return nil
			}

			if math.IsInf(floatVal, 1) {
				to.SetFloat(math.Inf(1))
				return nil
			}

			if math.IsInf(floatVal, -1) {
				to.SetFloat(math.Inf(-1))
				return nil
			}
		}

		// Check for potential float32 overflow
		if floatVal > math.MaxFloat32 || floatVal < -math.MaxFloat32 {
			return errors.Errorf("value %v is outside the range of float32", floatVal)
		}

		to.SetFloat(floatVal)

		return nil
	case reflect.Float64:
		// No need to check range for float64 as it can hold any float32 or float64 value
		to.SetFloat(floatVal)
		return nil

	case reflect.Complex64, reflect.Complex128:
		// Set float value to the real part, imaginary part is 0
		// Handle special values
		if isSpecial {
			var complexVal complex128
			switch {
			case math.IsNaN(floatVal):
				// For NaN, set both real and imaginary parts to NaN
				complexVal = complex(math.NaN(), math.NaN())
			case math.IsInf(floatVal, 1):
				complexVal = complex(math.Inf(1), 0)
			case math.IsInf(floatVal, -1):
				complexVal = complex(math.Inf(-1), 0)
			}
			to.SetComplex(complexVal)
			return nil
		}

		to.SetComplex(complex(floatVal, 0))
		return nil

	case reflect.Interface:
		// For any fields, use the original float type (either float32 or float64)
		to.Set(reflect.ValueOf(from))
		return nil

	case reflect.Ptr:
		// For pointer fields, check if valid and then dereference and call recursively
		if to.IsNil() {
			// Initialize nil pointers
			to.Set(reflect.New(to.Type().Elem()))
		}

		return Float[F](to.Elem(), from)
	}

	// If the to doesn't match any of the above types, return an error with more details
	return errors.Wrapf(ErrNotSupported, "cannot set float value %v to field of type %s", from, to.Type())
}

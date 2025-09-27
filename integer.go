package assign

import (
	"math"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// SignedValue is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type SignedValue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedValue is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type UnsignedValue interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// IntegerValue is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type IntegerValue interface {
	SignedValue | UnsignedValue
}

// Integer converts an integer value to the appropriate type for a target field.
// Handles conversion to strings, booleans, numeric types, and interfaces.
// Performs range checking to prevent overflow.
//
// Parameters:
//   - field: Target field to set (reflect.Value).
//   - number: Integer value to convert and set.
//
// Returns:
//   - error: If conversion or assignment fails.
func Integer[I IntegerValue](to reflect.Value, from I) error {
	// Check if the to is settable
	if !to.CanSet() {
		return errors.Errorf("field of type %s is not settable", to.Type())
	}

	// Determine if the from is signed or unsigned
	var (
		isSigned  bool
		int64Val  int64
		uint64Val uint64
	)

	// Convert from int64 or uint64 for standard processing
	switch any(from).(type) {
	case int, int8, int16, int32, int64:
		isSigned = true
		int64Val = int64(from)
	default:
		// Must be an unsigned integer
		uint64Val = uint64(from)
	}

	switch to.Kind() {
	case reflect.String:
		// Convert integer from string using appropriate formatting based on type
		if isSigned {
			to.SetString(strconv.FormatInt(int64Val, 10))
		} else {
			to.SetString(strconv.FormatUint(uint64Val, 10))
		}

		return nil

	case reflect.Bool:
		// Convert integer from bool (true only if > 0, false otherwise)
		if isSigned {
			to.SetBool(int64Val > 0)
		} else {
			to.SetBool(uint64Val > 0)
		}

		return nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if isSigned {
			// For signed input from signed to, check range
			if err := checkSignedIntegerRange(to, int64Val); err != nil {
				return err
			}
			to.SetInt(int64Val)
		} else {
			// For unsigned input from signed to, check overflow
			if uint64Val > math.MaxInt64 {
				return errors.Errorf("value %v overflows target type %s", uint64Val, to.Type())
			}
			// Then check range as a signed value
			if err := checkSignedIntegerRange(to, int64(uint64Val)); err != nil {
				return err
			}

			to.SetInt(int64(uint64Val))
		}

		return nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		// For unsigned types, ensure signed values are not negative
		if isSigned && int64Val < 0 {
			return errors.Errorf("cannot set negative value %v to unsigned type %s", int64Val, to.Type())
		}

		// Calculate the value from set
		var valueToSet uint64
		if isSigned {
			valueToSet = uint64(int64Val)
		} else {
			valueToSet = uint64Val
		}

		// Check for range
		if err := checkUnsignedIntegerRange(to, valueToSet); err != nil {
			return err
		}

		to.SetUint(valueToSet)
		return nil

	case reflect.Float32, reflect.Float64:
		// Convert integer from float
		if isSigned {
			to.SetFloat(float64(int64Val))
		} else {
			to.SetFloat(float64(uint64Val))
		}
		return nil

	case reflect.Complex64, reflect.Complex128:
		// Value integer value from the real part, imaginary part is 0
		if isSigned {
			to.SetComplex(complex(float64(int64Val), 0))
		} else {
			to.SetComplex(complex(float64(uint64Val), 0))
		}
		return nil

	case reflect.Interface:
		// For any fields, just set the integer value
		if isSigned {
			to.Set(reflect.ValueOf(int64Val))
		} else {
			to.Set(reflect.ValueOf(uint64Val))
		}
		return nil

	case reflect.Ptr:
		// For pointer fields, check if valid and then dereference and call recursively
		if to.IsNil() {
			// Initialize nil pointers
			to.Set(reflect.New(to.Type().Elem()))
		}
		return Integer[I](to.Elem(), from)
	}

	// If the to doesn't match any of the above types, return an error with more details
	var valueStr string
	if isSigned {
		valueStr = strconv.FormatInt(int64Val, 10)
	} else {
		valueStr = strconv.FormatUint(uint64Val, 10)
	}

	return errors.Wrapf(ErrNotSupported, "cannot set integer value %s to field of type %s", valueStr, to.Type())
}

// checkSignedIntegerRange verifies a signed integer value is within range for the target field.
// Prevents data loss or overflow from type conversions.
func checkSignedIntegerRange(field reflect.Value, number int64) error {
	// Store field kind to avoid multiple calls to field.Kind()
	kind := field.Kind()

	// Quick exit if we're not dealing with a signed integer type
	if kind < reflect.Int8 || (kind > reflect.Int64 && kind != reflect.Int) {
		return nil
	}

	// For signed values, check against the min/max range of the target type
	bitSize := field.Type().Bits()
	maxVal := int64(1)<<(bitSize-1) - 1
	minVal := -int64(1) << (bitSize - 1)

	if number > maxVal || number < minVal {
		return errors.Errorf("value %v is outside the range of target type %s [%d, %d]",
			number, field.Type(), minVal, maxVal)
	}

	return nil
}

// checkUnsignedIntegerRange verifies an unsigned integer value is within range for the target field.
// Checks for potential overflow based on bit size.
func checkUnsignedIntegerRange(field reflect.Value, number uint64) error {
	// For unsigned types, check for potential overflow
	bitSize := field.Type().Bits()
	maxVal := uint64(1)<<bitSize - 1

	if number > maxVal {
		return errors.Errorf("value %v overflows target type %s (max: %d)", number, field.Type(), maxVal)
	}

	return nil
}

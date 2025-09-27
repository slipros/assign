// Package assign provides utilities for type conversion and setting values in Go structs
// with automatic type conversion, validation, and nil pointer initialization.
//
// The library offers seamless conversion between compatible types, supports
// Go 1.18+ generics with type constraints, and provides an extension system
// for custom type converters. It's designed for high performance with
// optimized fast paths for common conversions.
//
// Example usage:
//
//	var target int
//	field := reflect.ValueOf(&target).Elem()
//	err := assign.Value(field, "42") // Converts string "42" to int 42
//
// The library supports:
//   - Automatic type conversion between compatible types
//   - Generic type constraints for compile-time safety
//   - Nil pointer initialization when needed
//   - TextUnmarshaler and BinaryUnmarshaler interface support
//   - Extension functions for custom type conversions
//   - Optimized performance with fast paths for common types
package assign

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// ExtensionFunc defines a function type for custom type conversion extensions.
// It takes a value of any type and returns a conversion function and a boolean
// indicating whether the conversion is supported.
//
// The returned function performs the actual conversion and returns an error if it fails.
// The boolean indicates whether this extension can handle the given value type.
type ExtensionFunc func(any) (func(to reflect.Value) error, bool)

// Value assigns a value to a reflect.Value field with automatic type conversion.
// This is the primary entry point for the assign library, handling initialization
// of nil pointers, various primitive types, and interfaces.
//
// The function automatically determines the best conversion strategy based on
// the source and target types. Extension functions can be provided to handle
// custom type conversions that aren't built into the library.
//
// Key features:
//   - Converts between compatible types (string ↔ int, float ↔ string, etc.)
//   - Initializes nil pointers when needed
//   - Supports fmt.Stringer interface for string conversion
//   - Handles slices, arrays, and common primitive types
//   - Uses extension functions for custom type conversions
//   - Provides detailed error messages for debugging
//
// Example:
//
//	var target struct {
//		ID   int
//		Name string
//		Age  *int
//	}
//	v := reflect.ValueOf(&target).Elem()
//
//	// Convert string to int
//	assign.Value(v.FieldByName("ID"), "123")
//
//	// Set string directly
//	assign.Value(v.FieldByName("Name"), "John")
//
//	// Initialize nil pointer and set value
//	assign.Value(v.FieldByName("Age"), "30")
//
// Parameters:
//   - to: Target field to set (must be settable reflect.Value)
//   - from: Source value to convert and assign (any type)
//   - extensions: Optional extension functions for custom type conversions
//
// Returns:
//   - error: If value could not be set due to type incompatibility or other issues
func Value(to reflect.Value, from any, extensions ...ExtensionFunc) error {
	// Check if to can be set
	if !to.CanSet() {
		return errors.Errorf("field of type %s is not settable", to.Type())
	}

	// Handle nil from early
	if from == nil {
		to.Set(reflect.Zero(to.Type()))

		return nil
	}

	kind := to.Kind()

	// Initialize and dereference nil pointers recursively
	if kind == reflect.Pointer {
		if to.IsNil() {
			// Initialize the pointer with a new from of the appropriate type
			to.Set(reflect.New(to.Type().Elem()))
		}

		// Recursively call Value on the dereferenced pointer
		return Value(to.Elem(), from)
	}

	for _, ext := range extensions {
		assign, ok := ext(from)
		if !ok {
			continue
		}

		return assign(to)
	}

	// Handle various types using specialized setters
	switch t := from.(type) {
	case string:
		return String(to, t)
	case *string:
		if t == nil {
			// For nil string pointers, set zero from
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return String(to, *t)
	case bool:
		// Add direct support for boolean values
		if kind == reflect.Bool {
			to.SetBool(t)
			return nil
		}
		return String(to, strconv.FormatBool(t))
	case *bool:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		if kind == reflect.Bool {
			to.SetBool(*t)
			return nil
		}
		return String(to, strconv.FormatBool(*t))
	case int:
		return Integer(to, t)
	case *int:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case int8:
		return Integer(to, t)
	case *int8:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case int16:
		return Integer(to, t)
	case *int16:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case int32:
		return Integer(to, t)
	case *int32:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case int64:
		return Integer(to, t)
	case *int64:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case uint:
		return Integer(to, t)
	case *uint:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case uint8:
		return Integer(to, t)
	case *uint8:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case uint16:
		return Integer(to, t)
	case *uint16:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case uint32:
		return Integer(to, t)
	case *uint32:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case uint64:
		return Integer(to, t)
	case *uint64:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Integer(to, *t)
	case float32:
		return Float(to, t)
	case *float32:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Float(to, *t)
	case float64:
		return Float(to, t)
	case *float64:
		if t == nil {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		return Float(to, *t)
	case []string:
		return SliceString(to, t)
	case []any:
		// Handle []any differently based on the to type
		if kind == reflect.Slice {
			return handleInterfaceSlice(to, t)
		}
	}

	// Handle types that implement fmt.Stringer
	if stringer, ok := from.(fmt.Stringer); ok {
		return String(to, stringer.String())
	}

	// Handle general assignable types
	valueValue := reflect.ValueOf(from)
	valueType := valueValue.Type()

	// Handle pointers for general types
	if valueType.Kind() == reflect.Pointer {
		// Check if the pointer is nil
		if valueValue.IsNil() {
			to.Set(reflect.Zero(to.Type()))

			return nil
		}

		// Dereference pointer and use indirect from
		valueType = valueType.Elem()
		valueValue = valueValue.Elem()
	}

	fieldType := to.Type()

	// Check if the from's type can be assigned to the target to's type
	if valueType.AssignableTo(fieldType) {
		to.Set(valueValue)

		return nil
	}

	// Try if the from's type is convertible to the target to's type
	if valueType.ConvertibleTo(fieldType) {
		// Perform safe conversion with panic recovery
		var convertErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					convertErr = errors.Errorf("conversion panic: %v", r)
				}
			}()

			convertedValue := valueValue.Convert(fieldType)
			to.Set(convertedValue)
		}()

		if convertErr != nil {
			return convertErr
		}

		return nil
	}

	// If we reach here, the from couldn't be set
	return errors.Wrapf(ErrNotSupported, "cannot set value of type %T to field of type %s", from, fieldType)
}

// handleInterfaceSlice handles conversion from []any to a target slice type.
// This function creates a new slice of the appropriate target type and converts
// each element from the any slice to the target element type.
//
// Parameters:
//   - field: The target slice field to populate (must be of slice kind).
//   - values: Slice of any values to convert and assign.
//
// Returns:
//   - error: An error if conversion fails for any element, or nil if successful.
func handleInterfaceSlice(field reflect.Value, values []any) error {
	if field.Kind() != reflect.Slice {
		return errors.Wrapf(ErrNotSupported, "cannot set []any to non-slice field of type %s", field.Type())
	}

	// Create a new slice of the appropriate type
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))

	// Convert each element and set directly by index
	for i, val := range values {
		elem := reflect.New(elemType).Elem()
		if err := Value(elem, val); err != nil {
			return errors.Wrapf(err, "failed to convert element %d of []any to %s", i, elemType)
		}

		slice.Index(i).Set(elem)
	}

	field.Set(slice)

	return nil
}

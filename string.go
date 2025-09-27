package assign

import (
	"encoding"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

// Pre-allocated reflect.Type values for common types to avoid repeated Type() calls
var byteType = reflect.TypeOf(byte(0))

// String converts a string value to the appropriate type for a target field.
// This is an optimized version that uses fast paths for common type conversions
// while maintaining full compatibility with all supported types.
// Handles conversion to numeric types, booleans, slices, time.Time and other types.
// Supports types implementing TextUnmarshaler or BinaryUnmarshaler interfaces.
//
// Parameters:
//   - field: Target field to set (reflect.Value).
//   - str: String value to convert and set.
//
// Returns:
//   - error: If conversion or assignment fails.
func String(to reflect.Value, from string) error {
	if !to.CanSet() {
		return errors.Errorf("field of type %s is not settable", to.Type())
	}

	if from == "" {
		return handleEmptyString(to)
	}

	kind := to.Kind()

	// Fast paths for most common types
	switch kind {
	case reflect.String:
		to.SetString(from)
		return nil

	case reflect.Int:
		return setIntFromString(to, from, 0)
	case reflect.Int8:
		return setIntFromString(to, from, 8)
	case reflect.Int16:
		return setIntFromString(to, from, 16)
	case reflect.Int32:
		return setIntFromString(to, from, 32)
	case reflect.Int64:
		return setIntFromString(to, from, 64)

	case reflect.Bool:
		return setBoolFromString(to, from)

	case reflect.Float32:
		return setFloatFromString(to, from, 32)
	case reflect.Float64:
		return setFloatFromString(to, from, 64)

	case reflect.Slice:
		return setSliceFromString(to, from)

	case reflect.Ptr:
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		return String(to.Elem(), from)

	case reflect.Interface:
		to.Set(reflect.ValueOf(from))
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintFromString(to, from)

	case reflect.Complex64, reflect.Complex128:
		return setComplexFromString(to, from)

	case reflect.Struct:
		return setStructFromString(to, from)

	case reflect.Map:
		return setMapFromString(to, from)
	}

	// Try to use custom unmarshalers for other types
	return tryUnmarshalers(to, from)
}

// handleEmptyString handles empty strings with optimized zero value assignment
func handleEmptyString(field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString("")
		return nil
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// Use unsafe zero assignment for primitive types
		field.Set(reflect.Zero(field.Type()))
		return nil
	case reflect.Slice:
		// Don't allocate for empty slices
		return nil
	case reflect.Interface:
		field.Set(reflect.ValueOf(""))
		return nil
	case reflect.Map:
		// For maps, create a new map if nil
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}
		return nil
	case reflect.Ptr:
		// For pointers, set to a new zero value if nil
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return nil
	case reflect.Complex64, reflect.Complex128:
		// For numeric and bool types, empty string means zero value
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	// Try unmarshalers with empty string
	return tryUnmarshalers(field, "")
}

// setIntFromString provides optimized integer parsing with comprehensive security validation.
// This function implements proper bitSize validation for all code paths to prevent
// silent data corruption and security vulnerabilities.
func setIntFromString(field reflect.Value, str string, bitSize int) error {
	// Security Fix: Handle common single-digit cases WITH bitSize validation
	// Previous vulnerability: single digits bypassed all validation
	if len(str) == 1 && str[0] >= '0' && str[0] <= '9' {
		val := int64(str[0] - '0')
		// Critical: Validate bitSize constraints for single digits too
		if err := validateIntRange(val, bitSize); err != nil {
			return err
		}
		field.SetInt(val)
		return nil
	}

	// Check for special prefixes first before doing simple digit parsing
	// Determine base automatically (0 means auto-detect: 0x for hex, 0 for octal)
	base := 10
	if len(str) >= 2 {
		switch {
		case str[0] == '0' && (str[1] == 'x' || str[1] == 'X'):
			base = 0 // Auto-detect will use base 16 for 0x prefix
		case str[0] == '0' && (str[1] == 'b' || str[1] == 'B'):
			base = 0 // Auto-detect will use base 2 for 0b prefix
		case str[0] == '0' && isAllDigits(str[1:]):
			base = 0 // Auto-detect will use base 8 for 0 prefix
		}
	}

	// Security Fix: Handle simple multi-digit positive numbers WITH proper validation
	// Only do this if we're using base 10 (no special prefixes) and no negative sign
	if base == 10 && len(str) <= 3 && isAllDigits(str) {
		val := int64(0)
		for i := 0; i < len(str); i++ {
			val = val*10 + int64(str[i]-'0')
		}

		// Critical Security Fix: Always validate bitSize constraints
		// Previous vulnerability: incomplete overflow checking with logic errors
		if err := validateIntRange(val, bitSize); err != nil {
			return err
		}

		field.SetInt(val)
		return nil
	}

	// Fall back to strconv for complex cases (negative numbers, hex, octal, large values)
	parsed, err := strconv.ParseInt(str, base, bitSize)
	if err != nil {
		return errors.Wrapf(err, "cannot convert string '%s' to int", str)
	}
	field.SetInt(parsed)
	return nil
}

// validateIntRange validates that a value fits within the specified integer bit size.
// This function prevents silent data corruption and security vulnerabilities
// by ensuring proper range checking for all integer types.
func validateIntRange(val int64, bitSize int) error {
	if bitSize == 0 {
		// For int type (bitSize 0), we accept any value that fits in int64
		// as the actual bit size depends on the platform (32 or 64 bit)
		return nil
	}

	// Define valid ranges for each bit size
	var minVal, maxVal int64
	switch bitSize {
	case 8:
		minVal, maxVal = math.MinInt8, math.MaxInt8 // -128 to 127
	case 16:
		minVal, maxVal = math.MinInt16, math.MaxInt16 // -32768 to 32767
	case 32:
		minVal, maxVal = math.MinInt32, math.MaxInt32 // -2147483648 to 2147483647
	case 64:
		minVal, maxVal = math.MinInt64, math.MaxInt64 // Full int64 range
	default:
		return errors.Errorf("unsupported integer bit size: %d", bitSize)
	}

	// Security validation: ensure value is within valid range
	if val < minVal || val > maxVal {
		return errors.Errorf("value %d out of range for %d-bit signed integer (range: %d to %d)", val, bitSize, minVal, maxVal)
	}

	return nil
}

// setBoolFromString provides optimized boolean parsing for common values
func setBoolFromString(field reflect.Value, str string) error {
	// Handle common string boolean representations with optimized paths
	switch len(str) {
	case 1:
		switch str[0] {
		case '1', 't', 'T', 'y', 'Y':
			field.SetBool(true)
			return nil
		case '0', 'f', 'F', 'n', 'N':
			field.SetBool(false)
			return nil
		}
	case 2:
		if equalFold(str, "on") || equalFold(str, "no") {
			field.SetBool(str[0] == 'o' || str[0] == 'O') // "on" = true, "no" = false
			return nil
		}
	case 3:
		if equalFold(str, "yes") || equalFold(str, "off") {
			field.SetBool(str[0] == 'y' || str[0] == 'Y') // "yes" = true, "off" = false
			return nil
		}
	case 4:
		if equalFold(str, "true") {
			field.SetBool(true)
			return nil
		}
	case 5:
		if equalFold(str, "false") {
			field.SetBool(false)
			return nil
		}
	}

	// Fall back to standard parsing
	parsed, err := strconv.ParseBool(str)
	if err != nil {
		return errors.Wrapf(err, "cannot convert string '%s' to bool", str)
	}
	field.SetBool(parsed)
	return nil
}

// setFloatFromString provides optimized float parsing
func setFloatFromString(field reflect.Value, str string, bitSize int) error {
	parsed, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return errors.Wrapf(err, "cannot convert string '%s' to float", str)
	}
	field.SetFloat(parsed)
	return nil
}

// setSliceFromString handles slice conversion with pre-allocation and optimized splitting
func setSliceFromString(field reflect.Value, str string) error {
	elemType := field.Type().Elem()

	// Optimize []byte case
	if elemType == byteType {
		field.SetBytes([]byte(str))
		return nil
	}

	// For string slices, use optimized splitting
	if elemType.Kind() == reflect.String {
		// Estimate parts based on comma count + 1
		estimatedParts := 1
		for i := 0; i < len(str); i++ {
			if str[i] == ',' {
				estimatedParts++
			}
		}

		parts := strings.Split(str, ",")
		slice := reflect.MakeSlice(field.Type(), 0, len(parts))

		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}

			// Create a properly typed reflect.Value for string element
			stringValue := reflect.ValueOf(trimmed).Convert(elemType)
			slice = reflect.Append(slice, stringValue)
		}

		field.Set(slice)

		return nil
	}

	// For other slice types, try to parse the string as a comma-separated list
	if strings.Contains(str, ",") {
		parts := strings.Split(str, ",")
		// Create a new slice with initial capacity
		slice := reflect.MakeSlice(field.Type(), 0, len(parts))

		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}

			// Create a new element for the slice
			newElem := reflect.New(elemType).Elem()

			// Try to set the new element with the string value
			if err := String(newElem, trimmed); err != nil {
				return err
			}

			slice = reflect.Append(slice, newElem)
		}

		field.Set(slice)

		return nil
	}

	// Try to set a single element
	newElem := reflect.New(elemType).Elem()
	if err := String(newElem, str); err != nil {
		return err
	}

	// Create a new slice with the single element
	slice := reflect.MakeSlice(field.Type(), 0, 1)
	slice = reflect.Append(slice, newElem)

	field.Set(slice)

	return nil
}

// setUintFromString converts a string to an unsigned integer and sets the field value.
// It handles decimal, hexadecimal (0x prefix), and octal (0 prefix) formats.
func setUintFromString(field reflect.Value, str string) error {
	// Determine base automatically
	base := 10
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		base = 0 // Auto-detect will use base 16 for 0x prefix
	} else if len(str) > 1 && strings.HasPrefix(str, "0") {
		base = 0 // Auto-detect will use base 8 for 0 prefix
	}

	// Get bit size based on field type
	bitSize := 0 // Default for uint
	switch field.Kind() {
	case reflect.Uint8:
		bitSize = 8
	case reflect.Uint16:
		bitSize = 16
	case reflect.Uint32:
		bitSize = 32
	case reflect.Uint64:
		bitSize = 64
	}

	// Parse the string
	parsed, err := strconv.ParseUint(str, base, bitSize)
	if err != nil {
		return errors.Wrapf(err, "cannot convert string '%s' to %s", str, field.Kind())
	}

	// Value the value
	field.SetUint(parsed)
	return nil
}

// setComplexFromString converts a string to a complex number and sets the field value.
func setComplexFromString(field reflect.Value, str string) error {
	// Get bit size based on field type
	bitSize := 128 // Default for complex128
	if field.Kind() == reflect.Complex64 {
		bitSize = 64
	}

	// Parse the string
	parsed, err := strconv.ParseComplex(str, bitSize)
	if err != nil {
		return errors.Wrapf(err, "cannot convert string '%s' to %s", str, field.Kind())
	}

	// Value the value
	field.SetComplex(parsed)
	return nil
}

// setStructFromString handles conversion from string to struct types.
// Currently supports time.Time directly.
func setStructFromString(field reflect.Value, str string) error {
	// Special handling for time.Time
	if field.Type() == typeTime {
		return setTimeFromString(field, str)
	}

	// For other structs, try using unmarshalers
	return tryUnmarshalers(field, str)
}

// setTimeFromString parses a string into a time.Time value using various layouts.
func setTimeFromString(field reflect.Value, str string) error {
	t, err := parseTime(str)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(t))

	return nil
}

// setMapFromString attempts to set a map value from a string.
// Currently, this only supports maps with string keys and a simple format like "key1:value1,key2:value2".
func setMapFromString(field reflect.Value, str string) error {
	// Only support maps with string keys for now
	if field.Type().Key().Kind() != reflect.String {
		return errors.Wrapf(ErrNotSupported, "only maps with string keys are supported for string conversion")
	}

	// Check if the map needs to be initialized
	if field.IsNil() {
		field.Set(reflect.MakeMap(field.Type()))
	}

	// Parse key-value pairs (format: "key1:value1,key2:value2")
	if !strings.Contains(str, ":") {
		return errors.Wrapf(ErrNotSupported, "invalid map format, expected 'key:value' pairs separated by commas")
	}

	pairs := strings.Split(str, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return errors.Errorf("invalid map key-value pair: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// Create a new value for the map
		elemType := field.Type().Elem()
		elem := reflect.New(elemType).Elem()

		// Value the value
		if err := String(elem, value); err != nil {
			return errors.Wrapf(err, "cannot convert '%s' to map value type %s", value, elemType)
		}

		// Value the key-value pair in the map
		field.SetMapIndex(reflect.ValueOf(key), elem)
	}

	return nil
}

// tryUnmarshalers attempts to use TextUnmarshaler or BinaryUnmarshaler interfaces
// to unmarshal the string into the field.
func tryUnmarshalers(field reflect.Value, str string) error {
	// For types that implement TextUnmarshaler or BinaryUnmarshaler,
	// we need to get a pointer to the field
	if !field.CanAddr() {
		return errors.WithStack(ErrNotSupported)
	}

	ptr := field.Addr()
	if !ptr.CanInterface() {
		return errors.WithStack(ErrNotSupported)
	}

	return implementsBytesUnmarshaler(ptr.Interface(), str)
}

// implementsBytesUnmarshaler uses TextUnmarshaler or BinaryUnmarshaler interfaces
// to unmarshal a string into a value, if supported.
func implementsBytesUnmarshaler(ptr any, str string) error {
	switch i := ptr.(type) {
	case encoding.TextUnmarshaler:
		// For types like time.Time that implement TextUnmarshaler
		if err := i.UnmarshalText([]byte(str)); err != nil {
			return errors.Wrapf(err, "TextUnmarshaler failed for '%s'", str)
		}
		return nil

	case encoding.BinaryUnmarshaler:
		// For types that implement BinaryUnmarshaler
		if err := i.UnmarshalBinary([]byte(str)); err != nil {
			return errors.Wrapf(err, "BinaryUnmarshaler failed for '%s'", str)
		}
		return nil
	}

	return errors.Wrapf(ErrNotSupported, "type %T does not implement UnmarshalText or UnmarshalBinary", ptr)
}

// isAllDigits checks if string contains only digits
func isAllDigits(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}

	return true
}

// equalFold performs case-insensitive string comparison for ASCII strings.
// This function is optimized for parsing common boolean string representations
// and only handles ASCII characters. For comprehensive Unicode support,
// use strings.EqualFold from the standard library.
//
// The unsafe operations are justified here because:
// 1. This is a performance-critical path for string comparison operations
// 2. We only read from string data (no modification)
// 3. Length bounds are checked before unsafe access
// 4. Only used with known, short ASCII strings like "true", "false", etc.
//
//nolint:gosec // G103: Justified use of unsafe for performance in hot path
func equalFold(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	// For empty strings or single characters, avoid unsafe overhead
	if len(s1) <= 1 {
		return strings.EqualFold(s1, s2)
	}

	// Use unsafe for performance with ASCII strings
	// This is safe because we only read and lengths are verified
	s1bytes := unsafe.Slice(unsafe.StringData(s1), len(s1))
	s2bytes := unsafe.Slice(unsafe.StringData(s2), len(s2))

	for i := 0; i < len(s1); i++ {
		c1, c2 := s1bytes[i], s2bytes[i]
		if c1 != c2 {
			// Simple ASCII case folding - convert uppercase to lowercase
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				return false
			}
		}
	}
	return true
}

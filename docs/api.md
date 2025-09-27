# API Reference

## Overview

The assign library provides a comprehensive set of functions for type conversion and value assignment in Go. The API is designed around reflection-based value assignment with automatic type conversion, generic type constraints for safety, and extensibility through custom conversion functions.

## Core Functions

### Value

```go
func Value(to reflect.Value, from any, extensions ...ExtensionFunc) error
```

The primary entry point for type conversion and assignment. This function handles the conversion and assignment of any value type to a target field with automatic type detection and conversion.

**Parameters:**
- `to` (reflect.Value): Target field that must be settable
- `from` (any): Source value of any type to be converted and assigned
- `extensions` (...ExtensionFunc): Optional custom conversion functions

**Returns:**
- `error`: Returns an error if the conversion fails or the target field is not settable

**Supported Conversions:**
- String to numeric types (int, uint, float, complex)
- String to boolean with flexible format support
- String to time.Time with multiple layout detection
- String to slices with CSV parsing
- Numeric type conversions with overflow protection
- Boolean conversions from various string representations
- Nil pointer initialization and assignment
- Interface{} assignments with type preservation

**Example:**
```go
var target struct {
    ID   int
    Name string
    Age  *int
}
v := reflect.ValueOf(&target).Elem()

assign.Value(v.FieldByName("ID"), "123")     // String to int
assign.Value(v.FieldByName("Name"), "John")  // String assignment
assign.Value(v.FieldByName("Age"), "30")     // Nil pointer initialization + assignment
```

### Integer

```go
func Integer[I IntegerValue](to reflect.Value, from I) error
```

Converts integer values to appropriate target types with comprehensive overflow protection and range validation.

**Type Constraint:**
- `I` must satisfy `IntegerValue` interface (any signed or unsigned integer type)

**Parameters:**
- `to` (reflect.Value): Target field for integer assignment
- `from` (I): Source integer value of any integer type

**Returns:**
- `error`: Returns an error if conversion fails, target is not settable, or value is out of range

**Features:**
- Automatic range checking for target type
- Overflow protection with detailed error messages
- Support for all Go integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
- Conversion to string, boolean, float, and complex types
- Interface{} assignment with type preservation

**Example:**
```go
var targets struct {
    Small  int8
    Medium int16
    Large  int64
}
v := reflect.ValueOf(&targets).Elem()

assign.Integer(v.FieldByName("Small"), 42)        // int to int8
assign.Integer(v.FieldByName("Medium"), int32(1000)) // int32 to int16
assign.Integer(v.FieldByName("Large"), uint64(999))  // uint64 to int64 (with checks)

// This would return an error:
// assign.Integer(v.FieldByName("Small"), 300) // Value 300 > int8 max (127)
```

### Float

```go
func Float[F FloatValue](to reflect.Value, from F) error
```

Handles floating-point conversions with special value support (NaN, Infinity) and proper range checking.

**Type Constraint:**
- `F` must satisfy `FloatValue` interface (float32 or float64)

**Parameters:**
- `to` (reflect.Value): Target field for float assignment
- `from` (F): Source floating-point value

**Returns:**
- `error`: Returns an error if conversion fails or target is not settable

**Features:**
- Special value handling (NaN, +Inf, -Inf)
- Conversion to string with appropriate formatting
- Conversion to integer types with truncation warnings
- Conversion to boolean (non-zero = true)
- Complex number conversion (real part only)
- Range checking for float32 targets

**Example:**
```go
var targets struct {
    Regular float32
    Double  float64
    Text    string
    Number  int
}
v := reflect.ValueOf(&targets).Elem()

assign.Float(v.FieldByName("Regular"), 3.14159)      // float64 to float32
assign.Float(v.FieldByName("Double"), 2.71828)       // Direct assignment
assign.Float(v.FieldByName("Text"), 123.456)         // Float to string
assign.Float(v.FieldByName("Number"), 98.6)          // Float to int (truncated)

// Special values
assign.Float(v.FieldByName("Double"), math.NaN())    // NaN handling
assign.Float(v.FieldByName("Double"), math.Inf(1))   // +Infinity handling
```

### String

```go
func String(to reflect.Value, from string) error
```

Optimized string parsing and conversion with fast paths for common type conversions.

**Parameters:**
- `to` (reflect.Value): Target field for string conversion
- `from` (string): Source string value

**Returns:**
- `error`: Returns an error if conversion fails or target is not settable

**Features:**
- Fast paths for common conversions (single digits, short strings)
- Comprehensive boolean string recognition ("true", "1", "yes", "on", etc.)
- Automatic base detection for integers (decimal, hex, octal, binary)
- CSV parsing for slice types
- Time parsing with multiple format support
- Support for TextUnmarshaler and BinaryUnmarshaler interfaces
- Map parsing with key:value format
- Complex number parsing

**Example:**
```go
var targets struct {
    Number   int
    Active   bool
    Score    float64
    Items    []string
    Bytes    []byte
    Time     time.Time
    Complex  complex128
}
v := reflect.ValueOf(&targets).Elem()

assign.String(v.FieldByName("Number"), "42")                    // String to int
assign.String(v.FieldByName("Active"), "true")                  // String to bool
assign.String(v.FieldByName("Score"), "98.6")                   // String to float
assign.String(v.FieldByName("Items"), "apple,banana,cherry")    // CSV to slice
assign.String(v.FieldByName("Bytes"), "hello")                  // String to []byte
assign.String(v.FieldByName("Time"), "2023-01-15T10:30:00Z")    // String to time.Time
assign.String(v.FieldByName("Complex"), "3+4i")                 // String to complex
```

**Boolean String Recognition:**
- True values: "true", "1", "t", "T", "yes", "y", "Y", "on"
- False values: "false", "0", "f", "F", "no", "n", "N", "off"

**Integer Base Detection:**
- Decimal: "123", "-456"
- Hexadecimal: "0x1A", "0X1a"
- Octal: "0123"
- Binary: "0b1010", "0B1010"

### SliceString

```go
func SliceString(to reflect.Value, from []string, options ...SliceOptionFunc) error
```

Converts string slices to various target types with configurable separator options.

**Parameters:**
- `to` (reflect.Value): Target field for slice assignment
- `from` ([]string): Source string slice
- `options` (...SliceOptionFunc): Optional configuration functions

**Returns:**
- `error`: Returns an error if conversion fails or target is not settable

**Features:**
- Conversion to numeric slices ([]int, []float64, etc.)
- String joining with configurable separators
- Interface{} slice assignment
- Element-wise type conversion with error reporting
- Support for pointer target types

**Example:**
```go
var targets struct {
    Numbers []int
    Joined  string
    Custom  string
}
v := reflect.ValueOf(&targets).Elem()
stringSlice := []string{"1", "2", "3", "4", "5"}

// Convert to int slice
assign.SliceString(v.FieldByName("Numbers"), stringSlice)

// Join with default separator (comma)
assign.SliceString(v.FieldByName("Joined"), []string{"apple", "banana", "cherry"})

// Join with custom separator
assign.SliceString(v.FieldByName("Custom"), []string{"red", "green", "blue"},
    assign.WithSeparator(" | "))
```

## Type Constraints

### SignedValue

```go
type SignedValue interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

Constraint for signed integer types. Used with the `Integer` function to ensure type safety.

### UnsignedValue

```go
type UnsignedValue interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
```

Constraint for unsigned integer types. Used with the `Integer` function to ensure type safety.

### IntegerValue

```go
type IntegerValue interface {
    SignedValue | UnsignedValue
}
```

Combined constraint for all integer types (signed and unsigned).

### FloatValue

```go
type FloatValue interface {
    ~float32 | ~float64
}
```

Constraint for floating-point types. Used with the `Float` function to ensure type safety.

## Extension System

### ExtensionFunc

```go
type ExtensionFunc func(any) (func(to reflect.Value) error, bool)
```

Function type for creating custom type conversion extensions.

**Parameters:**
- Input: `any` - The source value to potentially convert

**Returns:**
- `func(to reflect.Value) error` - Conversion function to execute if this extension handles the type
- `bool` - Whether this extension can handle the given value type

**Example:**
```go
import (
    "net/http"
    "reflect"
    "github.com/slipros/assign"
)

// Extension function for HTTP Cookie conversion
func cookieExtension(value any) (func(to reflect.Value) error, bool) {
    cookie, ok := value.(*http.Cookie)
    if !ok {
        return nil, false
    }

    return func(to reflect.Value) error {
        return assign.String(to, cookie.Value)
    }, true
}

// Usage
cookie := &http.Cookie{Name: "session", Value: "abc123"}
var sessionID string
field := reflect.ValueOf(&sessionID).Elem()

if err := assign.Value(field, cookie, cookieExtension); err != nil {
    panic(err)
}
// sessionID now contains "abc123"
```

## Configuration Options

### SliceOptionFunc

```go
type SliceOptionFunc func(*sliceOptions)
```

Function type for configuring `SliceString` behavior.

### WithSeparator

```go
func WithSeparator(sep string) SliceOptionFunc
```

Sets a custom separator for joining string slices.

**Parameters:**
- `sep` (string): Custom separator string

**Returns:**
- `SliceOptionFunc`: Configuration function for `SliceString`

**Example:**
```go
// Use pipe separator instead of comma
assign.SliceString(field, []string{"a", "b", "c"}, assign.WithSeparator(" | "))
// Result: "a | b | c"
```

## Error Types

### ErrNotSupported

```go
var ErrNotSupported = errors.New("not supported")
```

Sentinel error returned when a type conversion is not supported by the available conversion functions.

**Usage:**
```go
err := assign.Value(field, someValue)
if errors.Is(err, assign.ErrNotSupported) {
    // Handle unsupported conversion
}
```

## Interface Support

The library automatically detects and uses these standard Go interfaces:

### encoding.TextUnmarshaler

Types implementing `UnmarshalText([]byte) error` will be automatically detected and used for string conversions.

### encoding.BinaryUnmarshaler

Types implementing `UnmarshalBinary([]byte) error` will be automatically detected and used for byte slice conversions.

### fmt.Stringer

Types implementing `String() string` will be automatically converted to strings when needed.

## Performance Notes

- **Fast Paths**: Common conversions use optimized code paths
- **Zero Allocations**: Primitive type conversions avoid memory allocations
- **Caching**: Time format detection uses thread-safe caching
- **Regex Compilation**: Time parsing regex patterns are pre-compiled
- **String Operations**: Optimized for minimal allocations during string processing

## Thread Safety

All functions in the assign library are thread-safe and can be used concurrently from multiple goroutines. Internal caching mechanisms use appropriate synchronization primitives.
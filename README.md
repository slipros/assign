# assign

[![Go Reference](https://pkg.go.dev/badge/github.com/slipros/assign.svg)](https://pkg.go.dev/github.com/slipros/assign)
[![Go Report Card](https://goreportcard.com/badge/github.com/slipros/assign)](https://goreportcard.com/report/github.com/slipros/assign)
[![Coverage Status](https://coveralls.io/repos/github/slipros/assign/badge.svg)](https://coveralls.io/github/slipros/assign)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/slipros/assign)](https://golang.org/doc/devel/release.html)

A powerful Go library for type conversion and value assignment with automatic type conversion, validation, and nil pointer initialization.

## Features

- **Automatic Type Conversion**: Seamlessly convert between compatible types (string ↔ int, float ↔ string, etc.)
- **Generic Type Safety**: Leverage Go 1.18+ generics with type constraints for compile-time safety
- **Nil Pointer Handling**: Automatically initialize nil pointers when needed
- **Interface Support**: Built-in support for `encoding.TextUnmarshaler` and `encoding.BinaryUnmarshaler`
- **Extension System**: Add custom type converters through extension functions
- **High Performance**: Optimized fast paths for common type conversions
- **Comprehensive Time Parsing**: Support for multiple time formats with intelligent layout detection
- **Zero Dependencies**: Only depends on `github.com/pkg/errors` for enhanced error handling
- **Battle Tested**: 95.2% test coverage with comprehensive edge case handling

## Installation

```bash
go get github.com/slipros/assign@latest
```

## Quick Start

```go
package main

import (
    "fmt"
    "reflect"
    "time"

    "github.com/slipros/assign"
)

type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Age      int       `json:"age"`
    IsActive bool      `json:"is_active"`
    JoinedAt time.Time `json:"joined_at"`
    Scores   []int     `json:"scores"`
}

func main() {
    var user User
    v := reflect.ValueOf(&user).Elem()

    // Convert string to int
    assign.Value(v.FieldByName("ID"), "123")

    // Set string directly
    assign.Value(v.FieldByName("Name"), "John Doe")

    // Convert string to int
    assign.Value(v.FieldByName("Age"), "30")

    // Convert string to bool
    assign.Value(v.FieldByName("IsActive"), "true")

    // Parse time from string
    assign.Value(v.FieldByName("JoinedAt"), "2023-01-15T10:30:00Z")

    // Convert string slice to int slice
    assign.SliceString(v.FieldByName("Scores"), []string{"85", "92", "78"})

    fmt.Printf("User: %+v\n", user)
    // Output: User: {ID:123 Name:John Doe Age:30 IsActive:true JoinedAt:2023-01-15 10:30:00 +0000 UTC Scores:[85 92 78]}
}
```

## Core Functions

### Value - Universal Assignment

The `Value` function is the primary entry point for type conversion and assignment:

```go
func Value(to reflect.Value, from any, extensions ...ExtensionFunc) error
```

**Example:**
```go
var target int
field := reflect.ValueOf(&target).Elem()

// Convert from various types
assign.Value(field, "42")        // string to int
assign.Value(field, 3.14)       // float to int
assign.Value(field, true)       // bool to int (true = 1)
assign.Value(field, []byte("5")) // []byte to int
```

### Integer - Numeric Type Conversion

```go
func Integer[I IntegerValue](to reflect.Value, from I) error
```

Supports all integer types with overflow protection:

```go
var target int8
field := reflect.ValueOf(&target).Elem()

assign.Integer(field, 42)     // int to int8
assign.Integer(field, int64(100)) // int64 to int8

// Overflow protection
assign.Integer(field, 300)    // Returns error: value outside range
```

### Float - Floating Point Conversion

```go
func Float[F FloatValue](to reflect.Value, from F) error
```

Handles special float values (NaN, Infinity) appropriately:

```go
var target float32
field := reflect.ValueOf(&target).Elem()

assign.Float(field, 3.14159)         // float64 to float32
assign.Float(field, math.Inf(1))     // Handles +Infinity
assign.Float(field, math.NaN())      // Handles NaN
```

### String - Text Conversion

```go
func String(to reflect.Value, from string) error
```

Optimized string parsing with fast paths for common types:

```go
var targets struct {
    Number  int
    Active  bool
    Score   float64
    Items   []string
}

v := reflect.ValueOf(&targets).Elem()

assign.String(v.FieldByName("Number"), "123")
assign.String(v.FieldByName("Active"), "true")
assign.String(v.FieldByName("Score"), "95.5")
assign.String(v.FieldByName("Items"), "apple,banana,cherry") // CSV parsing
```

### SliceString - Array and Slice Conversion

```go
func SliceString(to reflect.Value, from []string, options ...SliceOptionFunc) error
```

Convert string slices to various types with configurable separators:

```go
var intSlice []int
var joinedString string

field1 := reflect.ValueOf(&intSlice).Elem()
field2 := reflect.ValueOf(&joinedString).Elem()

// String slice to int slice
assign.SliceString(field1, []string{"1", "2", "3", "4"})

// String slice to joined string
assign.SliceString(field2, []string{"apple", "banana", "cherry"})
// Result: "apple,banana,cherry"

// Custom separator
assign.SliceString(field2, []string{"a", "b", "c"}, assign.WithSeparator(" | "))
// Result: "a | b | c"
```

## Advanced Features

### Extension Functions

Create custom type converters for specialized types:

```go
// Custom type
type UserID int64

// Extension function for UserID conversion
func userIDExtension(value any) (func() error, bool) {
    str, ok := value.(string)
    if !ok || !strings.HasPrefix(str, "user:") {
        return nil, false
    }

    return func() error {
        idStr := strings.TrimPrefix(str, "user:")
        id, err := strconv.ParseInt(idStr, 10, 64)
        if err != nil {
            return err
        }
        *userID = UserID(id)
        return nil
    }, true
}

// Usage
var userID UserID
field := reflect.ValueOf(&userID).Elem()
assign.Value(field, "user:12345", userIDExtension)
```

### Interface Support

Built-in support for standard Go interfaces:

```go
// Custom type implementing TextUnmarshaler
type Email string

func (e *Email) UnmarshalText(text []byte) error {
    s := string(text)
    if !strings.Contains(s, "@") {
        return errors.New("invalid email format")
    }
    *e = Email(s)
    return nil
}

// Automatic interface detection
var email Email
field := reflect.ValueOf(&email).Elem()
assign.String(field, "user@example.com") // Uses UnmarshalText automatically
```

### Time Parsing

Intelligent time parsing with multiple format support:

```go
var timestamp time.Time
field := reflect.ValueOf(&timestamp).Elem()

// Supports multiple formats automatically
assign.String(field, "2023-01-15T10:30:00Z")       // RFC3339
assign.String(field, "2023-01-15")                 // Date only
assign.String(field, "01/15/2023")                 // US format
assign.String(field, "2023-01-15 10:30:00")        // DateTime
```

## Type Constraints

The library uses Go 1.18+ generics with type constraints for compile-time safety:

```go
// Integer constraints
type SignedValue interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedValue interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type IntegerValue interface {
    SignedValue | UnsignedValue
}

// Float constraints
type FloatValue interface {
    ~float32 | ~float64
}
```

## Error Handling

The library provides detailed error information for debugging:

```go
var target int8
field := reflect.ValueOf(&target).Elem()

err := assign.Integer(field, 300) // Value too large for int8
if err != nil {
    fmt.Println(err) // "value 300 is outside the range of target type int8 [-128, 127]"
}

err = assign.String(field, "not_a_number")
if errors.Is(err, assign.ErrNotSupported) {
    fmt.Println("Conversion not supported")
}
```

## Performance

The library is optimized for performance with:

- **Fast paths** for common type conversions
- **Zero allocations** for primitive type conversions
- **Compiled regex patterns** for time format detection
- **Thread-safe caching** for successful time format matches
- **Optimized string operations** with minimal allocations

## Real-World Example

Here's a practical example of using the library for configuration parsing:

```go
package main

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"

    "github.com/slipros/assign"
)

type Config struct {
    Host     string        `env:"HOST"`
    Port     int           `env:"PORT"`
    Debug    bool          `env:"DEBUG"`
    Timeout  time.Duration `env:"TIMEOUT"`
    Features []string      `env:"FEATURES"`
}

func parseConfig(envVars map[string]string) (*Config, error) {
    config := &Config{}
    v := reflect.ValueOf(config).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        envTag := fieldType.Tag.Get("env")
        if envTag == "" {
            continue
        }

        envValue, exists := envVars[envTag]
        if !exists {
            continue
        }

        // Handle different field types
        switch fieldType.Type.String() {
        case "[]string":
            parts := strings.Split(envValue, ",")
            if err := assign.SliceString(field, parts); err != nil {
                return nil, err
            }
        case "time.Duration":
            // Custom extension for duration
            duration, err := time.ParseDuration(envValue)
            if err != nil {
                return nil, err
            }
            field.Set(reflect.ValueOf(duration))
        default:
            if err := assign.String(field, envValue); err != nil {
                return nil, err
            }
        }
    }

    return config, nil
}

func main() {
    envVars := map[string]string{
        "HOST":     "localhost",
        "PORT":     "8080",
        "DEBUG":    "true",
        "TIMEOUT":  "30s",
        "FEATURES": "auth,logging,metrics",
    }

    config, err := parseConfig(envVars)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Config: %+v\n", config)
}
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/slipros/assign)
- [Examples](https://pkg.go.dev/github.com/slipros/assign#pkg-examples)
- [GitHub Pages Documentation](https://slipros.github.io/assign)

---

Made with ❤️ for the Go community
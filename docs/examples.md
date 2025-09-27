---
layout: default
title: Examples
---

# Examples

This page provides comprehensive examples of using the assign library for various type conversion scenarios.

## Table of Contents

- [Basic Usage](#basic-usage)
- [String Conversions](#string-conversions)
- [Numeric Conversions](#numeric-conversions)
- [Boolean Conversions](#boolean-conversions)
- [Time Parsing](#time-parsing)
- [Slice Operations](#slice-operations)
- [Nil Pointer Handling](#nil-pointer-handling)
- [Extension Functions](#extension-functions)
- [Error Handling](#error-handling)
- [Real-World Examples](#real-world-examples)

## Basic Usage

### Simple Type Conversion

```go
package main

import (
    "fmt"
    "reflect"
    "github.com/slipros/assign"
)

func main() {
    var target struct {
        ID   int
        Name string
        Age  int
    }

    v := reflect.ValueOf(&target).Elem()

    // Convert string to int
    if err := assign.Value(v.FieldByName("ID"), "123"); err != nil {
        panic(err)
    }

    // Set string directly
    if err := assign.Value(v.FieldByName("Name"), "John Doe"); err != nil {
        panic(err)
    }

    // Convert string to int
    if err := assign.Value(v.FieldByName("Age"), "30"); err != nil {
        panic(err)
    }

    fmt.Printf("Target: %+v\n", target)
    // Output: Target: {ID:123 Name:John Doe Age:30}
}
```

### Universal Value Assignment

```go
func ExampleUniversalAssignment() {
    var data struct {
        StringField  string
        IntField     int
        FloatField   float64
        BoolField    bool
        SliceField   []string
    }

    v := reflect.ValueOf(&data).Elem()

    // Various source types
    if err := assign.Value(v.FieldByName("StringField"), "hello"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("IntField"), "42"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("FloatField"), "3.14"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("BoolField"), "true"); err != nil {
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("SliceField"), []string{"a", "b", "c"}); err != nil {
        panic(err)
    }

    fmt.Printf("Data: %+v\n", data)
}
```

## String Conversions

### String to Various Types

```go
func ExampleStringConversions() {
    var targets struct {
        Number     int
        Float      float64
        Bool       bool
        Complex    complex128
        Bytes      []byte
        Items      []string
    }

    v := reflect.ValueOf(&targets).Elem()

    // String to different types
    if err := assign.String(v.FieldByName("Number"), "42"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Float"), "3.14159"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Bool"), "true"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Complex"), "3+4i"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Bytes"), "hello world"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Items"), "apple,banana,cherry"); err != nil {
        panic(err)
    }

    fmt.Printf("Number: %d\n", targets.Number)
    fmt.Printf("Float: %.3f\n", targets.Float)
    fmt.Printf("Bool: %t\n", targets.Bool)
    fmt.Printf("Complex: %v\n", targets.Complex)
    fmt.Printf("Bytes: %s\n", string(targets.Bytes))
    fmt.Printf("Items: %v\n", targets.Items)
}
```

### Number Base Detection

```go
func ExampleNumberBases() {
    var targets struct {
        Decimal     int
        Hexadecimal int
        Octal       int
        Binary      int
    }

    v := reflect.ValueOf(&targets).Elem()

    // Different number bases
    if err := assign.String(v.FieldByName("Decimal"), "123"); err != nil {      // Base 10
        panic(err)
    }
    if err := assign.String(v.FieldByName("Hexadecimal"), "0x1A"); err != nil { // Base 16
        panic(err)
    }
    if err := assign.String(v.FieldByName("Octal"), "0123"); err != nil {       // Base 8
        panic(err)
    }
    if err := assign.String(v.FieldByName("Binary"), "0b1010"); err != nil {    // Base 2
        panic(err)
    }

    fmt.Printf("Decimal: %d\n", targets.Decimal)         // 123
    fmt.Printf("Hexadecimal: %d\n", targets.Hexadecimal) // 26
    fmt.Printf("Octal: %d\n", targets.Octal)             // 83
    fmt.Printf("Binary: %d\n", targets.Binary)           // 10
}
```

## Numeric Conversions

### Integer Type Conversion with Safety

```go
func ExampleIntegerConversions() {
    var targets struct {
        Small   int8
        Medium  int16
        Large   int64
        Unsigned uint32
    }

    v := reflect.ValueOf(&targets).Elem()

    // Safe integer conversions
    if err := assign.Integer(v.FieldByName("Small"), 42); err != nil {
        panic(err)
    }
    if err := assign.Integer(v.FieldByName("Medium"), int32(1000)); err != nil {
        panic(err)
    }
    if err := assign.Integer(v.FieldByName("Large"), int64(999999999)); err != nil {
        panic(err)
    }
    if err := assign.Integer(v.FieldByName("Unsigned"), uint64(123456)); err != nil {
        panic(err)
    }

    fmt.Printf("Small: %d\n", targets.Small)
    fmt.Printf("Medium: %d\n", targets.Medium)
    fmt.Printf("Large: %d\n", targets.Large)
    fmt.Printf("Unsigned: %d\n", targets.Unsigned)

    // This would cause an error (overflow):
    // assign.Integer(v.FieldByName("Small"), 300) // Error: value 300 > int8 max (127)
}
```

### Float Conversions with Special Values

```go
func ExampleFloatConversions() {
    var targets struct {
        Regular    float64
        ToInt      int
        ToString   string
        SpecialNaN float64
        SpecialInf float64
    }

    v := reflect.ValueOf(&targets).Elem()

    // Regular float conversions
    if err := assign.Float(v.FieldByName("Regular"), 3.14159); err != nil {
        panic(err)
    }
    if err := assign.Float(v.FieldByName("ToInt"), 42.7); err != nil {      // Truncated to 42
        panic(err)
    }
    if err := assign.Float(v.FieldByName("ToString"), 123.456); err != nil {
        panic(err)
    }

    // Special values
    if err := assign.Float(v.FieldByName("SpecialNaN"), math.NaN()); err != nil {
        panic(err)
    }
    if err := assign.Float(v.FieldByName("SpecialInf"), math.Inf(1)); err != nil {
        panic(err)
    }

    fmt.Printf("Regular: %.3f\n", targets.Regular)
    fmt.Printf("ToInt: %d\n", targets.ToInt)
    fmt.Printf("ToString: %s\n", targets.ToString)
    fmt.Printf("NaN: %v\n", math.IsNaN(targets.SpecialNaN))
    fmt.Printf("Inf: %v\n", math.IsInf(targets.SpecialInf, 1))
}
```

## Boolean Conversions

### Multiple Boolean Formats

```go
func ExampleBooleanFormats() {
    var results struct {
        A, B, C, D, E, F, G, H, I, J bool
    }

    v := reflect.ValueOf(&results).Elem()
    fields := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

    // Various boolean representations
    values := []string{
        "true", "1", "yes", "on", "t",    // True values
        "false", "0", "no", "off", "f",   // False values
    }

    for i, value := range values {
        if err := assign.String(v.FieldByName(fields[i]), value); err != nil {
            panic(err)
        }
    }

    fmt.Printf("Boolean conversions:\n")
    fmt.Printf("'true' -> %t\n", results.A)
    fmt.Printf("'1' -> %t\n", results.B)
    fmt.Printf("'yes' -> %t\n", results.C)
    fmt.Printf("'on' -> %t\n", results.D)
    fmt.Printf("'t' -> %t\n", results.E)
    fmt.Printf("'false' -> %t\n", results.F)
    fmt.Printf("'0' -> %t\n", results.G)
    fmt.Printf("'no' -> %t\n", results.H)
    fmt.Printf("'off' -> %t\n", results.I)
    fmt.Printf("'f' -> %t\n", results.J)
}
```

## Time Parsing

### Multiple Time Formats

```go
func ExampleTimeFormats() {
    var timestamps struct {
        RFC3339    time.Time
        RFC3339Nano time.Time
        DateOnly   time.Time
        DateTime   time.Time
        USFormat   time.Time
        Custom     time.Time
    }

    v := reflect.ValueOf(&timestamps).Elem()

    // Various time formats (automatically detected)
    if err := assign.String(v.FieldByName("RFC3339"), "2023-01-15T10:30:00Z"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("RFC3339Nano"), "2023-01-15T10:30:00.123456789Z"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("DateOnly"), "2023-01-15"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("DateTime"), "2023-01-15 10:30:00"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("USFormat"), "01/15/2023"); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Custom"), "15 Jan 2023 10:30:00 GMT"); err != nil {
        panic(err)
    }

    fmt.Printf("RFC3339: %s\n", timestamps.RFC3339.Format("2006-01-02 15:04:05"))
    fmt.Printf("RFC3339Nano: %s\n", timestamps.RFC3339Nano.Format("2006-01-02 15:04:05.999999999"))
    fmt.Printf("DateOnly: %s\n", timestamps.DateOnly.Format("2006-01-02"))
    fmt.Printf("DateTime: %s\n", timestamps.DateTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("USFormat: %s\n", timestamps.USFormat.Format("2006-01-02"))
    fmt.Printf("Custom: %s\n", timestamps.Custom.Format("2006-01-02 15:04:05"))
}
```

## Slice Operations

### String Slice Conversions

```go
func ExampleSliceConversions() {
    var targets struct {
        IntSlice    []int
        FloatSlice  []float64
        BoolSlice   []bool
        JoinedComma string
        JoinedPipe  string
        JoinedSpace string
    }

    v := reflect.ValueOf(&targets).Elem()

    // Convert to typed slices
    if err := assign.SliceString(v.FieldByName("IntSlice"), []string{"1", "2", "3", "4", "5"}); err != nil {
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("FloatSlice"), []string{"1.1", "2.2", "3.3"}); err != nil {
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("BoolSlice"), []string{"true", "false", "1", "0"}); err != nil {
        panic(err)
    }

    // Join with different separators
    fruits := []string{"apple", "banana", "cherry"}
    if err := assign.SliceString(v.FieldByName("JoinedComma"), fruits); err != nil { // Default comma
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("JoinedPipe"), fruits, assign.WithSeparator(" | ")); err != nil {
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("JoinedSpace"), fruits, assign.WithSeparator(" ")); err != nil {
        panic(err)
    }

    fmt.Printf("IntSlice: %v\n", targets.IntSlice)
    fmt.Printf("FloatSlice: %v\n", targets.FloatSlice)
    fmt.Printf("BoolSlice: %v\n", targets.BoolSlice)
    fmt.Printf("JoinedComma: %s\n", targets.JoinedComma)
    fmt.Printf("JoinedPipe: %s\n", targets.JoinedPipe)
    fmt.Printf("JoinedSpace: %s\n", targets.JoinedSpace)
}
```

### Nested Slice Conversion

```go
func ExampleNestedSliceConversion() {
    var data struct {
        Matrix [][]int
        Items  []any
    }

    v := reflect.ValueOf(&data).Elem()

    // Convert slice of any to slice of slice of int
    sourceData := []any{
        []string{"1", "2", "3"},
        []string{"4", "5", "6"},
        []string{"7", "8", "9"},
    }

    if err := assign.Value(v.FieldByName("Matrix"), sourceData); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("Items"), []string{"hello", "42", "true", "3.14"}); err != nil {
        panic(err)
    }

    fmt.Printf("Matrix: %v\n", data.Matrix)
    fmt.Printf("Items: %v\n", data.Items)
}
```

## Nil Pointer Handling

### Automatic Pointer Initialization

```go
func ExampleNilPointers() {
    var data struct {
        Name     *string
        Age      *int
        Score    *float64
        Active   *bool
        Tags     *[]string
        Metadata *map[string]string
    }

    v := reflect.ValueOf(&data).Elem()

    // All fields start as nil, but assign will initialize them
    if err := assign.Value(v.FieldByName("Name"), "John Doe"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("Age"), "30"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("Score"), "95.5"); err != nil {
        panic(err)
    }
    if err := assign.Value(v.FieldByName("Active"), "true"); err != nil {
        panic(err)
    }
    if err := assign.SliceString(v.FieldByName("Tags"), []string{"developer", "golang", "backend"}); err != nil {
        panic(err)
    }
    if err := assign.String(v.FieldByName("Metadata"), "key1:value1,key2:value2"); err != nil {
        panic(err)
    }

    fmt.Printf("Name: %s\n", *data.Name)
    fmt.Printf("Age: %d\n", *data.Age)
    fmt.Printf("Score: %.1f\n", *data.Score)
    fmt.Printf("Active: %t\n", *data.Active)
    fmt.Printf("Tags: %v\n", *data.Tags)
    fmt.Printf("Metadata: %v\n", *data.Metadata)
}
```

### Deep Pointer Nesting

```go
func ExampleDeepPointers() {
    var data struct {
        Level1 *struct {
            Level2 *struct {
                Value *string
            }
        }
    }

    v := reflect.ValueOf(&data).Elem()

    // Automatically initializes all nested nil pointers
    field := v.FieldByName("Level1").FieldByName("Level2").FieldByName("Value")
    if err := assign.Value(field, "deep value"); err != nil {
        panic(err)
    }

    fmt.Printf("Deep value: %s\n", ***data.Level1.Level2.Value)
}
```

## Extension Functions

### HTTP Cookie Extension

```go
import (
    "fmt"
    "net/http"
    "reflect"
    "time"
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

// Extension for converting Cookie to custom types
func cookieToIntExtension(value any) (func(to reflect.Value) error, bool) {
    cookie, ok := value.(*http.Cookie)
    if !ok || cookie.Name != "user_id" {
        return nil, false
    }

    return func(to reflect.Value) error {
        return assign.String(to, cookie.Value)
    }, true
}

func ExampleExtensions() {
    // Create sample cookies
    sessionCookie := &http.Cookie{Name: "session", Value: "abc123def456"}
    userIDCookie := &http.Cookie{Name: "user_id", Value: "42"}

    var sessionID string
    var userID int

    sessionField := reflect.ValueOf(&sessionID).Elem()
    userField := reflect.ValueOf(&userID).Elem()

    // Use extensions with assign.Value
    if err := assign.Value(sessionField, sessionCookie, cookieExtension); err != nil {
        panic(err)
    }

    if err := assign.Value(userField, userIDCookie, cookieToIntExtension); err != nil {
        panic(err)
    }

    fmt.Printf("Session ID: %s\n", sessionID)
    fmt.Printf("User ID: %d\n", userID)

    // Output:
    // Session ID: abc123def456
    // User ID: 42
}
```

### Multiple Extensions

```go
func ExampleMultipleExtensions() {
    // Multiple extensions can be chained
    var sessionID string
    field := reflect.ValueOf(&sessionID).Elem()

    extensions := []assign.ExtensionFunc{
        cookieExtension,
        cookieToIntExtension,
        // Add more custom extensions as needed
    }

    // The first matching extension will be used
    cookie := &http.Cookie{Name: "session", Value: "xyz789"}
    if err := assign.Value(field, cookie, extensions...); err != nil {
        panic(err)
    }

    fmt.Printf("Session from cookie: %s\n", sessionID)
    // Output: Session from cookie: xyz789
}
```

## Error Handling

### Comprehensive Error Handling

```go
func ExampleErrorHandling() {
    var targets struct {
        ValidInt    int8
        InvalidInt  int8
        ValidFloat  float32
        InvalidTime time.Time
    }

    v := reflect.ValueOf(&targets).Elem()

    // Successful conversion
    if err := assign.Integer(v.FieldByName("ValidInt"), 42); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("ValidInt: %d\n", targets.ValidInt)
    }

    // Overflow error
    if err := assign.Integer(v.FieldByName("InvalidInt"), 300); err != nil {
        fmt.Printf("Overflow error: %v\n", err)
    }

    // Successful float conversion
    if err := assign.Float(v.FieldByName("ValidFloat"), 3.14); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("ValidFloat: %.2f\n", targets.ValidFloat)
    }

    // Invalid time format
    if err := assign.String(v.FieldByName("InvalidTime"), "not a time"); err != nil {
        fmt.Printf("Time parse error: %v\n", err)

        // Check for specific error types
        if errors.Is(err, assign.ErrNotSupported) {
            fmt.Println("This conversion is not supported")
        }
    }
}
```

### Error Type Checking

```go
func ExampleErrorTypes() {
    var target int
    field := reflect.ValueOf(&target).Elem()

    err := assign.Value(field, complex(1, 2)) // Complex to int not supported

    switch {
    case err == nil:
        fmt.Println("Conversion successful")
    case errors.Is(err, assign.ErrNotSupported):
        fmt.Printf("Unsupported conversion: %v\n", err)
    default:
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Real-World Examples

### Configuration Parsing

```go
type AppConfig struct {
    Host         string        `env:"HOST"`
    Port         int           `env:"PORT"`
    Debug        bool          `env:"DEBUG"`
    Timeout      time.Duration `env:"TIMEOUT"`
    DatabaseURL  string        `env:"DATABASE_URL"`
    Features     []string      `env:"FEATURES"`
    Limits       map[string]int `env:"LIMITS"`
}

func parseConfigFromEnv(envVars map[string]string) (*AppConfig, error) {
    config := &AppConfig{}
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
                return nil, fmt.Errorf("failed to parse %s: %w", envTag, err)
            }
        case "time.Duration":
            duration, err := time.ParseDuration(envValue)
            if err != nil {
                return nil, fmt.Errorf("failed to parse duration %s: %w", envTag, err)
            }
            field.Set(reflect.ValueOf(duration))
        case "map[string]int":
            if err := assign.String(field, envValue); err != nil {
                return nil, fmt.Errorf("failed to parse map %s: %w", envTag, err)
            }
        default:
            if err := assign.String(field, envValue); err != nil {
                return nil, fmt.Errorf("failed to parse %s: %w", envTag, err)
            }
        }
    }

    return config, nil
}

func ExampleConfigParsing() {
    envVars := map[string]string{
        "HOST":         "localhost",
        "PORT":         "8080",
        "DEBUG":        "true",
        "TIMEOUT":      "30s",
        "DATABASE_URL": "postgres://localhost/mydb",
        "FEATURES":     "auth,logging,metrics",
        "LIMITS":       "requests:1000,connections:100",
    }

    config, err := parseConfigFromEnv(envVars)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Config: %+v\n", config)
}
```

### JSON-like Data Mapping

```go
type User struct {
    ID       int64     `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Age      int       `json:"age"`
    IsActive bool      `json:"is_active"`
    JoinedAt time.Time `json:"joined_at"`
    Tags     []string  `json:"tags"`
    Metadata map[string]string `json:"metadata"`
}

func mapDataToStruct(data map[string]any, target any) error {
    v := reflect.ValueOf(target).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" || jsonTag == "-" {
            continue
        }

        value, exists := data[jsonTag]
        if !exists {
            continue
        }

        switch v := value.(type) {
        case []string:
            if err := assign.SliceString(field, v); err != nil {
                return fmt.Errorf("failed to assign slice %s: %w", jsonTag, err)
            }
        case map[string]string:
            // Convert map to string representation for parsing
            var parts []string
            for k, val := range v {
                parts = append(parts, fmt.Sprintf("%s:%s", k, val))
            }
            if err := assign.String(field, strings.Join(parts, ",")); err != nil {
                return fmt.Errorf("failed to assign map %s: %w", jsonTag, err)
            }
        default:
            if err := assign.Value(field, value); err != nil {
                return fmt.Errorf("failed to assign %s: %w", jsonTag, err)
            }
        }
    }

    return nil
}

func ExampleDataMapping() {
    data := map[string]any{
        "id":         "12345",
        "name":       "John Doe",
        "email":      "john@example.com",
        "age":        "30",
        "is_active":  "true",
        "joined_at":  "2023-01-15T10:30:00Z",
        "tags":       []string{"developer", "golang", "backend"},
        "metadata":   map[string]string{"department": "engineering", "level": "senior"},
    }

    var user User
    if err := mapDataToStruct(data, &user); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("User: %+v\n", user)
}
```

### Form Data Processing

```go
type ContactForm struct {
    Name      string    `form:"name"`
    Email     string    `form:"email"`
    Age       int       `form:"age"`
    Subscribe bool      `form:"subscribe"`
    Interests []string  `form:"interests"`
    Message   string    `form:"message"`
    SubmitTime time.Time
}

func processFormData(formData map[string][]string) (*ContactForm, error) {
    form := &ContactForm{
        SubmitTime: time.Now(),
    }
    v := reflect.ValueOf(form).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        formTag := fieldType.Tag.Get("form")
        if formTag == "" {
            continue
        }

        values, exists := formData[formTag]
        if !exists || len(values) == 0 {
            continue
        }

        // Handle multi-value fields (like checkboxes)
        if len(values) > 1 && field.Kind() == reflect.Slice {
            if err := assign.SliceString(field, values); err != nil {
                return nil, fmt.Errorf("failed to assign %s: %w", formTag, err)
            }
        } else {
            // Single value
            if err := assign.String(field, values[0]); err != nil {
                return nil, fmt.Errorf("failed to assign %s: %w", formTag, err)
            }
        }
    }

    return form, nil
}

func ExampleFormProcessing() {
    // Simulate form data (like from http.Request.Form)
    formData := map[string][]string{
        "name":      {"John Doe"},
        "email":     {"john@example.com"},
        "age":       {"30"},
        "subscribe": {"true"},
        "interests": {"golang", "web-development", "databases"},
        "message":   {"Hello, I'm interested in your services."},
    }

    form, err := processFormData(formData)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Form: %+v\n", form)
}
```

These examples demonstrate the flexibility and power of the assign library for various real-world scenarios involving type conversion and data mapping.
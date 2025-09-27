package assign_test

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/slipros/assign"
)

// ExampleValue demonstrates basic usage of the Value function
// for converting between different types.
func ExampleValue() {
	var target struct {
		ID       int
		Name     string
		Age      *int
		IsActive bool
		Score    float64
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

	// Initialize nil pointer and set value
	if err := assign.Value(v.FieldByName("Age"), "30"); err != nil {
		panic(err)
	}

	// Convert string to bool
	if err := assign.Value(v.FieldByName("IsActive"), "true"); err != nil {
		panic(err)
	}

	// Convert string to float
	if err := assign.Value(v.FieldByName("Score"), "95.5"); err != nil {
		panic(err)
	}

	fmt.Printf("ID: %d\n", target.ID)
	fmt.Printf("Name: %s\n", target.Name)
	fmt.Printf("Age: %d\n", *target.Age)
	fmt.Printf("IsActive: %t\n", target.IsActive)
	fmt.Printf("Score: %.1f\n", target.Score)

	// Output:
	// ID: 123
	// Name: John Doe
	// Age: 30
	// IsActive: true
	// Score: 95.5
}

// ExampleInteger demonstrates type-safe integer conversion with
// overflow protection and range checking.
func ExampleInteger() {
	var targets struct {
		Small  int8
		Medium int16
		Large  int64
	}

	v := reflect.ValueOf(&targets).Elem()

	// Safe conversions within range
	if err := assign.Integer(v.FieldByName("Small"), int(42)); err != nil {
		panic(err)
	}
	if err := assign.Integer(v.FieldByName("Medium"), int32(1000)); err != nil {
		panic(err)
	}
	if err := assign.Integer(v.FieldByName("Large"), int64(999999999)); err != nil {
		panic(err)
	}

	fmt.Printf("Small: %d\n", targets.Small)
	fmt.Printf("Medium: %d\n", targets.Medium)
	fmt.Printf("Large: %d\n", targets.Large)

	// Output:
	// Small: 42
	// Medium: 1000
	// Large: 999999999
}

// ExampleFloat demonstrates floating-point conversion with
// special value handling.
func ExampleFloat() {
	var targets struct {
		Regular float32
		Double  float64
		Text    string
	}

	v := reflect.ValueOf(&targets).Elem()

	// Float conversions
	if err := assign.Float(v.FieldByName("Regular"), 3.14159); err != nil {
		panic(err)
	}
	if err := assign.Float(v.FieldByName("Double"), 2.71828182845); err != nil {
		panic(err)
	}

	// Float to string conversion
	if err := assign.Float(v.FieldByName("Text"), 123.456); err != nil {
		panic(err)
	}

	fmt.Printf("Regular: %.2f\n", targets.Regular)
	fmt.Printf("Double: %.5f\n", targets.Double)
	fmt.Printf("Text: %s\n", targets.Text)

	// Output:
	// Regular: 3.14
	// Double: 2.71828
	// Text: 123.456
}

// ExampleString demonstrates string parsing and conversion to various types
// including optimized fast paths for common conversions.
func ExampleString() {
	var targets struct {
		Number  int
		Active  bool
		Score   float64
		Items   []string
		Bytes   []byte
		Complex complex128
	}

	v := reflect.ValueOf(&targets).Elem()

	// String to various types
	if err := assign.String(v.FieldByName("Number"), "42"); err != nil {
		panic(err)
	}
	if err := assign.String(v.FieldByName("Active"), "true"); err != nil {
		panic(err)
	}
	if err := assign.String(v.FieldByName("Score"), "98.6"); err != nil {
		panic(err)
	}
	if err := assign.String(v.FieldByName("Items"), "apple,banana,cherry"); err != nil {
		panic(err)
	}
	if err := assign.String(v.FieldByName("Bytes"), "hello"); err != nil {
		panic(err)
	}
	if err := assign.String(v.FieldByName("Complex"), "3+4i"); err != nil {
		panic(err)
	}

	fmt.Printf("Number: %d\n", targets.Number)
	fmt.Printf("Active: %t\n", targets.Active)
	fmt.Printf("Score: %.1f\n", targets.Score)
	fmt.Printf("Items: %v\n", targets.Items)
	fmt.Printf("Bytes: %s\n", string(targets.Bytes))
	fmt.Printf("Complex: %v\n", targets.Complex)

	// Output:
	// Number: 42
	// Active: true
	// Score: 98.6
	// Items: [apple banana cherry]
	// Bytes: hello
	// Complex: (3+4i)
}

// ExampleSliceString demonstrates converting string slices to various types
// and joining them with custom separators.
func ExampleSliceString() {
	var targets struct {
		Numbers []int
		Joined  string
		Custom  string
	}

	v := reflect.ValueOf(&targets).Elem()
	stringSlice := []string{"1", "2", "3", "4", "5"}

	// Convert to int slice
	if err := assign.SliceString(v.FieldByName("Numbers"), stringSlice); err != nil {
		panic(err)
	}

	// Join with default separator (comma)
	if err := assign.SliceString(v.FieldByName("Joined"), []string{"apple", "banana", "cherry"}); err != nil {
		panic(err)
	}

	// Join with custom separator
	if err := assign.SliceString(v.FieldByName("Custom"), []string{"red", "green", "blue"}, assign.WithSeparator(" | ")); err != nil {
		panic(err)
	}

	fmt.Printf("Numbers: %v\n", targets.Numbers)
	fmt.Printf("Joined: %s\n", targets.Joined)
	fmt.Printf("Custom: %s\n", targets.Custom)

	// Output:
	// Numbers: [1 2 3 4 5]
	// Joined: apple,banana,cherry
	// Custom: red | green | blue
}

// ExampleValue_extensionFunc demonstrates how to create and use
// extension functions for custom type conversions.
func ExampleValue_extensionFunc() {
	// Extension function for HTTP Cookie conversion
	cookieExtension := func(value any) (func(to reflect.Value) error, bool) {
		cookie, ok := value.(*http.Cookie)
		if !ok {
			return nil, false
		}

		return func(to reflect.Value) error {
			return assign.String(to, cookie.Value)
		}, true
	}

	// Create a sample cookie
	cookie := &http.Cookie{Name: "session", Value: "abc123def456"}
	var sessionID string
	field := reflect.ValueOf(&sessionID).Elem()

	// Use the extension function with assign.Value
	if err := assign.Value(field, cookie, cookieExtension); err != nil {
		panic(err)
	}

	fmt.Printf("Session ID: %s\n", sessionID)

	// Output:
	// Session ID: abc123def456
}

// ExampleString_timeConversion demonstrates automatic time parsing
// with multiple format support.
func ExampleString_timeConversion() {
	var timestamps struct {
		RFC3339  time.Time
		DateOnly time.Time
		DateTime time.Time
		USFormat time.Time
	}

	v := reflect.ValueOf(&timestamps).Elem()

	// Various time formats
	if err := assign.String(v.FieldByName("RFC3339"), "2023-01-15T10:30:00Z"); err != nil {
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

	fmt.Printf("RFC3339: %s\n", timestamps.RFC3339.Format("2006-01-02 15:04:05"))
	fmt.Printf("DateOnly: %s\n", timestamps.DateOnly.Format("2006-01-02"))
	fmt.Printf("DateTime: %s\n", timestamps.DateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("USFormat: %s\n", timestamps.USFormat.Format("2006-01-02"))

	// Output:
	// RFC3339: 2023-01-15 10:30:00
	// DateOnly: 2023-01-15
	// DateTime: 2023-01-15 10:30:00
	// USFormat: 2023-01-15
}

// ExampleValue_nilPointers demonstrates automatic nil pointer
// initialization during value assignment.
func ExampleValue_nilPointers() {
	var data struct {
		Name  *string
		Age   *int
		Score *float64
		Tags  *[]string
	}

	v := reflect.ValueOf(&data).Elem()

	// All fields start as nil, but assign will initialize them
	if err := assign.Value(v.FieldByName("Name"), "John"); err != nil {
		panic(err)
	}
	if err := assign.Value(v.FieldByName("Age"), "30"); err != nil {
		panic(err)
	}
	if err := assign.Value(v.FieldByName("Score"), "95.5"); err != nil {
		panic(err)
	}
	if err := assign.SliceString(v.FieldByName("Tags"), []string{"developer", "golang"}); err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", *data.Name)
	fmt.Printf("Age: %d\n", *data.Age)
	fmt.Printf("Score: %.1f\n", *data.Score)
	fmt.Printf("Tags: %v\n", *data.Tags)

	// Output:
	// Name: John
	// Age: 30
	// Score: 95.5
	// Tags: [developer golang]
}

// ExampleString_booleanConversions demonstrates various ways
// to convert strings to boolean values.
func ExampleString_booleanConversions() {
	var results struct {
		A, B, C, D, E, F, G, H bool
	}

	v := reflect.ValueOf(&results).Elem()
	fields := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	values := []string{"true", "1", "yes", "on", "false", "0", "no", "off"}

	for i, value := range values {
		if err := assign.String(v.FieldByName(fields[i]), value); err != nil {
			panic(err)
		}
	}

	fmt.Printf("'true' -> %t\n", results.A)
	fmt.Printf("'1' -> %t\n", results.B)
	fmt.Printf("'yes' -> %t\n", results.C)
	fmt.Printf("'on' -> %t\n", results.D)
	fmt.Printf("'false' -> %t\n", results.E)
	fmt.Printf("'0' -> %t\n", results.F)
	fmt.Printf("'no' -> %t\n", results.G)
	fmt.Printf("'off' -> %t\n", results.H)

	// Output:
	// 'true' -> true
	// '1' -> true
	// 'yes' -> true
	// 'on' -> true
	// 'false' -> false
	// '0' -> false
	// 'no' -> false
	// 'off' -> false
}

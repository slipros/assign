package assign

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTryUnmarshaleersSpecificCoverage tests specific uncovered lines in tryUnmarshalers
func TestTryUnmarshaleersSpecificCoverage(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (reflect.Value, string)
		expectError bool
	}{
		{
			name: "field that cannot get interface",
			setup: func() (reflect.Value, string) {
				// This is hard to create in practice, but let's try with an invalid field
				var target struct{ unexported string }
				field := reflect.ValueOf(target).FieldByName("unexported")
				return field, "test"
			},
			expectError: true, // Should error because field is not addressable
		},
		{
			name: "type that doesn't implement unmarshalers",
			setup: func() (reflect.Value, string) {
				// Test with a custom struct which doesn't implement TextUnmarshaler or BinaryUnmarshaler
				var target struct{ Name string }
				field := reflect.ValueOf(&target).Elem()
				return field, "test data"
			},
			expectError: true, // Should hit the "does not implement" error case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, str := tt.setup()
			err := String(field, str)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetTimeFormatsByLengthCoverage tests missing coverage in getTimeFormatsByLength
func TestGetTimeFormatsByLengthCoverage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "length 6 format",
			input:       "123456", // 6 characters - no valid format
			expectError: true,
		},
		{
			name:        "length 7 format",
			input:       "1234567", // 7 characters - no valid format
			expectError: true,
		},
		{
			name:        "length 11 format",
			input:       "12345678901", // 11 characters - should trigger case 11
			expectError: true,          // No valid format for this case
		},
		{
			name:        "length 16 format",
			input:       "1234567890123456", // 16 characters - should trigger case 16
			expectError: true,               // No valid format for this case
		},
		{
			name:        "length 15 format",
			input:       "123456789012345", // 15 characters - should trigger case 15
			expectError: true,              // No valid format for this case
		},
		{
			name:        "length 9 format",
			input:       "123456789", // 9 characters - should trigger case 9
			expectError: true,        // No valid format for this case
		},
		{
			name:        "length 20 format",
			input:       "12345678901234567890", // 20 characters - should trigger case 20 (RFC3339)
			expectError: true,                   // Invalid format but triggers the case
		},
		{
			name:        "length 25 format",
			input:       "1234567890123456789012345", // 25 characters - should trigger case 25 (RFC3339)
			expectError: true,                        // Invalid format but triggers the case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target struct{ TimeField time.Time }
			field := reflect.ValueOf(&target).Elem().FieldByName("TimeField")
			err := String(field, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEqualFoldSpecificCoverage tests uncovered branches in equalFold function
func TestEqualFoldSpecificCoverage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "mixed case TRUE",
			input:    "TrUe",
			expected: true,
		},
		{
			name:     "mixed case FALSE",
			input:    "FaLsE",
			expected: true,
		},
		{
			name:     "mixed case YES",
			input:    "YeS",
			expected: true,
		},
		{
			name:     "mixed case NO",
			input:    "nO",
			expected: true,
		},
		{
			name:     "mixed case ON",
			input:    "On",
			expected: true,
		},
		{
			name:     "mixed case OFF",
			input:    "OfF",
			expected: true,
		},
		{
			name:     "single character Y - short string path",
			input:    "Y",
			expected: true,
		},
		{
			name:     "single character N - short string path",
			input:    "N",
			expected: true,
		},
		{
			name:     "empty string - should trigger short string path",
			input:    "",
			expected: true, // Empty string gets handled by handleEmptyString and sets to false (no error)
		},
		{
			name:     "case mismatch requiring character-by-character comparison",
			input:    "TruE", // Mixed case to trigger character comparison path
			expected: true,
		},
		{
			name:     "case where characters are the same",
			input:    "true", // Should succeed with no case conversion needed
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target bool
			field := reflect.ValueOf(&target).Elem()
			err := String(field, tt.input)

			if tt.expected {
				// Should succeed (no error)
				assert.NoError(t, err)
			} else {
				// Should fail
				assert.Error(t, err)
			}
		})
	}
}

// TestFloatSpecificCoverage tests missing coverage in Float function
func TestFloatSpecificCoverage(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (reflect.Value, any)
		expectError bool
	}{
		{
			name: "very small float32 to complex64",
			setup: func() (reflect.Value, any) {
				var target complex64
				return reflect.ValueOf(&target).Elem(), float32(1e-10)
			},
			expectError: false,
		},
		{
			name: "very large float64 to complex128",
			setup: func() (reflect.Value, any) {
				var target complex128
				return reflect.ValueOf(&target).Elem(), float64(1e100)
			},
			expectError: false,
		},
		{
			name: "small positive float to uint - truncation warning path",
			setup: func() (reflect.Value, any) {
				var target uint32
				return reflect.ValueOf(&target).Elem(), float64(0.5) // Will be truncated to 0
			},
			expectError: false,
		},
		{
			name: "small positive float to int - truncation warning path",
			setup: func() (reflect.Value, any) {
				var target int32
				return reflect.ValueOf(&target).Elem(), float64(0.5) // Will trigger warning path but succeed
			},
			expectError: false,
		},
		{
			name: "large float64 causing int64 overflow",
			setup: func() (reflect.Value, any) {
				var target int64
				return reflect.ValueOf(&target).Elem(), float64(1e100) // Exceeds int64 range
			},
			expectError: true,
		},
		{
			name: "large float64 causing uint64 overflow",
			setup: func() (reflect.Value, any) {
				var target uint64
				return reflect.ValueOf(&target).Elem(), float64(1e100) // Exceeds uint64 range
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, input := tt.setup()

			var err error
			switch v := input.(type) {
			case float32:
				err = Float(target, v)
			case float64:
				err = Float(target, v)
			}

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestHandleEmptyStringCoverage tests specific cases for handleEmptyString
func TestHandleEmptyStringCoverage(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() reflect.Value
		input       string
		expectError bool
	}{
		{
			name: "empty string to chan type - unsupported",
			setup: func() reflect.Value {
				var target chan int
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: true,
		},
		{
			name: "empty string to func type - unsupported",
			setup: func() reflect.Value {
				var target func()
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: true,
		},
		{
			name: "empty string to array type",
			setup: func() reflect.Value {
				var target [3]int
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: true, // Arrays aren't handled in handleEmptyString
		},
		{
			name: "empty string to struct type",
			setup: func() reflect.Value {
				var target struct{ Name string }
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: true, // Structs aren't handled in handleEmptyString
		},
		{
			name: "empty string to initialized map",
			setup: func() reflect.Value {
				target := make(map[string]string)
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: false, // Maps with existing value should succeed
		},
		{
			name: "empty string to initialized pointer",
			setup: func() reflect.Value {
				initial := 42
				target := &initial
				return reflect.ValueOf(&target).Elem()
			},
			input:       "",
			expectError: false, // Initialized pointers should succeed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.setup()
			err := String(field, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

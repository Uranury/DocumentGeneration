package services_test

import (
	"RBKproject4/internal/services"
	"reflect"
	"testing"
)

func TestToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "simple struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "John",
				Age:  30,
			},
			expected: map[string]interface{}{
				"name": "John",
				"age":  float64(30), // JSON numbers unmarshal as float64
			},
			wantErr: false,
		},
		{
			name: "nested struct",
			input: struct {
				User struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"user"`
				Active bool `json:"active"`
			}{
				User: struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					Name:  "Jane",
					Email: "jane@example.com",
				},
				Active: true,
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "Jane",
					"email": "jane@example.com",
				},
				"active": true,
			},
			wantErr: false,
		},
		{
			name: "map input",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": float64(42),
				"key3": true,
			},
			wantErr: false,
		},
		{
			name: "slice with mixed types",
			input: struct {
				Items []interface{} `json:"items"`
			}{
				Items: []interface{}{"string", 123, true, nil},
			},
			expected: map[string]interface{}{
				"items": []interface{}{"string", float64(123), true, nil},
			},
			wantErr: false,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "empty struct",
			input:    struct{}{},
			expected: map[string]interface{}{},
			wantErr:  false,
		},
		{
			name: "struct with json tags",
			input: struct {
				PublicField  string `json:"public"`
				privateField string `json:"private"`
				IgnoredField string `json:"-"`
				OmitEmpty    string `json:"omit_empty,omitempty"`
			}{
				PublicField:  "visible",
				privateField: "hidden", // won't be marshaled (private)
				IgnoredField: "ignored",
				OmitEmpty:    "", // will be omitted
			},
			expected: map[string]interface{}{
				"public": "visible",
			},
			wantErr: false,
		},
		{
			name: "struct with pointers",
			input: struct {
				Name    *string `json:"name"`
				Age     *int    `json:"age"`
				Nothing *string `json:"nothing"`
			}{
				Name:    stringPtr("Alice"),
				Age:     intPtr(25),
				Nothing: nil,
			},
			expected: map[string]interface{}{
				"name":    "Alice",
				"age":     float64(25),
				"nothing": nil,
			},
			wantErr: false,
		},
		{
			name: "unmarshalable type - function",
			input: func() {
				// functions cannot be marshaled to JSON
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "unmarshalable type - channel",
			input:    make(chan int),
			expected: nil,
			wantErr:  true,
		},
		{
			name: "complex nested structure",
			input: struct {
				Data map[string]interface{} `json:"data"`
				Meta struct {
					Count int      `json:"count"`
					Tags  []string `json:"tags"`
				} `json:"meta"`
			}{
				Data: map[string]interface{}{
					"nested": map[string]interface{}{
						"deep": "value",
					},
					"array": []int{1, 2, 3},
				},
				Meta: struct {
					Count int      `json:"count"`
					Tags  []string `json:"tags"`
				}{
					Count: 3,
					Tags:  []string{"tag1", "tag2"},
				},
			},
			expected: map[string]interface{}{
				"data": map[string]interface{}{
					"nested": map[string]interface{}{
						"deep": "value",
					},
					"array": []interface{}{float64(1), float64(2), float64(3)},
				},
				"meta": map[string]interface{}{
					"count": float64(3),
					"tags":  []interface{}{"tag1", "tag2"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := services.ToMap(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ToMap() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ToMap() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToMap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test edge cases separately
func TestToMapEdgeCases(t *testing.T) {
	t.Run("large numbers", func(t *testing.T) {
		input := struct {
			BigInt   int64   `json:"big_int"`
			BigFloat float64 `json:"big_float"`
		}{
			BigInt:   9223372036854775807,     // max int64
			BigFloat: 1.7976931348623157e+308, // close to max float64
		}

		result, err := services.ToMap(input)
		if err != nil {
			t.Fatalf("ToMap() error: %v", err)
		}

		if result["big_int"] != float64(9223372036854775807) {
			t.Errorf("big_int not converted correctly")
		}

		if result["big_float"] != 1.7976931348623157e+308 {
			t.Errorf("big_float not converted correctly")
		}
	})

	t.Run("unicode strings", func(t *testing.T) {
		input := struct {
			Unicode string `json:"unicode"`
			Emoji   string `json:"emoji"`
		}{
			Unicode: "Hello 世界",
			Emoji:   "🚀 🌟",
		}

		result, err := services.ToMap(input)
		if err != nil {
			t.Fatalf("ToMap() error: %v", err)
		}

		if result["unicode"] != "Hello 世界" {
			t.Errorf("unicode string not preserved")
		}

		if result["emoji"] != "🚀 🌟" {
			t.Errorf("emoji string not preserved")
		}
	})
}

// Benchmark the function
func BenchmarkToMap(b *testing.B) {
	testStruct := struct {
		Name  string                 `json:"name"`
		Age   int                    `json:"age"`
		Data  map[string]interface{} `json:"data"`
		Items []string               `json:"items"`
	}{
		Name: "Benchmark Test",
		Age:  30,
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Items: []string{"item1", "item2", "item3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = services.ToMap(testStruct)
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

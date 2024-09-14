package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionAnalyzer(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected []FunctionInfo
	}{
		{
			name: "Simple function",
			code: `
				package main
				func add(a, b int) int {
					return a + b
				}
			`,
			expected: []FunctionInfo{
				{
					Name: "add",
					Parameters: []ParameterInfo{
						{Name: "a", Type: "int"},
						{Name: "b", Type: "int"},
					},
					Results: []ParameterInfo{
						{Type: "int"},
					},
				},
			},
		},
		{
			name: "Method with receiver",
			code: `
				package main
				type Calculator struct{}
				func (c *Calculator) Multiply(a, b int) int {
					return a * b
				}
			`,
			expected: []FunctionInfo{
				{
					Name:     "Multiply",
					Receiver: "*Calculator",
					IsMethod: true,
					Parameters: []ParameterInfo{
						{Name: "a", Type: "int"},
						{Name: "b", Type: "int"},
					},
					Results: []ParameterInfo{
						{Type: "int"},
					},
				},
			},
		},
		{
			name: "Generic function",
			code: `
				package main
				func identity[T any](value T) T {
					return value
				}
			`,
			expected: []FunctionInfo{
				{
					Name:      "identity",
					IsGeneric: true,
					Parameters: []ParameterInfo{
						{Name: "value", Type: "T"},
					},
					Results: []ParameterInfo{
						{Type: "T"},
					},
				},
			},
		},
		{
			name: "Function with complex types",
			code: `
				package main
				func process(data map[string][]int, callback func(int) bool) (result *string, err error) {
					return nil, nil
				}
			`,
			expected: []FunctionInfo{
				{
					Name: "process",
					Parameters: []ParameterInfo{
						{Name: "data", Type: "map[string][]int"},
						{Name: "callback", Type: "func()"},
					},
					Results: []ParameterInfo{
						{Name: "result", Type: "*string"},
						{Name: "err", Type: "error"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := NewFunctionAnalyzer()
			err := analyzer.AnalyzeFile("test.go")
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, analyzer.Functions)
		})
	}
}

func TestAnalyzeDirectory(t *testing.T) {
	// Create a temporary directory with some Go files
	// Run AnalyzeDirectory on this temporary directory
	// Assert that all files are analyzed correctly
	// Clean up the temporary directory
}
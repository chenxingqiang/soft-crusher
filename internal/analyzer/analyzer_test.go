package analyzer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
						{Name: "callback", Type: "func(int) bool"},
					},
					Results: []ParameterInfo{
						{Name: "result", Type: "*string"},
						{Name: "err", Type: "error"},
					},
				},
			},
		},
		{
			name: "Multiple functions",
			code: `
				package main
				func foo() {}
				func bar(x int) string { return "" }
				func (s *Service) baz(y float64) (bool, error) { return false, nil }
			`,
			expected: []FunctionInfo{
				{
					Name: "foo",
				},
				{
					Name: "bar",
					Parameters: []ParameterInfo{
						{Name: "x", Type: "int"},
					},
					Results: []ParameterInfo{
						{Type: "string"},
					},
				},
				{
					Name:     "baz",
					Receiver: "*Service",
					IsMethod: true,
					Parameters: []ParameterInfo{
						{Name: "y", Type: "float64"},
					},
					Results: []ParameterInfo{
						{Type: "bool"},
						{Type: "error"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := ioutil.TempFile("", "test_*.go")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			_, err = tempFile.Write([]byte(tc.code))
			require.NoError(t, err)
			tempFile.Close()

			analyzer := NewFunctionAnalyzer()
			err = analyzer.AnalyzeFile(tempFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, analyzer.Functions)
		})
	}
}

func TestAnalyzeDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "test_analyzer")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	files := map[string]string{
		"file1.go": `
			package main
			func foo() {}
			func bar(x int) string { return "" }
		`,
		"file2.go": `
			package main
			func (s *Service) baz(y float64) (bool, error) { return false, nil }
		`,
		"ignored.txt": "This file should be ignored",
	}

	for filename, content := range files {
		err := ioutil.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Run AnalyzeDirectory
	analyzer := NewFunctionAnalyzer()
	err = analyzer.AnalyzeDirectory(tempDir)
	assert.NoError(t, err)

	// Check results
	expected := []FunctionInfo{
		{
			Name: "foo",
		},
		{
			Name: "bar",
			Parameters: []ParameterInfo{
				{Name: "x", Type: "int"},
			},
			Results: []ParameterInfo{
				{Type: "string"},
			},
		},
		{
			Name:     "baz",
			Receiver: "*Service",
			IsMethod: true,
			Parameters: []ParameterInfo{
				{Name: "y", Type: "float64"},
			},
			Results: []ParameterInfo{
				{Type: "bool"},
				{Type: "error"},
			},
		},
	}

	assert.ElementsMatch(t, expected, analyzer.Functions)
}

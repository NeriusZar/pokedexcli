package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			"hello world",
			[]string{"hello", "world"},
		},
		{
			"   hello world   ",
			[]string{"hello", "world"},
		},
		{
			"Hello world",
			[]string{"hello", "world"},
		},
		{
			"Hello world",
			[]string{"hello", "world"},
		},
	}

	for _, c := range testCases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			failTest(c.input, c.expected, actual, t)
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				failTest(c.input, c.expected, actual, t)
				break
			}
		}
	}
}

func failTest(input string, expected []string, actual []string, t *testing.T) {
	t.Errorf(`
	Test Failed:
	  input: %s
	  expected: %s
	  actual: %s
	`, input, formatSlice(expected), formatSlice(actual))
}

func formatSlice(words []string) string {
	result := "["
	for i, w := range words {
		if i == len(words)-1 {
			result += w
		} else {
			result += fmt.Sprintf("%s, ", w)
		}
	}

	result += "]"

	return result
}

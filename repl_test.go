package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "		 hello world	 ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Error to take the clean input: unexpected length")
		}
		for i := range actual {
			world := actual[i]
			expectedWord := c.expected[i]

			if world != expectedWord {
				t.Errorf("%s and %s are different, error to parse", world, expectedWord)
			}
		}
	}
}

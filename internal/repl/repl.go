package repl

import "strings"

func CleanInput(text string) []string {
	sLower := strings.ToLower(text)
	return strings.Fields(sLower)
}

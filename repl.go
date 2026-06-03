package main

import "strings"

func cleanInput(text string) []string {
	sLower := strings.ToLower(text)
	return strings.Fields(sLower)
}

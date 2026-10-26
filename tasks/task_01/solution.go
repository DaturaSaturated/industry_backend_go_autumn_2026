package main

import (
	"strings"
)

func greet(name string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) != 0 {
		return "Hello, " + trimmed + "!"
	}
	return "Hello, World!"

}

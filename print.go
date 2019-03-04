package main

import (
	"fmt"
	"os"
	"strings"
)

func msg(s string, args ...interface{}) {
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	fmt.Printf(s, args...)
}

func warn(s string, args ...interface{}) {
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	fmt.Fprintf(os.Stderr, s, args...)
}

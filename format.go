package main

import (
	"fmt"
	"os"
)

const (
	NONE = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
)

func printError(text string, err error) {
	fmt.Fprintf(os.Stderr, format(RED, "%s: %s\n"), text, err)
}

func format(c int, text any) string {
	const escape = "\x1b"

	color := func(c int) string {
		var term string
		if c != NONE {
			term = "3"
		}

		return fmt.Sprintf("%s[%s%dm", escape, term, c)
	}

	return fmt.Sprintf("%s%s%s", color(c), text, color(NONE))
}

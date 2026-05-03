//go:build !windows

package main

import (
	"bufio"
	"os"

	"golang.org/x/term"
)

// readLineHidden reads a line with echo disabled, regardless of ANSI support.
func readLineHidden() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err == nil {
			defer term.Restore(fd, state)
		}
	}
	// In raw mode the terminal will not echo; we read until newline manually.
	r := bufio.NewReader(os.Stdin)
	var buf []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		switch b {
		case '\r', '\n':
			return string(buf), nil
		case 0x7f, 0x08: // DEL / BS
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		case 0x03: // Ctrl+C
			return "", errInterrupted
		case 0x04: // Ctrl+D
			if len(buf) == 0 {
				return "", errInterrupted
			}
		default:
			if b >= 32 || b == '\t' {
				buf = append(buf, b)
			}
		}
	}
}

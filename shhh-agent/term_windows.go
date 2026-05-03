//go:build windows

package main

import (
	"bufio"
	"os"

	"golang.org/x/term"
)

// readLineHidden uses Windows console raw mode so even paste is invisible.
func readLineHidden() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err == nil {
			defer term.Restore(fd, state)
		}
	}
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
		case 0x7f, 0x08:
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		case 0x03:
			return "", errInterrupted
		case 0x04:
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

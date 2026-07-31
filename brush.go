// Package brush styles terminal output using ANSI escape codes.
package brush

import (
	"os"
	"strconv"
)

var Disabled = os.Getenv("NO_COLOR") != ""

const (
	esc   = "\x1b["
	end   = "m"
	reset = "\x1b[0m"
)

type Color int8

const (
	NoColor Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func (c Color) fgCode() int {
	// Normal colors
	if c < BrightBlack {
		return int(c-Black) + 30
	}
	// Bright colors
	return int(c-BrightBlack) + 90
}

func (c Color) Paint(s string) string {
	if s == "" || Disabled || c == NoColor {
		return s
	}
	return esc + strconv.Itoa(c.fgCode()) + end + s + reset
}

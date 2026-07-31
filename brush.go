// Package brush styles terminal output using ANSI escape codes.
package brush

import "strconv"

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
	if s == "" || c == NoColor {
		return s
	}
	return "\x1b[" + strconv.Itoa(c.fgCode()) + "m" + s + "\x1b[0m"
}

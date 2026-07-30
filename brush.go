// Package brush styles terminal text with ANSI colors.
package brush

import "fmt"

type Color int

const (
	Black Color = iota + 30
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
	if c >= BrightBlack {
		return int(c-BrightBlack) + 90
	}
	return int(c-Black) + 30
}

func (c Color) bgCode() int {
	if c >= BrightBlack {
		return int(c-BrightBlack) + 100
	}
	return int(c-Black) + 40
}

func (c Color) Style(text string) string {
	code := c.fgCode()
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, text)
}

type Styler interface {
	Style(text string) string
}

func Paint(s Styler, a ...any) string {
	return s.Style(fmt.Sprint(a...))
}

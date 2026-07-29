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
	_ // Code 38 is for 256/24-bit foreground color.
	Default

	BgBlack
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	_ // Code 48 is for 256/24-bit background color.
	BgDefault
)

const (
	BrightBlack Color = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

const (
	BgBrightBlack Color = iota + 100
	BgBrightRed
	BgBrightGreen
	BgBrightYellow
	BgBrightBlue
	BgBrightMagenta
	BgBrightCyan
	BgBrightWhite
)

func (c Color) Style(text string) string {
	if c == Default {
		return text
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, text)
}

type Styler interface {
	Style(text string) string
}

func Paint(s Styler, a ...any) string {
	return s.Style(fmt.Sprint(a...))
}

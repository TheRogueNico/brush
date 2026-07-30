// Package brush styles terminal text with ANSI colors.
package brush

import "fmt"

type Color int

const (
	Black Color = iota
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
	return int(c) + 30
}

func (c Color) bgCode() int {
	if c >= BrightBlack {
		return int(c-BrightBlack) + 100
	}
	return int(c) + 40
}

type Attribute int

const (
	Bold Attribute = 1 << iota
	Dim
	Italic
	Underline
	Strikethrough
)

type Look struct {
	Foreground Color
	Background Color
	Attributes Attribute
}

func (c Color) Style(text string) string {
	code := c.fgCode()
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, text)
}

func (l Look) Style(text string) string {
	// TODO: Write the actual logic
	return ""
}

type Styler interface {
	Style(text string) string
}

func Paint(s Styler, a ...any) string {
	return s.Style(fmt.Sprint(a...))
}

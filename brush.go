// Package brush styles terminal text with ANSI colors.
package brush

import (
	"fmt"
	"strconv"
	"strings"
)

type Color int

const (
	Default Color = iota
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
	if c == Default {
		return text
	}
	code := c.fgCode()
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, text)
}

type Attribute int

const (
	Bold Attribute = 1 << iota
	Dim
	Italic
	Underline
	Strikethrough
)

var attributeCodes = []struct {
	attr Attribute
	code int
}{
	{Bold, 1},
	{Dim, 2},
	{Italic, 3},
	{Underline, 4},
	{Strikethrough, 9},
}

type Look struct {
	Foreground Color
	Background Color
	Attributes Attribute
}

func (l Look) Style(text string) string {
	var codes []string

	if l.Foreground != Default {
		codes = append(codes, strconv.Itoa(l.Foreground.fgCode()))
	}
	if l.Background != Default {
		codes = append(codes, strconv.Itoa(l.Background.bgCode()))
	}

	for _, ac := range attributeCodes {
		// Checks if attribute bit is set using bitwise AND
		if l.Attributes&ac.attr != 0 {
			codes = append(codes, strconv.Itoa(ac.code))
		}
	}

	if len(codes) == 0 {
		return text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", strings.Join(codes, ";"), text)
}

type Styler interface {
	Style(text string) string
}

func Paint(s Styler, a ...any) string {
	return s.Style(fmt.Sprint(a...))
}

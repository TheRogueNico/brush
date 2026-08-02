// Package brush styles terminal output using ANSI escape codes.
package brush

import (
	"os"
	"strconv"
	"strings"
)

func noColorEnvSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

// Disabled controls whether Paint methods apply color.
// This does not disable italic, bold or other styles.
var Disabled = noColorEnvSet()

// SGR escape sequences
const (
	esc   = "\x1b["
	end   = "m"
	reset = "\x1b[0m"
)

// Color represents one of the 16 standard ANSI terminal colors.
type Color int8

// The 16 standard ANSI terminal colors.
// NoColor is the zero value and represents the absence of a color.
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

// fgCode returns the SGR foreground code for c.
// c must not be NoColor. It is up to the caller to handle this case.
func (c Color) fgCode() int {
	// Normal colors
	if c < BrightBlack {
		return int(c-Black) + 30
	}
	// Bright colors
	return int(c-BrightBlack) + 90
}

func (c Color) bgCode() int {
	// Normal colors
	if c < BrightBlack {
		return int(c-Black) + 40
	}
	// Bright colors
	return int(c-BrightBlack) + 100
}

// Paint wraps s in the ANSI escape sequence for the foreground color c,
// followed by a reset sequence. If c is NoColor, coloring is disabled
// globally via Disabled, or s is empty, s is returned unchanged.
func (c Color) Paint(s string) string {
	if s == "" || Disabled || c == NoColor {
		return s
	}
	return esc + strconv.Itoa(c.fgCode()) + end + s + reset
}

var colorNames = [...]string{
	"NoColor",
	"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White",
	"BrightBlack", "BrightRed", "BrightGreen", "BrightYellow",
	"BrightBlue", "BrightMagenta", "BrightCyan", "BrightWhite",
}

// String returns the name of the color, used by fmt.Stringer.
func (c Color) String() string {
	if c < 0 || int(c) >= len(colorNames) {
		return "Color(" + strconv.Itoa(int(c)) + ")"
	}
	return colorNames[c]
}

type Attribute uint8

const (
	Bold Attribute = 1 << iota
	Faint
	Italic
	Underline
	Blink
	Reverse
	Strikethrough
)

var attributeTable = []struct {
	attr Attribute
	code int
	name string
}{
	{Bold, 1, "Bold"},
	{Faint, 2, "Faint"},
	{Italic, 3, "Italic"},
	{Underline, 4, "Underline"},
	{Blink, 5, "Blink"},
	{Reverse, 7, "Reverse"},
	{Strikethrough, 9, "Strikethrough"},
}

type Style struct {
	Foreground Color
	Background Color
	Attributes Attribute
}

func (a Attribute) String() string {
	var names []string
	for _, t := range attributeTable {
		if a&t.attr != 0 {
			names = append(names, t.name)
			a &^= t.attr
		}
	}
	if a != 0 {
		names = append(names, "Attribute("+strconv.Itoa(int(a))+")")
	}
	return strings.Join(names, "|")
}

func (s Style) Paint(text string) string {
	if text == "" {
		return text
	}

	var codes []string
	for _, a := range attributeTable {
		if s.Attributes&a.attr != 0 {
			codes = append(codes, strconv.Itoa(a.code))
		}
	}
	if !Disabled {
		if s.Foreground != NoColor {
			codes = append(codes, strconv.Itoa(s.Foreground.fgCode()))
		}
		if s.Background != NoColor {
			codes = append(codes, strconv.Itoa(s.Background.bgCode()))
		}
	}

	if len(codes) == 0 {
		return text
	}
	return esc + strings.Join(codes, ";") + end + text + reset
}

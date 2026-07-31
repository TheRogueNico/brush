// Package brush styles terminal output using ANSI escape codes.
package brush

import (
	"os"
	"strconv"
	"strings"
)

// Disabled controls whether Paint methods apply color.
// This does not disable italic, bold or other styles.
var Disabled = os.Getenv("NO_COLOR") != ""

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

type Style struct {
	Foreground Color
	Background Color
	Attributes Attribute
}

func (s Style) Paint(text string) string {
	var codes []string
	if text != "" && !Disabled && s.Foreground != NoColor {
		codes = append(codes, strconv.Itoa(s.Foreground.fgCode()))
	}

	return esc + strings.Join(codes, ";") + end + text + reset
}

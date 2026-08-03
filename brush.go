// Package brush wraps text in ANSI escape codes for color and text attributes.
//
// A Color applies one of the 16 standard ANSI colors:
//
//	fmt.Println(brush.Red.Paint("error"))
//
// A Style combines a foreground Color, a background Color, and any
// number of Attributes:
//
//	warn := brush.Style{Foreground: brush.Yellow, Attributes: brush.Bold | brush.Italic}
//	fmt.Println(warn.Paint("careful"))
//
// A Stencil colors runs of uppercase letters, lowercase letters,
// digits, punctuation, and symbols with separate colors:
//
//	rule := brush.Stencil{Upper: brush.Red, Digit: brush.Blue}
//	fmt.Println(rule.Paint("Error 403"))
//
// Color, Style and Stencil all implement Painter, the package's only
// interface, so a custom coloring rule needs only a Paint method to be
// used the same way as the built-in ones. PaintFunc adapts an ordinary
// function into a Painter:
//
//	shout := brush.PaintFunc(func(text string) string {
//		return brush.Red.Paint(strings.ToUpper(text))
//	})
//	fmt.Println(shout.Paint("hello"))
//
// Output is colored unless Disabled is true, which happens
// automatically when the NO_COLOR environment variable is set to a
// non-empty value. See https://no-color.org.
package brush

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

// noColorEnvSet reports whether NO_COLOR is set to a non-empty value.
func noColorEnvSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

// Disabled controls whether Paint methods apply color.
// This does not disable italic, bold or other attributes.
var Disabled = noColorEnvSet()

// SGR escape sequences
const (
	esc   = "\x1b["
	end   = "m"
	reset = "\x1b[0m"
)

// Painter applies a coloring rule to text. Color and Style both
// implement Painter, so a custom coloring rule only needs a Paint
// method to be used interchangeably with them.
type Painter interface {
	Paint(text string) string
}

// PaintFunc adapts an ordinary function to a Painter.
type PaintFunc func(text string) string

// Paint calls f(text).
func (f PaintFunc) Paint(text string) string {
	return f(text)
}

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

// bgCode returns the SGR background code for c.
// c must not be NoColor. It is up to the caller to handle this case.
func (c Color) bgCode() int {
	// Normal colors
	if c < BrightBlack {
		return int(c-Black) + 40
	}
	// Bright colors
	return int(c-BrightBlack) + 100
}

// Paint wraps text in the ANSI escape sequence for the foreground color c,
// followed by a reset sequence. If c is NoColor, coloring is disabled
// globally via Disabled, or text is empty, text is returned unchanged.
func (c Color) Paint(text string) string {
	if text == "" || Disabled || c == NoColor {
		return text
	}
	return esc + strconv.Itoa(c.fgCode()) + end + text + reset
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

// Attribute represents a text attribute, such as bold or italic, that
// can be combined with a Color to form a Style. Attributes are bit
// flags, combined with the bitwise OR operator, e.g. Bold|Italic.
type Attribute uint8

// The supported text attributes. Not every terminal renders every
// attribute; Blink and Strikethrough support varies most.
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

// Style combines a foreground Color, a background Color, and any
// number of Attributes into a single, reusable coloring rule.
//
// The zero value applies nothing; Paint returns its argument unchanged.
type Style struct {
	Foreground Color
	Background Color
	Attributes Attribute
}

// String returns the pipe-separated names of the set attributes, used by fmt.Stringer.
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

// Paint wraps text in the ANSI escape sequence for the style's colors
// and attributes, followed by a reset sequence. Disabled suppresses
// only the color codes; attributes are always applied. Empty text is
// returned unchanged.
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

// Stencil colors uppercase letters, lowercase letters, digits,
// punctuation, and symbols using separate colors.
// Runes outside these five classes are left unstyled.
//
// The zero value leaves everything unstyled.
type Stencil struct {
	Upper  Color
	Lower  Color
	Digit  Color
	Punct  Color
	Symbol Color
}

// colorFor returns the Color for r's class or NoColor if r matches none of them.
func (s Stencil) colorFor(r rune) Color {
	switch {
	case unicode.IsUpper(r):
		return s.Upper
	case unicode.IsLower(r):
		return s.Lower
	case unicode.IsDigit(r):
		return s.Digit
	case unicode.IsPunct(r):
		return s.Punct
	case unicode.IsSymbol(r):
		return s.Symbol
	default:
		return NoColor
	}
}

// Paint colors each run of same-class characters, switching colors
// at class boundaries and ending with a reset if any color was
// applied. If Disabled or text is empty, text is returned unchanged.
func (s Stencil) Paint(text string) string {
	if text == "" || Disabled {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))

	start := 0
	current := NoColor
	for i, r := range text {
		class := s.colorFor(r)
		if class == current {
			continue
		}
		b.WriteString(text[start:i])
		start = i
		if class == NoColor {
			b.WriteString(reset)
		} else {
			b.WriteString(esc)
			b.WriteString(strconv.Itoa(class.fgCode()))
			b.WriteString(end)
		}
		current = class
	}
	b.WriteString(text[start:])
	if current != NoColor {
		b.WriteString(reset)
	}
	return b.String()
}

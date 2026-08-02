package brush

import (
	"os"
	"strings"
	"testing"
)

func TestNoColorEnvSet(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		original, ok := os.LookupEnv("NO_COLOR")
		os.Unsetenv("NO_COLOR")
		t.Cleanup(func() {
			if ok {
				os.Setenv("NO_COLOR", original)
			}
		})
		if got := noColorEnvSet(); got != false {
			t.Errorf("noColorEnvSet() = %v, want false", got)
		}
	})

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"empty", "", false},
		{"set to 1", "1", true},
		{"set to anything", "yes", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tt.value)
			if got := noColorEnvSet(); got != tt.want {
				t.Errorf("noColorEnvSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Paint(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		in    string
		want  string
	}{
		{"red", Red, "Hello, World!", "\x1b[31mHello, World!\x1b[0m"},
		{"black", Black, "Hello, World!", "\x1b[30mHello, World!\x1b[0m"},
		{"bright white", BrightWhite, "Hello, World!", "\x1b[97mHello, World!\x1b[0m"},
		{"zero value is NoColor", Color(0), "Hello, World!", "Hello, World!"},
		{"empty string", White, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.Paint(tt.in); got != tt.want {
				t.Errorf("Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_String(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{"no color", NoColor, "NoColor"},
		{"red", Red, "Red"},
		{"bright white", BrightWhite, "BrightWhite"},
		{"undefined color", Color(99), "Color(99)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStyle_Paint(t *testing.T) {
	infoStyle := Style{Foreground: Blue}
	warningStyle := Style{Foreground: Yellow, Attributes: Italic}
	errorStyle := Style{Foreground: Black, Background: BrightRed, Attributes: Italic | Bold}

	tests := []struct {
		name  string
		style Style
		in    string
		want  string
	}{
		{
			name:  "zero value",
			style: Style{},
			in:    "no style",
			want:  "no style",
		},
		{
			name:  "foreground only",
			style: infoStyle,
			in:    "starting up...",
			want:  "\x1b[34mstarting up...\x1b[0m",
		},
		{
			name:  "background only",
			style: Style{Background: Blue},
			in:    "Highlight",
			want:  "\x1b[44mHighlight\x1b[0m",
		},
		{
			name:  "attributes only",
			style: Style{Attributes: Underline},
			in:    "diceware",
			want:  "\x1b[4mdiceware\x1b[0m",
		},
		{
			name:  "attribute and foreground",
			style: warningStyle,
			in:    "Disk is nearly full",
			want:  "\x1b[3;33mDisk is nearly full\x1b[0m",
		},
		{
			name:  "attributes foreground and background",
			style: errorStyle,
			in:    "Connection lost!",
			want:  "\x1b[1;3;30;101mConnection lost!\x1b[0m",
		},
		{
			name:  "explicit NoColor is equivalent to unset",
			style: Style{Foreground: NoColor, Background: NoColor, Attributes: Bold},
			in:    "x",
			want:  "\x1b[1mx\x1b[0m",
		},
		{
			name:  "empty string",
			style: Style{Foreground: Red, Attributes: Bold},
			in:    "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Paint(tt.in)
			if got != tt.want {
				t.Errorf("Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStyle_PaintDisabled(t *testing.T) {
	Disabled = true
	defer func() { Disabled = false }()

	tests := []struct {
		name  string
		style Style
		in    string
		want  string
	}{
		{
			name:  "color is suppressed",
			style: Style{Foreground: Red, Background: Black},
			in:    "Hello",
			want:  "Hello",
		},
		{
			name:  "attributes not suppressed",
			style: Style{Attributes: Bold | Italic},
			in:    "Hello",
			want:  "\x1b[1;3mHello\x1b[0m",
		},
		{
			name:  "attributes alongside suppressed color",
			style: Style{Foreground: Red, Attributes: Bold},
			in:    "Hello",
			want:  "\x1b[1mHello\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.style.Paint(tt.in); got != tt.want {
				t.Errorf("Paint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAttribute_String(t *testing.T) {
	tests := []struct {
		name string
		attr Attribute
		want string
	}{
		{"zero", 0, ""},
		{"single attribute", Bold, "Bold"},
		{"two attributes", Bold | Italic, "Bold|Italic"},
		{"number", 7, "Bold|Faint|Italic"},
		{"out of range", 128, "Attribute(128)"},
		{"mixed valid and out of range", 129, "Bold|Attribute(128)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.attr.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaintFunc_Paint(t *testing.T) {
	tests := []struct {
		name string
		f    PaintFunc
		in   string
		want string
	}{
		{
			name: "identity",
			f:    func(s string) string { return s },
			in:   "x",
			want: "x",
		},
		{
			name: "wraps in brackets",
			f:    func(s string) string { return "[" + s + "]" },
			in:   "x",
			want: "[x]",
		},
		{
			name: "delegates to a Color",
			f:    Red.Paint,
			in:   "x",
			want: "\x1b[31mx\x1b[0m",
		},
		{
			name: "composes with another Painter",
			f: func(s string) string {
				return Style{Attributes: Bold}.Paint(strings.ToUpper(s))
			},
			in:   "hello",
			want: "\x1b[1mHELLO\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Paint(tt.in); got != tt.want {
				t.Errorf("Paint(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCharClass_Paint(t *testing.T) {
	rule := Stencil{Upper: Red, Lower: Green, Digit: Blue, Punct: Yellow, Symbol: Magenta}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty string", "", ""},
		{"single class", "ABC", "\x1b[31mABC\x1b[0m"},
		{"mixed classes no reset between classes", "AB12!!", "\x1b[31mAB\x1b[34m12\x1b[33m!!\x1b[0m"},
		{"symbol class", "1+1=2", "\x1b[34m1\x1b[35m+\x1b[34m1\x1b[35m=\x1b[34m2\x1b[0m"},
		{"reset on return to unstyled", "AB cd", "\x1b[31mAB\x1b[0m \x1b[32mcd\x1b[0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.Paint(tt.in)
			if got != tt.want {
				t.Errorf("Paint(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}

	t.Run("zero value", func(t *testing.T) {
		var zero Stencil
		got := zero.Paint("AB 12 !!")
		if got != "AB 12 !!" {
			t.Errorf("Paint() = %q, want unchanged", got)
		}
	})
}

// Compile-time assertions that all types satisfy Painter.
var (
	_ Painter = Red
	_ Painter = Style{}
	_ Painter = PaintFunc(nil)
	_ Painter = Stencil{}
)

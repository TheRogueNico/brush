package brush

import (
	"os"
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
		name  string // description of this test case
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
		name  string // description of this test case
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
		name  string // description of this test case
		style Style
		text  string
		want  string
	}{
		{
			name:  "zero value",
			style: Style{},
			text:  "no style",
			want:  "no style",
		},
		{
			name:  "foreground only",
			style: infoStyle,
			text:  "starting up...",
			want:  "\x1b[34mstarting up...\x1b[0m",
		},
		{
			name:  "background only",
			style: Style{Background: Blue},
			text:  "Highlight",
			want:  "\x1b[44mHighlight\x1b[0m",
		},
		{
			name:  "attributes only",
			style: Style{Attributes: Underline},
			text:  "diceware",
			want:  "\x1b[4mdiceware\x1b[0m",
		},
		{
			name:  "attribute and foreground",
			style: warningStyle,
			text:  "Disk is nearly full",
			want:  "\x1b[3;33mDisk is nearly full\x1b[0m",
		},
		{
			name:  "attributes foreground and background",
			style: errorStyle,
			text:  "Connection lost!",
			want:  "\x1b[1;3;30;101mConnection lost!\x1b[0m",
		},
		{
			name:  "explicit NoColor is equivalent to unset",
			style: Style{Foreground: NoColor, Background: NoColor, Attributes: Bold},
			text:  "x",
			want:  "\x1b[1mx\x1b[0m",
		},
		{
			name:  "empty string",
			style: Style{Foreground: Red, Attributes: Bold},
			text:  "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Paint(tt.text)
			if got != tt.want {
				t.Errorf("Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttribute_String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
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

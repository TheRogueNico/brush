package brush

import "testing"

func TestColorPaint(t *testing.T) {
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

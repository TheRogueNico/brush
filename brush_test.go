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

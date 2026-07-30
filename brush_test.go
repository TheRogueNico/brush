package brush

import "testing"

func TestPaint_Color(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    Styler
		a    []any
		want string
	}{
		{
			name: "red foreground",
			s:    Red,
			a:    []any{"Hello, ", "World", "!"},
			want: "\x1b[31mHello, World!\x1b[0m",
		},
		{
			name: "bright red foreground",
			s:    BrightRed,
			a:    []any{"Hello, ", "World", "!"},
			want: "\x1b[91mHello, World!\x1b[0m",
		},
		{
			name: "default color no styling",
			s:    Default,
			a:    []any{"Hello, ", "World", "!"},
			want: "Hello, World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Paint(tt.s, tt.a...)
			if got != tt.want {
				t.Errorf("Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaint_Look(t *testing.T) {
	infoLook := Look{
		Foreground: BrightBlue,
		Attributes: Italic,
	}
	warningLook := Look{
		Foreground: Yellow,
		Attributes: Bold,
	}
	errorLook := Look{
		Foreground: White,
		Background: Red,
		Attributes: Bold | Italic,
	}
	allAttributes := Look{
		Attributes: Bold | Dim | Italic | Underline | Strikethrough,
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    Styler
		a    []any
		want string
	}{
		{
			name: "bright blue italic",
			s:    infoLook,
			a:    []any{"starting up..."},
			want: "\x1b[94;3mstarting up...\x1b[0m",
		},
		{
			name: "yellow bold",
			s:    warningLook,
			a:    []any{"disk usage at ", 90, "%"},
			want: "\x1b[33;1mdisk usage at 90%\x1b[0m",
		},
		{
			name: "white on red bold italic",
			s:    errorLook,
			a:    []any{"FATAL: ", "connection refused"},
			want: "\x1b[37;41;1;3mFATAL: connection refused\x1b[0m",
		},
		{
			name: "all attributes no color",
			s:    allAttributes,
			a:    []any{"kitchen sink"},
			want: "\x1b[1;2;3;4;9mkitchen sink\x1b[0m",
		},
		{
			name: "zero value Look no styling",
			s:    Look{},
			a:    []any{"plain ", "text"},
			want: "plain text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Paint(tt.s, tt.a...)
			if got != tt.want {
				t.Errorf("Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}

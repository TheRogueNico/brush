package brush

import "testing"

func TestPaint(t *testing.T) {
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

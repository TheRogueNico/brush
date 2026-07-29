package brush

import "testing"

func TestPaint(t *testing.T) {
	got := Paint(Red, "Hello,", " World", "!")
	want := "\033[31mHello, World!\033[0m"

	if got != want {
		t.Errorf("Paint: got %v want %v", got, want)
	}
}

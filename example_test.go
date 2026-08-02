package brush_test

import (
	"fmt"
	"strings"

	"github.com/TheRogueNico/brush"
)

func ExampleColor_Paint() {
	fmt.Println(brush.Red.Paint("error"))
}

func ExampleStyle_Paint() {
	warn := brush.Style{Foreground: brush.Yellow, Attributes: brush.Bold | brush.Italic}
	fmt.Println(warn.Paint("Warning!"))
}

// ExamplePaintFunc shows a custom coloring rule written as a plain
// function and adapted into a Painter.
func ExamplePaintFunc() {
	// To simplify output
	brush.Disabled = true
	defer func() { brush.Disabled = false }()

	shout := brush.PaintFunc(func(text string) string {
		return brush.Red.Paint(strings.ToUpper(text))
	})
	fmt.Println(shout.Paint("hello"))
	// Output: HELLO
}

// lowercaseOnly is a custom coloring rule written as a type with its own
// Paint method: it colors only the lowercase runs of a string, leaving
// everything else untouched. Because it implements Painter, it's used
// exactly like a Color or a Style.
type lowercaseOnly struct {
	Color brush.Color
}

func (l lowercaseOnly) Paint(text string) string {
	var b strings.Builder
	start := -1
	flush := func(end int) {
		if start == -1 {
			return
		}
		b.WriteString(l.Color.Paint(text[start:end]))
		start = -1
	}
	for i, r := range text {
		if r >= 'a' && r <= 'z' {
			if start == -1 {
				start = i
			}
			continue
		}
		flush(i)
		b.WriteRune(r)
	}
	flush(len(text))
	return b.String()
}

func Example_customRule() {
	// To simplify output
	brush.Disabled = true
	defer func() { brush.Disabled = false }()

	rule := lowercaseOnly{brush.Green}
	fmt.Println(rule.Paint("Hello World"))
	// Output: Hello World
}

// Package brush styles terminal text with ANSI colors.
package brush

import "fmt"

type Color int

const Red Color = 31

func (c Color) Style(text string) string {
	return "\033[31m" + text + "\033[0m"
}

type Styler interface {
	Style(text string) string
}

func Paint(s Styler, a ...any) string {
	return s.Style(fmt.Sprint(a...))
}

# brush

A tiny, dependency-free ANSI styling library for Go.

`brush` wraps strings with terminal escape codes to add colors and text attributes without bringing in complexity. It's small by design and it won't try to print anything itself.

Built mostly for fun, learning, and small projects that need lightweight terminal styling.

## Installation

Requires Go 1.22 or newer.

```bash
go get github.com/TheRogueNico/brush
```

Import it:

```go
import "github.com/TheRogueNico/brush"
```

## Usage

### Colors

Apply a color directly:

```go
fmt.Println(brush.Red.Paint("error"))
fmt.Println(brush.Green.Paint("success"))
fmt.Println(brush.Yellow.Paint("warning"))
```

### Styles

Combine colors and attributes into reusable styles:

```go
warning := brush.Style{
	Foreground: brush.Yellow,
	Attributes: brush.Bold | brush.Italic,
}

fmt.Println(warning.Paint("Warning!"))
```

Styles can also apply background colors:

```go
highlight := brush.Style{
	Foreground: brush.White,
	Background: brush.Blue,
}

fmt.Println(highlight.Paint("highlighted text"))
```

### Getting your hands painty

Anything implementing `Painter` can be used as a text painter.

```go
type Painter interface {
	Paint(text string) string
}
```

Create custom behavior with `PaintFunc`:

```go
shout := brush.PaintFunc(func(text string) string {
	return brush.Red.Paint(strings.ToUpper(text))
})

fmt.Println(shout.Paint("hello"))
```

`PaintFunc` is for stateless rules. When a rule needs its own state, such as a color chosen at random on every call, or one that only touches part of the string, write it as a type with its own `Paint` method instead:

```go
type lowercaseOnly struct {
	Color brush.Color
}

func (l lowercaseOnly) Paint(text string) string {
	// return text with ANSI codes
 }
 ```

## Disable Colors

Color output can be disabled globally:

```go
brush.Disabled = true
```

When disabled, `Paint` returns text without color escape sequences. Text attributes such as bold and italic are still applied.

### NO_COLOR Support

`brush` follows the [NO_COLOR](https://no-color.org/) convention.

If the `NO_COLOR` environment variable is set to any non-empty value, colors are automatically disabled.

`brush` does not attempt to detect whether output is going to a terminal, a file, a pipe, or another destination. It is a simple ANSI wrapper and will not guess how its output will be used. If you need to disable colors, use `brush.Disabled` or the `NO_COLOR` environment variable.

## API Overview

### `Color`

Represents one of the standard ANSI colors.

Available colors:

```go
NoColor // This is the default value
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
```

Example:

```go
brush.Cyan.Paint("hello")
```

---

### `Attribute`

Text attributes can be combined with bitwise OR:

```go
brush.Bold | brush.Underline
```

Available attributes:

```go
Bold
Faint
Italic
Underline
Blink
Reverse
Strikethrough
```

---

### `Style`

A reusable combination of colors and attributes:

```go
type Style struct {
	Foreground Color
	Background Color
	Attributes Attribute
}
```

Example:

```go
myStyle := brush.Style{
	Foreground: brush.BrightGreen,
	Background: brush.Black,
	Attributes: brush.Bold | brush.Italic | brush.Underline,
}
```

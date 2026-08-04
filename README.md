# brush

A tiny, dependency-free ANSI styling library for Go.

`brush` wraps strings with terminal escape codes to add colors and text attributes without bringing in complexity. It's small by design and it won't try to print anything itself.

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

Apply one of the 16 ANSI colors directly:

```go
fmt.Println(brush.Red.Paint("error"))
fmt.Println(brush.Green.Paint("success"))
fmt.Println(brush.Yellow.Paint("warning"))
```

#### Available colors

`NoColor`, This is the default value

```
Black, Red, Green, Yellow, Blue, Magenta, Cyan, White,

BrightBlack, BrightRed, BrightGreen, BrightYellow,
BrightBlue, BrightMagenta, BrightCyan, BrightWhite
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

#### Available attributes

```
Bold, Faint, Italic, Underline, Blink, Reverse, Strikethrough
```

Styles can also apply background colors:

```go
highlight := brush.Style{
	Foreground: brush.White,
	Background: brush.Blue,
}

fmt.Println(highlight.Paint("highlighted text"))
```

### Stencil

Color uppercase letters, lowercase letters, digits, punctuation, and symbols each with their own color.
Anything outside these five classes, such as whitespace, is left unstyled.

```go
st := brush.Stencil{
	Upper:  brush.Red,
	Lower:  brush.Green,
	Digit:  brush.Blue,
	Punct:  brush.Yellow,
	Symbol: brush.Yellow,
}

fmt.Println(st.Paint("XYZ{F4K3_FL4G}"))
```

The zero value leaves everything unstyled, so set only the fields you care about.

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

## License

This is an educational repository under the [MIT license](LICENSE).

// Package colorize contains functions to stylize texts for terminal.
package colorize

import (
	"fmt"
)

type style int

// Colors
const (
	ColorBlack style = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Bright Colors
const (
	ColorBlackBright style = iota + 90
	ColorRedBright
	ColorGreenBright
	ColorYellowBright
	ColorBlueBright
	ColorMagentaBright
	ColorCyanBright
	ColorWhiteBright
)

// Text Styles
const (
	StyleBold style = iota + 1
	StyleDim
	StyleItalic
	StyleUnderline
)

// Configurers Text Styles
const (
	StyleInverse style = iota + 7
	StyleHidden
	StyleStrikethrough
)

type Colorize struct {
	Enabled bool
}

func (c *Colorize) Black(s string) string         { return c.stylize(s, ColorBlack, 39) }
func (c *Colorize) BlackBright(s string) string   { return c.stylize(s, ColorBlackBright, 39) }
func (c *Colorize) Red(s string) string           { return c.stylize(s, ColorRed, 39) }
func (c *Colorize) RedBright(s string) string     { return c.stylize(s, ColorRedBright, 39) }
func (c *Colorize) Green(s string) string         { return c.stylize(s, ColorGreen, 39) }
func (c *Colorize) GreenBright(s string) string   { return c.stylize(s, ColorGreenBright, 39) }
func (c *Colorize) Yellow(s string) string        { return c.stylize(s, ColorYellow, 39) }
func (c *Colorize) YellowBright(s string) string  { return c.stylize(s, ColorYellowBright, 39) }
func (c *Colorize) Blue(s string) string          { return c.stylize(s, ColorBlue, 39) }
func (c *Colorize) BlueBright(s string) string    { return c.stylize(s, ColorBlueBright, 39) }
func (c *Colorize) Magenta(s string) string       { return c.stylize(s, ColorMagenta, 39) }
func (c *Colorize) MagentaBright(s string) string { return c.stylize(s, ColorMagentaBright, 39) }
func (c *Colorize) Cyan(s string) string          { return c.stylize(s, ColorCyan, 39) }
func (c *Colorize) CyanBright(s string) string    { return c.stylize(s, ColorCyanBright, 39) }
func (c *Colorize) Gray(s string) string          { return c.stylize(s, ColorBlackBright, 39) }
func (c *Colorize) White(s string) string         { return c.stylize(s, ColorWhite, 39) }
func (c *Colorize) WhiteBright(s string) string   { return c.stylize(s, ColorWhiteBright, 39) }

func (c *Colorize) Bold(s string) string          { return c.stylize(s, StyleBold, 22) }
func (c *Colorize) Dim(s string) string           { return c.stylize(s, StyleDim, 22) }
func (c *Colorize) Italic(s string) string        { return c.stylize(s, StyleItalic, 23) }
func (c *Colorize) Underline(s string) string     { return c.stylize(s, StyleUnderline, 24) }
func (c *Colorize) Inverse(s string) string       { return c.stylize(s, StyleInverse, 27) }
func (c *Colorize) Hidden(s string) string        { return c.stylize(s, StyleHidden, 28) }
func (c *Colorize) Strikethrough(s string) string { return c.stylize(s, StyleStrikethrough, 29) }

func (c *Colorize) UseColors(value bool) {
	c.Enabled = value
}

func (c *Colorize) stylize(s string, open style, close int) string {
	if c.Enabled {
		return fmt.Sprintf("\x1b[%dm%s\x1b[%dm", open, s, close)
	}

	return s
}

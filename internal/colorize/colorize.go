// Package colorize contains functions to stylize texts for terminal.
package colorize

import (
	"fmt"
	"os"
)

const _colorEnv = "MOCHA_COLOR"

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

// Additional Text Styles
const (
	StyleInverse style = iota + 7
	StyleHidden
	StyleStrikethrough
)

var (
	color, useColor = os.LookupEnv(_colorEnv)
)

func Black(s string) string         { return stylize(s, ColorBlack, 39) }
func BlackBright(s string) string   { return stylize(s, ColorBlackBright, 39) }
func Red(s string) string           { return stylize(s, ColorRed, 39) }
func RedBright(s string) string     { return stylize(s, ColorRedBright, 39) }
func Green(s string) string         { return stylize(s, ColorGreen, 39) }
func GreenBright(s string) string   { return stylize(s, ColorGreenBright, 39) }
func Yellow(s string) string        { return stylize(s, ColorYellow, 39) }
func YellowBright(s string) string  { return stylize(s, ColorYellowBright, 39) }
func Blue(s string) string          { return stylize(s, ColorBlue, 39) }
func BlueBright(s string) string    { return stylize(s, ColorBlueBright, 39) }
func Magenta(s string) string       { return stylize(s, ColorMagenta, 39) }
func MagentaBright(s string) string { return stylize(s, ColorMagentaBright, 39) }
func Cyan(s string) string          { return stylize(s, ColorCyan, 39) }
func CyanBright(s string) string    { return stylize(s, ColorCyanBright, 39) }
func Gray(s string) string          { return stylize(s, ColorBlackBright, 39) }
func White(s string) string         { return stylize(s, ColorWhite, 39) }
func WhiteBright(s string) string   { return stylize(s, ColorWhiteBright, 39) }

func Bold(s string) string          { return stylize(s, StyleBold, 22) }
func Dim(s string) string           { return stylize(s, StyleDim, 22) }
func Italic(s string) string        { return stylize(s, StyleItalic, 23) }
func Underline(s string) string     { return stylize(s, StyleUnderline, 24) }
func Inverse(s string) string       { return stylize(s, StyleInverse, 27) }
func Hidden(s string) string        { return stylize(s, StyleHidden, 28) }
func Strikethrough(s string) string { return stylize(s, StyleStrikethrough, 29) }

func stylize(s string, open style, close int) string {
	if !useColor {
		return s
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[%dm", open, s, close)
}

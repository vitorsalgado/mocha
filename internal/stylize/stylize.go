// Package stylize contains functions to colorize and stylize texts for terminal
package stylize

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	Reset = "\u001b[0m"

	Bright       = 1
	ColorBlack   = 30
	ColorRed     = 31
	ColorGreen   = 32
	ColorYellow  = 33
	ColorBlue    = 34
	ColorMagenta = 35
	ColorCyan    = 36
	ColorWhite   = 97
	ColorGray    = 90

	StyleBold          = 1
	StyleDim           = 2
	StyleItalic        = 3
	StyleUnderline     = 4
	StyleInverse       = 7
	StyleHidden        = 8
	StyleStrikethrough = 9
)

var (
	noColor = os.Getenv("MOCHA_NO_COLOR") == "true"
	isWin   = runtime.GOOS == "windows"
	_, isCI = os.LookupEnv("CI")
)

// Stylize stylizes the text using provided text style parameters
func Stylize(s string, codes ...int) string {
	if noColor || isWin || isCI {
		return s
	}

	c := make([]string, len(codes))
	for i, code := range codes {
		c[i] = strconv.Itoa(code)
	}

	return fmt.Sprintf("\u001b[%sm%s%s", strings.Join(c, ";"), s, Reset)
}

func Black(s string) string         { return Stylize(s, ColorBlack) }
func BlackBright(s string) string   { return Stylize(s, ColorBlack) }
func Red(s string) string           { return Stylize(s, ColorRed) }
func RedBright(s string) string     { return Stylize(s, ColorRed, Bright) }
func Green(s string) string         { return Stylize(s, ColorGreen) }
func GreenBright(s string) string   { return Stylize(s, ColorGreen, Bright) }
func Yellow(s string) string        { return Stylize(s, ColorYellow) }
func YellowBright(s string) string  { return Stylize(s, ColorYellow, Bright) }
func Blue(s string) string          { return Stylize(s, ColorBlue) }
func BlueBright(s string) string    { return Stylize(s, ColorBlue, Bright) }
func Magenta(s string) string       { return Stylize(s, ColorMagenta) }
func MagentaBright(s string) string { return Stylize(s, ColorMagenta, Bright) }
func Cyan(s string) string          { return Stylize(s, ColorCyan) }
func CyanBright(s string) string    { return Stylize(s, ColorCyan, Bright) }
func Gray(s string) string          { return Stylize(s, ColorGray) }
func GrayBright(s string) string    { return Stylize(s, ColorGray, Bright) }
func White(s string) string         { return Stylize(s, ColorWhite) }
func WhiteBright(s string) string   { return Stylize(s, ColorWhite, Bright) }

func Bold(s string) string          { return Stylize(s, StyleBold) }
func Dim(s string) string           { return Stylize(s, StyleDim) }
func Italic(s string) string        { return Stylize(s, StyleItalic) }
func Underline(s string) string     { return Stylize(s, StyleUnderline) }
func Inverse(s string) string       { return Stylize(s, StyleInverse) }
func Hidden(s string) string        { return Stylize(s, StyleHidden) }
func Strikethrough(s string) string { return Stylize(s, StyleStrikethrough) }

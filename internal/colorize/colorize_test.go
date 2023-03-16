package colorize

import (
	"fmt"
	"testing"
)

func TestStylize(t *testing.T) {
	c := Colorize{}
	c.UseColors(true)

	fmt.Print(c.Black("black"))
	fmt.Print(c.BlackBright("black bright"))
	fmt.Print(c.Red("red"))
	fmt.Print(c.RedBright("red bright"))
	fmt.Print(c.Green("green"))
	fmt.Print(c.GreenBright("green bright"))
	fmt.Print(c.Yellow("yellow"))
	fmt.Print(c.YellowBright("yellow bright"))
	fmt.Print(c.Blue("blue"))
	fmt.Print(c.BlueBright("blue bright"))
	fmt.Print(c.Magenta("magenta"))
	fmt.Print(c.MagentaBright("magenta bright"))
	fmt.Print(c.Cyan("cyan"))
	fmt.Print(c.CyanBright("cyan bright"))
	fmt.Print(c.Gray("gray"))
	fmt.Print(c.White("white"))
	fmt.Print(c.WhiteBright("white bright"))

	fmt.Print(c.Bold("bold"))
	fmt.Print(c.Dim("dim"))
	fmt.Print(c.Italic("italic"))
	fmt.Print(c.Underline("underline"))
	fmt.Print(c.Inverse("inverse"))
	fmt.Print(c.Hidden("hidden"))
	fmt.Print(c.Strikethrough("strikethrough"))
}

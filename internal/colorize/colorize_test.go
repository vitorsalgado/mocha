package colorize

import (
	"fmt"
	"testing"
)

func TestStylize(t *testing.T) {
	c := Colorize{}
	c.UseColors(true)

	fmt.Println(c.Black("black"))
	fmt.Println(c.BlackBright("black bright"))
	fmt.Println(c.Red("red"))
	fmt.Println(c.RedBright("red bright"))
	fmt.Println(c.Green("green"))
	fmt.Println(c.GreenBright("green bright"))
	fmt.Println(c.Yellow("yellow"))
	fmt.Println(c.YellowBright("yellow bright"))
	fmt.Println(c.Blue("blue"))
	fmt.Println(c.BlueBright("blue bright"))
	fmt.Println(c.Magenta("magenta"))
	fmt.Println(c.MagentaBright("magenta bright"))
	fmt.Println(c.Cyan("cyan"))
	fmt.Println(c.CyanBright("cyan bright"))
	fmt.Println(c.Gray("gray"))
	fmt.Println(c.White("white"))
	fmt.Println(c.WhiteBright("white bright"))

	fmt.Println(c.Bold("bold"))
	fmt.Println(c.Dim("dim"))
	fmt.Println(c.Italic("italic"))
	fmt.Println(c.Underline("underline"))
	fmt.Println(c.Inverse("inverse"))
	fmt.Println(c.Hidden("hidden"))
	fmt.Println(c.Strikethrough("strikethrough"))
}

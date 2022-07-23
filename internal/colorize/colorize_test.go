package colorize

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv(_noColorEnv, "1")

	code := m.Run()

	_ = os.Unsetenv(_noColorEnv)

	os.Exit(code)
}

func TestStylize(t *testing.T) {
	fmt.Println(Black("black"))
	fmt.Println(BlackBright("black bright"))
	fmt.Println(Red("red"))
	fmt.Println(RedBright("red bright"))
	fmt.Println(Green("green"))
	fmt.Println(GreenBright("green bright"))
	fmt.Println(Yellow("yellow"))
	fmt.Println(YellowBright("yellow bright"))
	fmt.Println(Blue("blue"))
	fmt.Println(BlueBright("blue bright"))
	fmt.Println(Magenta("magenta"))
	fmt.Println(MagentaBright("magenta bright"))
	fmt.Println(Cyan("cyan"))
	fmt.Println(CyanBright("cyan bright"))
	fmt.Println(Gray("gray"))
	fmt.Println(White("white"))
	fmt.Println(WhiteBright("white bright"))

	fmt.Println(Bold("bold"))
	fmt.Println(Dim("dim"))
	fmt.Println(Italic("italic"))
	fmt.Println(Underline("underline"))
	fmt.Println(Inverse("inverse"))
	fmt.Println(Hidden("hidden"))
	fmt.Println(Strikethrough("strikethrough"))
}

func TestMultipleStyles(t *testing.T) {
	fmt.Println(Red(Bold("hello") + " world"))
}

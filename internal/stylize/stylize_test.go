package stylize

import (
	"fmt"
	"testing"
)

func TestStylize(t *testing.T) {
	msg := "hello world"

	fmt.Println(Black(msg))
	fmt.Println(BlackBright(msg))
	fmt.Println(Red(msg))
	fmt.Println(RedBright(msg))
	fmt.Println(Green(msg))
	fmt.Println(GreenBright(msg))
	fmt.Println(Yellow(msg))
	fmt.Println(YellowBright(msg))
	fmt.Println(Blue(msg))
	fmt.Println(BlueBright(msg))
	fmt.Println(Magenta(msg))
	fmt.Println(MagentaBright(msg))
	fmt.Println(Cyan(msg))
	fmt.Println(CyanBright(msg))
	fmt.Println(Gray(msg))
	fmt.Println(GrayBright(msg))
	fmt.Println(White(msg))
	fmt.Println(WhiteBright(msg))

	fmt.Println(Bold(msg))
	fmt.Println(Dim(msg))
	fmt.Println(Italic(msg))
	fmt.Println(Underline(msg))
	fmt.Println(Inverse(msg))
	fmt.Println(Hidden(msg))
	fmt.Println(Strikethrough(msg))
}

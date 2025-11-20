package commands

import (
	"github.com/fatih/color"
)

func ColorIt(str string,tcolor string) string {
	switch tcolor {
		case "green":
			green := color.New(color.FgGreen).SprintFunc()
			result := green(str)
			return result
		case "red":
			green := color.New(color.FgRed).SprintFunc()
			result := green(str)
			return result
		case "blue":
			green := color.New(color.FgBlue).SprintFunc()
			result := green(str)
			return result
		default:
			return tcolor
		}

}
func ResultColor(ok bool) string {
	if ok {
		return  ColorIt("[ok]" , "green")
	} else {
		return ColorIt("[X]" , "red")
	}
}
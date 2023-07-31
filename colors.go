package main

import "fmt"

type Color int
type Effect int
type Background int

const (
	Red Color = iota
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	Underline Effect = iota
	Strikethrough
	Bold
)

const (
	BgRed Background = iota
	BgGreen
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// TODO: Implement CSS (hex) VALUE TO ANSI COLOR CODE
// - Will need to grab the hex value from here:
// div.message-convo-left:nth-child(n) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > strong:nth-child(1)
func hexToANSI(hex string) string {
	return "" // TODO
}

// maps to translate enums to ASCII color codes
var colors = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
var effects = []string{"\033[4;31m", "\033[9;31m", "\033[1;31m"}
var backgrounds = []string{"\033[41m", "\033[42m", "\033[44m", "\033[45m", "\033[46m", "\033[47m"}
var reset = "\033[0m"

func asciiColor(color Color, effect *Effect, background *Background) string {
	var effectString, backgroundString string
	if effect != nil {
		effectString = effects[*effect]
	}
	if background != nil {
		backgroundString = backgrounds[*background]
	}
	return fmt.Sprintf("%v%v%v%v", effectString, colors[color], backgroundString, reset)
}

func printColoredUser(username string, usergroup int) string {
	var color Color
	// var effect *Effect
	// var background *Background
	switch usergroup {
	case 2:
		color = White
	case 3:
		color = Cyan
	case 4:
		color = Magenta
	case 7:
		color = Red
	case 9:
		color = Yellow
	case 28:
		color = Blue
	case 38:
		color = White
		// e := Strikethrough
		// effect = &e
		// b := BgRed
		// background = &b
	default:
		color = White
	}
	return fmt.Sprintf("%v%v%v", colors[color], username, reset)
}

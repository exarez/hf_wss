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

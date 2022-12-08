package day00

import (
	_ "embed"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput any
}

func BuildDay00() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func (d *day) parse() {}

func (d *day) Run1() string {
	return ""
}

func (d *day) Run2() string {
	return ""
}

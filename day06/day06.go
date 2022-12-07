package day06

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput []rune
}

func BuildDay06() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func (d *day) parse() {
	trimmed := strings.TrimSpace(input)
	d.parsedInput = []rune(trimmed)
}

type set map[rune]*any

func newSet(elems ...rune) set {
	s := set{}
	for _, e := range elems {
		s[e] = nil
	}
	return s
}

func (d *day) detect(n int) int {
	var i int
	for i = n - 1; i < len(d.parsedInput); i++ {
		if len(newSet(d.parsedInput[i-n+1:i+1]...)) == n {
			return i + 1
		}
	}
	panic("pattern not found")
}

func (d *day) Run1() string {
	return fmt.Sprint(d.detect(4))
}

func (d *day) Run2() string {
	return fmt.Sprint(d.detect(14))
}

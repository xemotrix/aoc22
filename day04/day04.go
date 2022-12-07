package day04

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput [][2]section
}

func BuildDay04() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

type section [2]int

func (d *day) parse() {
	splittedLines := strings.Split(strings.TrimSpace(input), "\n")

	splitTwice := Pipe(FCurry(strings.Split, ","), Curry(Map[string, []string], FCurry(strings.Split, "-")))
	splitAndInt := Pipe(splitTwice, Curry(Map[[]string, []int], Curry(Map[string, int], Expect(strconv.Atoi))))

	nestedSlices := Map(splitAndInt, splittedLines)

	toSection := func(s []int) section {
		if len(s) != 2 {
			panic("incorrect len for section")
		}
		return section{s[0], s[1]}
	}

	sliceToArr := func(s []section) [2]section {
		if len(s) != 2 {
			panic("incorrect len for slice of section")
		}
		return [2]section{s[0], s[1]}
	}

	toFinalType := Pipe(Curry(Map[[][]int, []section], Curry(Map[[]int, section], toSection)), Curry(Map[[]section, [2]section], sliceToArr))

	d.parsedInput = toFinalType(nestedSlices)
}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func (d *day) Run1() string {
	checkOverlap := func(a [2]section) bool {
		return (a[0][0] >= a[1][0] && a[0][1] <= a[1][1]) || (a[1][0] >= a[0][0] && a[1][1] <= a[0][1])
	}
	res := len(Filter(Identity[bool], Map(checkOverlap, d.parsedInput)))

	return fmt.Sprint(res)
}

func (d *day) Run2() string {
	checkOverlap := func(a [2]section) bool {
		return (a[0][0] >= a[1][0] && a[0][0] <= a[1][1]) ||
			(a[0][1] >= a[1][0] && a[0][1] <= a[1][1]) ||
			(a[1][0] >= a[0][0] && a[1][0] <= a[0][1]) ||
			(a[1][1] >= a[0][0] && a[1][1] <= a[0][1])
	}
	res := len(Filter(Identity[bool], Map(checkOverlap, d.parsedInput)))

	return fmt.Sprint(res)
}

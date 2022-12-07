package day03

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput [][][]int
}

func BuildDay03() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) parse() {

	splittedLines := strings.Split(strings.TrimSpace(input), "\n")

	runes := Map(func(s string) []rune { return []rune(s) }, splittedLines)

	adjust := func(r rune) int {
		if r >= 97 && r <= 122 {
			return int(r) - 96
		} else if r >= 65 && r <= 90 {
			return int(r) - 65 + 27
		}
		panic(fmt.Sprintf("not valid letter: '%v'", string(r)))
	}

	adjusted := Map(Curry(Map[rune, int], adjust), runes)

	splitHalf := func(s []int) [][]int {
		splitPoint := len(s) / 2
		return [][]int{s[:splitPoint], s[splitPoint:]}
	}

	spittedRunes := Map(splitHalf, adjusted)

	d.parsedInput = spittedRunes
}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

var sumReduce func([]int) int = Curry(Reduce[int], func(a, b int) int { return a + b })

type set map[int]*any

func toSet(s []int) set {
	mySet := set{}
	Map(func(i int) *any { mySet[i] = nil; return nil }, s)
	return mySet
}

func toSlice(m map[int]*any) []int {
	s := []int{}
	for k := range m {
		s = append(s, k)
	}
	return s
}

func checkInter(a []int, b []int) []int {
	set1 := toSet(a)
	inter := Filter(func(i int) bool { _, ok := set1[i]; return ok }, b)
	res := toSlice(toSet(inter))
	return res
}

func (d *day) Run1() string {

	takeFirst := func(s []int) int {
		if len(s) > 0 {
			return s[0]
		}
		panic("cant take first from empty list")
	}
	intersections := Map(Curry(Reduce[[]int], checkInter), d.parsedInput)
	res := sumReduce(Map(takeFirst, intersections))

	return fmt.Sprint(res)
}

func (d *day) Run2() string {
	concat := func(s [][]int) []int {
		return append(s[0], s[1]...)
	}

	concated := Map(concat, d.parsedInput)

	pivoted := [][][]int{}
	for i := 0; i < len(concated); i += 3 {
		pivoted = append(pivoted, concated[i:i+3])
	}

	intersections := Map(Curry(Reduce[[]int], checkInter), pivoted)
	score := sumReduce(Map(sumReduce, intersections))

	return fmt.Sprint(score)
}

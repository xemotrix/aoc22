package day01

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// // go:embed input2.txt
// var input2 string

type day struct {
	input       string
	parsedInput [][]int
}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func BuildDay01() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) parse() {
	arrStr := Map(
		Curry(Flip(strings.Split), "\n"),
		Map(strings.TrimSpace, strings.Split(d.input, "\n\n")),
	)

	strToInt := Curry(Map[string, int], Expect(strconv.Atoi))
	d.parsedInput = Map(strToInt, arrStr)
}

func (d *day) Run1() string {
	elfTotals := Map(
		Curry(Reduce[int], func(a, b int) int { return a + b }),
		d.parsedInput,
	)
	return fmt.Sprint(Reduce(Max[int], elfTotals))
}

func (d *day) Run2() string {

	elfTotals := Map(
		Curry(Reduce[int], func(a, b int) int { return a + b }),
		d.parsedInput,
	)
	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })

	initialSum := sumReduce(elfTotals)

	var popMax func([]int, int) []int
	popMax = func(in []int, times int) []int {
		max := Reduce(Max[int] , in, )
		elfTotals = Filter(func(v int) bool { return v != max },elfTotals, )
		if times == 1 {
			return elfTotals
		}
		return popMax(elfTotals, times-1)
	}

	result := initialSum - sumReduce(popMax(elfTotals, 3))
	return fmt.Sprint(result)

}

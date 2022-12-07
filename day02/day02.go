package day02

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
	parsedInput [][2]int
}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func BuildDay02() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d
}

const (
	Rock int = iota + 1
	Paper
	Scissors
)

func (d *day) parse() {
	toHand := func(s string) [2]int {
		var nameMap = map[rune]int{
			'X': Rock,
			'Y': Paper,
			'Z': Scissors,
			'A': Rock,
			'B': Paper,
			'C': Scissors,
		}
		if len(s) != 3 {
			panic("Error parsing, length of hand != 1")
		}
		runes := []rune(s)
		return [2]int{nameMap[runes[0]], nameMap[runes[2]]}
	}
	splitted := strings.Split(strings.TrimSpace(input), "\n")
	d.parsedInput = Map(toHand, splitted)
}

func checkMatchScore(pair [2]int) int {
	score := pair[1]
	if pair[1] == pair[0] {
		score += 3
	} else if pair[1] == (pair[0]%3)+1 {
		score += 6
	}
	return score
}

func (d *day) Run1() string {
	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })
	res := sumReduce(Map(checkMatchScore, d.parsedInput))
	return fmt.Sprint(res)
}

func applyStrat(pair [2]int) (myChoice int) {
	switch pair[1] {
	case 1: // lose
		myChoice = -((-pair[0]+4)%3 + 1) + 4
	case 2: // draw
		myChoice = pair[0]
	case 3: // win
		myChoice = (pair[0] % 3) + 1
	default:
		panic("invalid choice")
	}
	return
}

func (d *day) Run2() string {
	Map(applyStrat, d.parsedInput)

	return ""
}

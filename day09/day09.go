package day09

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput []instruction
}

func BuildDay09() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type instruction struct {
	dir rune
	num int
}

func (d *day) parse() {

	instructionFromStr := func(s string) instruction {
		if len(s) < 3 {
			panic("invalid instruction string")
		}
		return instruction{
			dir: rune(s[0]),
			num: Expect(strconv.Atoi)(s[2:]),
		}
	}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	d.parsedInput = Map(instructionFromStr, lines)
}

type loc struct {
	x, y int
}
type ropeLink struct {
	H loc
	T loc
}

func (r *ropeLink) checkValid() bool {
	distanceSq := math.Pow(float64(r.H.x-r.T.x), 2) + math.Pow(float64(r.H.y-r.T.y), 2)
	return distanceSq <= 2
}

func (r *ropeLink) move(ins instruction) []loc {
	tPositions := make([]loc, ins.num)
	for i := 0; i < ins.num; i++ {
		switch ins.dir {
		case 'R':
			r.H.x++
		case 'L':
			r.H.x--
		case 'U':
			r.H.y++
		case 'D':
			r.H.y--
		}
		if valid := r.checkValid(); !valid {
			switch ins.dir {
			case 'R':
				r.T = loc{
					x: r.H.x - 1,
					y: r.H.y,
				}
			case 'L':
				r.T = loc{
					x: r.H.x + 1,
					y: r.H.y,
				}
			case 'U':
				r.T = loc{
					x: r.H.x,
					y: r.H.y - 1,
				}
			case 'D':
				r.T = loc{
					x: r.H.x,
					y: r.H.y + 1,
				}
			}
		}
		tPositions[i] = r.T
	}
	return tPositions
}

type set[T comparable] map[T]*any

func (s set[T]) add(elems ...T) {
	for _, e := range elems {
		s[e] = nil
	}
}

func (d *day) Run1() string {
	r := ropeLink{}
	tailPositions := set[loc]{}
	for i, ins := range d.parsedInput {
		tailPositions.add(r.move(ins)...)
		fmt.Printf("Iter %d\n", i)
		fmt.Printf(" - H: xy(%d, %d)\n", r.H.x, r.H.y)
		fmt.Printf(" - T: xy(%d, %d)\n", r.T.x, r.T.y)
	}
	fmt.Println(len(tailPositions))
	return ""
}

type ropeKnot struct {
	loc  loc
	next *ropeKnot
}

func (d *day) Run2() string {
	return ""
}

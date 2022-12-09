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

type knot struct {
	loc  loc
	next *knot
}

func (r *knot) checkNext() {
	if r.next == nil {
		return
	}

	distX := r.loc.x - r.next.loc.x
	distY := r.loc.y - r.next.loc.y

	distanceSq := math.Pow(float64(distX), 2) + math.Pow(float64(distY), 2)
	if distanceSq <= 2 {
		return
	}

	r.next.loc = r.loc

	if distX >= 2 { // we are moving right
		r.next.loc.x--
	} else if distX <= -2 { // we are moving left
		r.next.loc.x++
	}

	if distY >= 2 { // we are moving up
		r.next.loc.y--
	} else if distY <= -2 { // we are moving down
		r.next.loc.y++
	}

	r.next.checkNext()
}

func (k *knot) step(dir rune) {
	switch dir {
	case 'R':
		k.loc.x++
	case 'L':
		k.loc.x--
	case 'U':
		k.loc.y++
	case 'D':
		k.loc.y--
	}
	k.checkNext()
	return
}

func (k *knot) getTailLoc() loc {
	if k.next != nil {
		return k.next.getTailLoc()
	}
	return k.loc
}

func (k *knot) move(ins instruction) (tPositions []loc) {
	for i := 0; i < ins.num; i++ {
		k.step(ins.dir)
		tPositions = append(tPositions, k.getTailLoc())
	}
	return tPositions
}

func buildRope(size int) *knot {
	if size == 0 {
		return nil
	}
	return &knot{
		next: buildRope(size - 1),
	}
}

type set[T comparable] map[T]*any

func (s set[T]) add(elems ...T) {
	for _, e := range elems {
		s[e] = nil
	}
}

func (d *day) runSimulation(ropeLength int) int {
	head := buildRope(ropeLength)
	tailPositions := set[loc]{}
	for _, ins := range d.parsedInput {
		tailPositions.add(head.move(ins)...)
	}
	return len(tailPositions)
}

func (d *day) Run1() string {
	return fmt.Sprint(d.runSimulation(2))
}

func (d *day) Run2() string {
	return fmt.Sprint(d.runSimulation(10))
}

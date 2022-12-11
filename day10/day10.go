package day10

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:embed input.txt
var input string

type day struct {
	input    string
	pnode    *progNode
	register int
	cycle    int
}

type progNode struct {
	instr    instruction
	waitTime int
	next     *progNode
}

type instruction struct {
	name  string
	value int
	setup int
}

func BuildDay10() *day {
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
	lines := strings.Split(strings.TrimSpace(input), "\n")
	parseInstr := func(s string) instruction {
		parts := strings.Split(s, " ")
		instr := instruction{
			name: parts[0],
		}
		if len(parts) > 1 {
			val := Expect(strconv.Atoi)(parts[1])
			instr.value = val
		}
		return instr
	}

	prev := &progNode{instr: parseInstr(lines[0])}

	head := prev

	for _, l := range lines[1:] {
		node := progNode{instr: parseInstr(l)}
		prev.next = &node
		prev = &node
	}

	d.register = 1
	d.cycle = 1
	d.pnode = head
}

func (d *day) runCycle() {
	switch d.pnode.instr.name {
	case "noop":
		d.pnode = d.pnode.next
	case "addx":
		d.pnode.instr.setup++
		if d.pnode.instr.setup == 2 {
			d.register += d.pnode.instr.value
			d.pnode = d.pnode.next
		}
	}
	d.cycle++
}

func (d *day) Run1() string {
	resSlice := []int{}

	for {
		if d.pnode == nil {
			break
		}
		if d.cycle == 20 || (d.cycle-20)%40 == 0 {
			resSlice = append(resSlice, d.register*d.cycle)
		}
		d.runCycle()
	}

	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })
	res := sumReduce(resSlice)

	return fmt.Sprint(res)
}

const (
	Reset = "\033[0m"
	Green = "\u001b[46m"
	Black = "\u001b[30m"
)

func (d *day) render(s []bool) string {
	res := ""
	for i, pix := range s {

		h := i % 40

		if h == 0 {
			res += "\n"
		}

		pixSym := ""
		if pix {
			pixSym = "#"
		} else {
			pixSym = "·"
		}

		if h >= d.register-1 && h <= d.register-1+2 {
			pixSym = Green + Black + pixSym + Reset
		}

		res += pixSym
	}
	return res
}

func (d *day) Run2() string {
	d.parse()

	screen := make([]bool, 40*6)

	for {
		if d.pnode == nil {
			break
		}
		time.Sleep(15 * time.Millisecond)

		drawingPos := (d.cycle - 1) % 40

		if drawingPos >= d.register-1 && drawingPos <= d.register-1+2 {
			screen[d.cycle] = true
		}

		d.runCycle()

		fmt.Println(d.render(screen))
	}

	return d.render(screen)
}

package day05

import (
	_ "embed"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input        string
	queues       []queue
	instructions []instruction
}

func BuildDay05() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type node struct {
	data rune
	next *node
}

type queue struct {
	top *node
}

func (q *queue) push(elements ...rune) {
	for _, ele := range elements {
		n := node{
			data: ele,
			next: q.top,
		}
		q.top = &n
	}
}

func (q *queue) pop() (rune, bool) {
	if q.top == nil {
		return 0, false
	}
	res := q.top.data
	q.top = q.top.next // GC do your stuff
	return res, true
}

type instruction struct {
	quantity int
	from     int
	to       int
}

func (d *day) parse() {
	rawSplitted := strings.Split(input, "\n\n")

	// firs part: stacks
	stacks := strings.Split(rawSplitted[0], "\n")

	var nStacksRune rune
	for _, char := range []rune(stacks[len(stacks)-1]) {
		if char != ' ' {
			nStacksRune = char
		}
	}

	numQueues, err := strconv.Atoi(string(nStacksRune))
	if err != nil {
		panic("Error parsing cant convert to number")
	}

	queues := make([]queue, numQueues)
	for i := len(stacks) - 2; i >= 0; i-- {

		for x, char := range []rune(stacks[i]) {
			if char >= 65 && char <= 90 {
				queues[(x-1)/4].push(char)
			}
		}
	}
	d.queues = queues

	// second part: instructions
	rawInstr := strings.Split(strings.TrimSpace(rawSplitted[1]), "\n")

	d.instructions = make([]instruction, len(rawInstr))
	for i, line := range rawInstr {
		regex := regexp.MustCompile(`\d+`)
		matches := regex.FindAllString(line, -1)

		quantity, _ := strconv.Atoi(matches[0])
		from, _ := strconv.Atoi(matches[1])
		to, _ := strconv.Atoi(matches[2])

		d.instructions[i] = instruction{
			quantity: quantity,
			from:     from,
			to:       to,
		}
	}
}

func (d *day) Run1() string {
	for _, ins := range d.instructions {
		for i := 0; i < ins.quantity; i++ {
			val, _ := d.queues[ins.from-1].pop()
			d.queues[ins.to-1].push(val)
		}
	}
	res := make([]rune, len(d.queues))
	for i, q := range d.queues {
		res[i], _ = q.pop()
	}
	return string(res)
}

func (d *day) Run2() string {
	d.parse()

	for _, ins := range d.instructions {
		toMove := make([]rune, ins.quantity)
		for i := ins.quantity - 1; i >= 0; i-- {
			toMove[i], _ = d.queues[ins.from-1].pop()
		}
		d.queues[ins.to-1].push(toMove...)
	}
	res := make([]rune, len(d.queues))
	for i, q := range d.queues {
		res[i], _ = q.pop()
	}
	return string(res)
}

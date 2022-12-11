package day11

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input     string
	monkeys   []monkey
	commonMod int
}

func BuildDay11() *day {
	d := day{
		input: input,
	}
	d.parse(1)
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type monkey struct {
	items        Dequeue[int]
	operation    func(int) int
	test         func(int) int
	inspectCount int
}

func (d *day) parse(part int) {
	monkeyStr := strings.Split(strings.TrimSpace(input), "\n\n")
	numPattern := regexp.MustCompile(`\d+`)
	opPattern := regexp.MustCompile(`\+|\*`)

	d.commonMod = 1

	parseMonkey := func(s string) monkey {
		lines := strings.Split(s, "\n")
		m := monkey{}

		startItems := Map(Expect(strconv.Atoi), numPattern.FindAllString(lines[1], -1))
		m.items = BuildDequeue(startItems...)

		op := opPattern.FindString(lines[2])

		opFun := func(n int) (res int) {
			intNum, err := strconv.Atoi(numPattern.FindString(lines[2]))
			if err != nil {
				intNum = n
			}
			opNum := intNum
			switch op {
			case "+":
				res = n + opNum
			case "*":
				res = n * opNum
			default:
				panic(fmt.Sprintf("Invalid operation %s", op))
			}
			if part == 1 {
				return res / 3
			}

			return res % d.commonMod
		}
		m.operation = opFun

		testDiv := Expect(strconv.Atoi)(numPattern.FindString(lines[3]))
		d.commonMod *= testDiv
		testTrue := Expect(strconv.Atoi)(numPattern.FindString(lines[4]))
		testFalse := Expect(strconv.Atoi)(numPattern.FindString(lines[5]))
		testFun := func(n int) int {
			if n%testDiv == 0 {
				return testTrue
			}
			return testFalse
		}
		m.test = testFun

		return m

	}

	d.monkeys = Map(parseMonkey, monkeyStr)
}

func (d *day) printMonkeys() {
	for i, m := range d.monkeys {
		fmt.Printf("Monkey %d: ", i)
		fmt.Printf("(inspected %d) current items -> ", m.inspectCount)
		printItem := func(item int) int {
			fmt.Printf("%d, ", item)
			return item
		}
		m.items.ApplyToAll(printItem)
		fmt.Println()
	}
	fmt.Println()
}

func (m *monkey) doItem() (worry, nextMonkeyIndex int) {
	item := m.items.PopHead()
	worry = m.operation(item)
	nextMonkeyIndex = m.test(worry)
	m.inspectCount++
	return
}

func (d *day) runSim(iters int) string {
	for i := 0; i < iters; i++ {
		for mi := range d.monkeys {
			m := &d.monkeys[mi]
			for !m.items.Empty() {
				worry, nextMonkeyIndex := m.doItem()
				d.monkeys[nextMonkeyIndex].items.PushTail(worry)
			}
		}
	}
	getInspectCount := func(m monkey) int { return m.inspectCount }
	counts := Map(getInspectCount, d.monkeys)
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	return fmt.Sprint(counts[0] * counts[1])

}

func (d *day) Run1() string {
	return d.runSim(20)
}

func (d *day) Run2() string {
	d.parse(2)
	return d.runSim(10000)
}

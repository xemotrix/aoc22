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
	input   string
	monkeys []monkey
}

func BuildDay11() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type monkey struct {
	items        Dequeue[uint64]
	operation    func(uint64) uint64
	test         func(uint64) uint64
	inspectCount int
}

func (d *day) parse() {
	monkeyStr := strings.Split(strings.TrimSpace(input), "\n\n")
	numPattern := regexp.MustCompile(`\d+`)
	opPattern := regexp.MustCompile(`\+|\*`)
	parseMonkey := func(s string) monkey {
		lines := strings.Split(s, "\n")
		m := monkey{}

		startItems := Map(Expect(strconv.Atoi), numPattern.FindAllString(lines[1], -1))
		toInt64 := func(num int) uint64 { return uint64(num) }
		m.items = BuildDequeue(Map(toInt64, startItems)...)

		op := opPattern.FindString(lines[2])

		opFun := func(n uint64) (res uint64) {
			intNum, err := strconv.Atoi(numPattern.FindString(lines[2]))
			if err != nil {
				intNum = int(n)
			}
			opNum := uint64(intNum)
			switch op {
			case "+":
				res = n + uint64(opNum)
			case "*":
				res = n * uint64(opNum)
			default:
				panic(fmt.Sprintf("Invalid operation %s", op))
			}
			worryFinal := res / 3
			return worryFinal
		}
		m.operation = opFun

		testDiv := uint64(Expect(strconv.Atoi)(numPattern.FindString(lines[3])))
		testTrue := uint64(Expect(strconv.Atoi)(numPattern.FindString(lines[4])))
		testFalse := uint64(Expect(strconv.Atoi)(numPattern.FindString(lines[5])))
		testFun := func(n uint64) uint64 {
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

		printItem := func(item uint64) uint64 {
			fmt.Printf("%d, ", item)
			return item
		}

		m.items.ApplyToAll(printItem)
		fmt.Println()
	}

	fmt.Println()
}

func (m *monkey) doItem() (worry, nextMonkeyIndex uint64) {
	item := m.items.PopHead()
	worry = m.operation(item)
	nextMonkeyIndex = m.test(worry)
	m.inspectCount++
	return
}

func (d *day) Run1() string {

	d.printMonkeys()

	for i := 0; i < 20; i++ {
		for mi := range d.monkeys {
			m := &d.monkeys[mi]
			for !m.items.Empty() {
				worry, nextMonkeyIndex := m.doItem()
				d.monkeys[nextMonkeyIndex].items.PushTail(worry)
			}
		}
	}

	d.printMonkeys()

	getInspectCount := func(m monkey) int { return m.inspectCount }
	counts := Map(getInspectCount, d.monkeys)
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	return fmt.Sprint(counts[0] * counts[1])
	// return ""
}

func (d *day) Run2() string {
	d.parse()
	d.printMonkeys()

	for i := 0; i < 20; i++ {
		for mi := range d.monkeys {
			m := &d.monkeys[mi]
			for !m.items.Empty() {
				worry, nextMonkeyIndex := m.doItem()
				d.monkeys[nextMonkeyIndex].items.PushTail(worry)
			}
		}
	}

	d.printMonkeys()

	getInspectCount := func(m monkey) int { return m.inspectCount }
	counts := Map(getInspectCount, d.monkeys)
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	return fmt.Sprint(counts[0] * counts[1])
}

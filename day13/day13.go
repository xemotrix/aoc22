package day13

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
	input       string
	parsedInput [][]signal
}

func BuildDay13() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type signal struct {
	value    *int
	contents []signal
}

func (s *signal) isValue() bool {
	return s.value != nil
}

func (s *signal) isEmpty() bool {
	return s.value == nil && len(s.contents) == 0
}

var digitRegex = regexp.MustCompile(`\d+`)
var breakpointRegex = regexp.MustCompile(`\]|,`)

func findSubSignal(s string) string {
	depth := 0
	for i, c := range s {
		if c == '[' {
			depth++
		} else if c == ']' {
			depth--
		}
		if depth == 0 {
			return s[:i+1]
		}
	}
	return ""
}

func parseLine(s string) signal {

	sig := signal{}
	if s == "" {
		return sig
	}

	if s[0] != '[' || s[len(s)-1] != ']' {
		panic("Error parsing, expected string enclosed by '[' and ']'")
	}

	stripped := s[1 : len(s)-1]

	for i := 0; i < len(stripped); i++ {
		c := stripped[i]
		if c >= '0' && c <= '9' {
			val := Expect(strconv.Atoi)(digitRegex.FindString(stripped[i:]))

			sig.contents = append(sig.contents, signal{value: &val})

			if nextBreakPoint := breakpointRegex.FindStringIndex(stripped[i:]); nextBreakPoint != nil {
				i += nextBreakPoint[0]
			}
			continue
		}
		if c == ',' {
			continue
		}
		if c == '[' {

			subSignalStr := findSubSignal(stripped[i:])
			endOfSubSignal := len(subSignalStr)

			subSignal := parseLine(subSignalStr)
			sig.contents = append(sig.contents, subSignal)
			i += endOfSubSignal
		}
	}
	return sig
}

func (d *day) parse() {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")
	lines := Map(FCurry(strings.Split, "\n"), blocks)
	d.parsedInput = Map(Curry(Map[string, signal], parseLine), lines)
}

func comparePair(a, b signal) *bool {
	if a.isValue() && b.isValue() {
		if *a.value == *b.value {
			return nil
		}
		res := *a.value < *b.value

		return &res
	}

	if !a.isValue() && !b.isValue() {
		f := false
		t := true
		for i := range a.contents {
			if i >= len(b.contents) {
				return &f
			}
			if ok := comparePair(a.contents[i], b.contents[i]); ok != nil {
				return ok
			}
		}
		if len(a.contents) == len(b.contents) {
			return nil
		}
		return &t
	}

	if a.isValue() {
		wrapped := signal{contents: []signal{a}}
		return comparePair(wrapped, b)
	}

	wrapped := signal{contents: []signal{b}}
	return comparePair(a, wrapped)
}

const (
	Reset   = "\033[0m"
	GreenBG = "\u001b[46m"
	RedBG   = "\u001b[41m"
	BlackFG = "\u001b[30m"
)

func (s *signal) toString() string {
	if s.isValue() {
		return strconv.Itoa(*s.value)
	}
	if s.isEmpty() {
		return "[]"
	}
	res := "["
	for _, c := range s.contents {
		res += c.toString() + ","
	}
	res = res[:len(res)-1] + "]"
	return res
}

func (d *day) Run1() string {
	res := []int{}
	for i, pair := range d.parsedInput {
		if *comparePair(pair[0], pair[1]) {
			res = append(res, i+1)

			continue
		}

	}
	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })
	return fmt.Sprint(sumReduce(res))
}

type signalSlice []signal

func (is signalSlice) Len() int {
	return len(is)
}

func (is signalSlice) Less(i, j int) bool {
	return *comparePair(is[i], is[j])
}

func (is signalSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func (d *day) Run2() string {

	divider1 := "[[2]]"
	divider2 := "[[6]]"
	lines := strings.Split(strings.Replace(input+"\n"+divider1+"\n"+divider2, "\n\n", "\n", -1), "\n")

	signals := Map(parseLine, lines)
	sort.Sort(signalSlice(signals))

	res := 1
	for i, s := range signals {
		str := s.toString()
		if str == divider1 || str == divider2 {
			res *= i + 1
		}
	}
	return fmt.Sprint(res)
}

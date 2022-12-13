package day13

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"regexp"
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

func (d *day) parse() {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")
	lines := Map(FCurry(strings.Split, "\n"), blocks)

	parseDepth := 0
	findSubSignal := func(s string) string {
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

	digitRegex := regexp.MustCompile(`\d+`)
	breakpointRegex := regexp.MustCompile(`\]|,`)

	var parseLine func(string) signal
	parseLine = func(s string) signal {
		parseDepth++
		defer func() { parseDepth-- }()

		depthStr := strings.Repeat("  ", parseDepth)
		fmt.Printf("%s<>Commencing parse: '%s'\n", depthStr, s)

		sig := signal{}

		if s == "" {
			fmt.Printf("%sEmpty string, returning\n", depthStr)
			return sig
		}

		if s[0] != '[' || s[len(s)-1] != ']' {
			panic("Error parsing, expected string enclosed by '[' and ']'")
		}

		stripped := s[1 : len(s)-1]

		for i := 0; i < len(stripped); i++ {
			c := stripped[i]
			fmt.Printf("%s·Check char %s\n", depthStr, string(c))
			if c >= '0' && c <= '9' {
				fmt.Printf("%s·>%s is number, added\n", depthStr, string(c))
				val := Expect(strconv.Atoi)(digitRegex.FindString(stripped[i:]))

				sig.contents = append(sig.contents, signal{value: &val})

				if nextBreakPoint := breakpointRegex.FindStringIndex(stripped[i:]); nextBreakPoint != nil {
					i += nextBreakPoint[0]
				}
				continue
			}
			if c == ',' {
				// fmt.Printf("%s·>%s is ',' skip\n", depthStr, string(c))
				continue
			}
			if c == '[' {

				subSignalStr := findSubSignal(stripped[i:])
				endOfSubSignal := len(subSignalStr)

				fmt.Printf("%sRecursion! -> '%s'\n", depthStr, subSignalStr)
				subSignal := parseLine(subSignalStr)
				fmt.Printf("%sLeft -> %s\n", depthStr, stripped[i+endOfSubSignal:])
				sig.contents = append(sig.contents, subSignal)
				i += endOfSubSignal
			}
		}
		return sig
	}

	d.parsedInput = Map(Curry(Map[string, signal], parseLine), lines)
}

func comparePair(a, b signal) *bool {
	// time.Sleep(100 * time.Millisecond)
	fmt.Printf("\nComparing: \n->'%v'\n->'%v'\n", a.toString(), b.toString())
	fmt.Printf("vals : %v, %v\n", a.isValue(), b.isValue())

	if a.isValue() && b.isValue() {
		fmt.Printf("Both are values: %d, %d\n", *a.value, *b.value)
		if *a.value == *b.value {
			fmt.Printf("·a==b: %d, %d\n", *a.value, *b.value)
			return nil
		}
		res := *a.value < *b.value

		fmt.Printf("·a<b: %v\n", res)
		return &res
	}

	if !a.isValue() && !b.isValue() {
		fmt.Printf("Both are lists\n")
		f := false
		t := true
		for i := range a.contents {
			if i >= len(b.contents) {
				fmt.Printf("· ended and b run out : false\n")
				return &f
			}
			if ok := comparePair(a.contents[i], b.contents[i]); ok != nil {
				fmt.Printf("· Inner pair compared and ok = %v\n", *ok)
				return ok
			}
		}
		if len(a.contents) == len(b.contents) {
			fmt.Printf("both lists are ==, returning nil \n")
			return nil
		}
		fmt.Printf("· ended and a run out : true\n")
		return &t
	}

	if a.isValue() {
		fmt.Printf("wrap a\n")
		wrapped := signal{contents: []signal{a}}
		return comparePair(wrapped, b)
	}

	fmt.Printf("wrap b\n")
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
		fmt.Printf("\n=== PAIR #%d ===\n", i+1)
		fmt.Printf("s1: %v\ns2: %v\n", pair[0].toString(), pair[1].toString())
		if *comparePair(pair[0], pair[1]) {
			fmt.Printf("|>INPUT (PAIR %d) IN %s%sRIGHT ORDER%s <|\n", i+1, GreenBG, BlackFG, Reset)
			res = append(res, i+1)

			continue
		}

		fmt.Printf("|>INPUT (PAIR %d) IN %s%sWRONG%s ORDER <|\n", i+1, RedBG, BlackFG, Reset)
	}
	fmt.Println(res)
	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })
	return fmt.Sprint(sumReduce(res))
}

func (d *day) Run2() string {
	return ""
}

// [4 7 9 10 11 14 15 16 18 19 20 23 25 27 28 30 33 34 37 41 42 43 46 49 50 53 54 61 63 65 66 69 70 71 73 74 78 80 81 82 84 85 88 89 90 91 94 95 96 99 100 101 103 110 111 112 116 117 118 119 120 122 123 126 128 130 132 137 138 139 140 144 145 147 149]

// [4 6 7 8 9 10 11 14 15 16 18 19 20 22 23 25 27 28 30 33 34 37 38 41 42 43 46 49 50 53 54 58 63 65 66 69 70 71 73 74 78 80 81 82 84 85 88 89 90 91 94 95 96 99 100 101 103 105 110 112 116 117 118 120 122 123 126 128 130 132 137 138 139 140 144 145 147 149]

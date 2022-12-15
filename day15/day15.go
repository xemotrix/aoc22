package day15

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type day struct {
	input   string
	sensors []sensor
}

func BuildDay15() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type loc struct {
	y, x int
}

type sensor struct {
	loc    loc
	beacon loc
	dist   int
}

func manhattan(a, b loc) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

func (d *day) parse() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	digitsRegex := regexp.MustCompile(`(-?)\d+`)
	matStr := Map(FCurry(digitsRegex.FindAllString, -1), lines)
	matInt := Map(Curry(Map[string, int], Expect(strconv.Atoi)), matStr)
	toSensor := func(s []int) sensor {
		sen := sensor{
			loc: loc{
				x: s[0],
				y: s[1],
			},
			beacon: loc{
				x: s[2],
				y: s[3],
			},
		}
		sen.dist = manhattan(sen.loc, sen.beacon)
		return sen
	}
	d.sensors = Map(toSensor, matInt)
}

type intervals []interval

type interval struct {
	min int
	max int
}

func (yi *interval) Len() int {
	return yi.max - yi.min
}

func (yi intervals) Len() int {
	return len(yi)
}

func (yi intervals) Less(i, j int) bool {
	return yi[i].min < yi[j].min
}

func (yi intervals) Swap(i, j int) {
	yi[i], yi[j] = yi[j], yi[i]
}

func (d *day) buildIntervals(row int) intervals {
	inter := []interval{}
	for _, s := range d.sensors {
		distY := int(math.Abs(float64(s.loc.y - row)))
		rad := s.dist - distY
		if rad < 0 {
			continue
		}
		newInter := interval{s.loc.x - rad, s.loc.x + rad}
		inter = append(inter, newInter)
	}
	sort.Sort(intervals(inter))
	cleanIntervals := []interval{}
	buildInter := inter[0]
	for i := 1; i < len(inter); i++ {
		if inter[i].min <= buildInter.max {
			buildInter.max = Max(inter[i].max, buildInter.max)
			continue
		}
		cleanIntervals = append(cleanIntervals, buildInter)
		buildInter = inter[i]
	}
	cleanIntervals = append(cleanIntervals, buildInter)
	return cleanIntervals
}

func (i intervals) cap(min, max int) {
	if i[0].min < min {
		i[0].min = min
	}
	if i[len(i)-1].max > max {
		i[len(i)-1].max = max
	}
}

func (d *day) Run1() string {
	row := 2000000
	cleanIntervals := d.buildIntervals(row)
	lens := Map(func(i interval) int { return i.Len() }, cleanIntervals)
	res := Reduce(func(a, b int) int { return a + b }, lens)
	return fmt.Sprint(res)
}

func printIntervals(num int, is intervals, min, max int) {
	str := []rune(strings.Repeat("·", max-min+1))
	for _, inter := range is {
		for j := inter.min; j <= inter.max; j++ {
			str[j] = '#'
		}
	}
	fmt.Println(num, "\t", string(str))
}

func (d *day) Run2() string {
	min, max := 0, 4000000
	for row := min; row <= max; row++ {
		cleanIntervals := d.buildIntervals(row)
		cleanIntervals.cap(min, max)

		lens := Map(func(i interval) int { return i.Len() + 1 }, cleanIntervals)
		holes := max + 1 - Reduce(func(a, b int) int { return a + b }, lens)

		if holes > 0 {
			return fmt.Sprint(4000000*(cleanIntervals[1].min-1) + row)
		}

	}

	return ""
}
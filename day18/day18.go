package day18

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
	input                              string
	cubes                              []cube
	cubeMap                            map[cube]*any
	gasCubes                           map[cube]*any
	maxX, minX, maxY, minY, maxZ, minZ int
	outsideSurface                     int
}

func BuildDay18() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type cube struct {
	x, y, z int
}

func (d *day) parse() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	splitCommas := FCurry(strings.Split, ",")
	toInt := Curry(Map[string, int], Expect(strconv.Atoi))
	toCube := func(s []int) cube {
		return cube{s[0], s[1], s[2]}
	}
	parseLine := Pipe(toInt, toCube)
	d.cubes = Map(Pipe(splitCommas, parseLine), lines)
	d.cubeMap = map[cube]*any{}
	for _, c := range d.cubes {
		d.cubeMap[c] = nil
	}
}

func Abs(i int) int {
	return int(math.Abs(float64(i)))
}

func (c cube) adjacent(other cube) bool {
	return Abs(c.x-other.x) <= 1 ||
		Abs(c.y-other.y) <= 1 ||
		Abs(c.z-other.z) <= 1
}

func (d *day) calcNonAdjacent(c cube) (res int) {

	cubesToCheck := []cube{
		{c.x - 1, c.y, c.z},
		{c.x + 1, c.y, c.z},
		{c.x, c.y - 1, c.z},
		{c.x, c.y + 1, c.z},
		{c.x, c.y, c.z - 1},
		{c.x, c.y, c.z + 1},
	}

	for _, ctc := range cubesToCheck {
		if _, ok := d.cubeMap[ctc]; !ok {
			res++
		}
	}

	return res
}

func (d *day) Run1() string {
	res := 0
	for _, c := range d.cubes {
		res += d.calcNonAdjacent(c)
	}
	return fmt.Sprint(res)
}

func (d *day) calcGas(start cube) {
	cubesToCheck := []cube{
		{start.x - 1, start.y, start.z},
		{start.x + 1, start.y, start.z},
		{start.x, start.y - 1, start.z},
		{start.x, start.y + 1, start.z},
		{start.x, start.y, start.z - 1},
		{start.x, start.y, start.z + 1},
	}

	outOfBounds := func(c cube) bool {
		return c.x > d.maxX+10 ||
			c.x < d.minX-10 ||
			c.y > d.maxY+10 ||
			c.y < d.minY-10 ||
			c.z > d.maxZ+10 ||
			c.z < d.minZ-10
	}

	for _, ctc := range cubesToCheck {
		if outOfBounds(ctc) {
			continue
		}
		if _, ok := d.gasCubes[ctc]; ok {
			continue
		}
		if _, ok := d.cubeMap[ctc]; ok {
			d.outsideSurface++
			continue
		}
		d.gasCubes[ctc] = nil
		d.calcGas(ctc)
	}
}

func (d *day) Run2() string {
	getX := func(c cube) int { return c.x }
	getY := func(c cube) int { return c.y }
	getZ := func(c cube) int { return c.z }

	xs := Map(getX, d.cubes)
	ys := Map(getY, d.cubes)
	zs := Map(getZ, d.cubes)

	max := Curry(Reduce[int], Max[int])
	min := Curry(Reduce[int], Min[int])

	d.maxX, d.minX = max(xs), min(xs)
	d.maxY, d.minY = max(ys), min(ys)
	d.maxZ, d.minZ = max(zs), min(zs)

	start := cube{max(xs), max(ys), max(zs)}
	d.gasCubes = map[cube]*any{}
	d.calcGas(start)

	return fmt.Sprint(d.outsideSurface)
}

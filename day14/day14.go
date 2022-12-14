package day14

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/gosuri/uilive"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput [][][2]int
	cave        cave
}

func BuildDay14() *day {
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
	mat := Map(FCurry(strings.Split, " -> "), lines)
	toInstr := func(s string) [2]int {
		splitted := strings.Split(s, ",")
		return [2]int{
			Expect(strconv.Atoi)(splitted[0]),
			Expect(strconv.Atoi)(splitted[1]),
		}
	}
	d.parsedInput = Map(Curry(Map[string, [2]int], toInstr), mat)
}

type material int

const (
	Air material = iota
	Rock
	Sand
	Floor
)

type cell struct {
	material material
}

type cave struct {
	caveMap    map[int]map[int]cell
	activeSand *struct {
		x int
		y int
	}
	sandStable int
	voidLvl    int
	part       int
	minX       int
	maxX       int
	minY       int
	maxY       int
}

func buildCave(part int) cave {
	return cave{
		caveMap: map[int]map[int]cell{},
		activeSand: &struct {
			x int
			y int
		}{500, 0},
		part: part,
		minX: 500,
		maxX: 500,
		minY: 0,
		maxY: 0,
	}
}

func (c *cave) Update(x, y int, mat material) {

	if x > c.maxX {
		c.maxX = x
	}
	if x < c.minX {
		c.minX = x
	}
	if y > c.maxY {
		c.maxY = y
	}
	if y < c.minY {
		c.minY = y
	}

	if _, ok := c.caveMap[x]; !ok {
		c.caveMap[x] = map[int]cell{}
	}
	if _, ok := c.caveMap[x][y]; !ok {
		c.caveMap[x][y] = cell{material: Air}
	}

	cmap := c.caveMap[x][y]
	cmap.material = mat
	c.caveMap[x][y] = cmap

}

func (c *cave) addRockPath(part int, a, b [2]int) {
	xStart := Min(a[0], b[0])
	xEnd := Max(a[0], b[0])
	yStart := Min(a[1], b[1])
	yEnd := Max(a[1], b[1])

	if part == 1 && yEnd >= c.voidLvl {
		c.voidLvl = yEnd + 1
		c.maxY = c.voidLvl
	}
	if part == 2 && yEnd+2 > c.voidLvl {
		c.voidLvl = yEnd + 2
		c.maxY = c.voidLvl
	}

	for x := xStart; x <= xEnd; x++ {
		for y := yStart; y <= yEnd; y++ {
			c.Update(x, y, Rock)
		}
	}
}

func (c *cave) look(x, y int) material {
	if y == c.voidLvl {
		return Floor
	}
	if _, ok := c.caveMap[x]; !ok {
		return Air
	}
	if _, ok := c.caveMap[x][y]; !ok {
		return Air
	}
	return c.caveMap[x][y].material

}

const (
	Reset       = "\033[0m"
	SandColor   = "\u001b[43m"
	RockColor   = "\u001b[47;1m"
	ActiveColor = "\u001b[41;1m"
	FloorColor  = "\u001b[46m"
	Space       = "  "
)

func (c *cave) render() string {

	minX := c.minX - 3
	maxX := c.maxX + 3

	minY := c.minY - 3
	maxY := c.maxY + 3

	lines := make([]string, maxY-minY)
	for x := minX - 1; x < maxX; x++ {
		for y := 0; y < maxY-minY; y++ {

			if c.activeSand.x == x && c.activeSand.y == y {
				lines[y] += ActiveColor + Space + Reset
				continue
			}

			switch c.look(x, y) {
			case Rock:
				lines[y] += RockColor + Space + Reset
			case Sand:
				lines[y] += SandColor + Space + Reset
			case Air:
				lines[y] += Space
			case Floor:
				lines[y] += FloorColor + Space + Reset
			}
		}
	}

	return "\n" + strings.Join(lines, "\n")
}

var focus int

func (c *cave) renderScreen(rows int) string {
	s := strings.Split(c.render(), "\n")[1:]

	if c.activeSand.y > focus {
		focus = c.activeSand.y
	}
	bot := Max(focus+10, rows)
	top := bot - rows

	screen := s[top:Min(bot+1, len(s))]

	w := Min(len(screen[1]), len(screen[2]))
	surround := func(s string) string {
		return "┃" + s + "┃" + fmt.Sprint(w)
	}
	screen = Map(surround, screen)

	screen = append(
		[]string{"┏" + strings.Repeat("━", w) + "┓"},
		screen...,
	)

	screen = append(
		screen,
		"┗"+strings.Repeat("━", w)+"┛",
	)
	return "\n" + strings.Join(screen, "\n") + "\n"

}

func (c *cave) simStep(part int) bool {

	as := c.activeSand
	matDown := c.look(as.x, as.y+1)
	if matDown == Air {
		as.y++
		return as.y < c.voidLvl
	}
	if part == 1 && matDown == Floor {
		return false
	}

	matLeftDown := c.look(as.x-1, as.y+1)
	if matLeftDown == Air {
		as.x--
		as.y++
		return as.y < c.voidLvl
	}

	matRightDown := c.look(as.x+1, as.y+1)
	if matRightDown == Air {
		as.x++
		as.y++
		return as.y < c.voidLvl
	}

	c.sandStable++
	c.Update(as.x, as.y, Sand)

	pos := c.activeSand.y

	as.x = 500
	as.y = 0

	if part == 2 {
		return pos != 0
	}

	return true

}

func (d *day) runSim(part int) string {
	c := buildCave(part)

	for _, line := range d.parsedInput {
		for i := 1; i < len(line); i++ {
			c.addRockPath(part, line[i-1], line[i])
		}
	}

	writer := uilive.New()
	writer.Start()

	for i := 0; true; i++ {
		// if i%30 == 0 {
		// 	fmt.Fprintf(writer, c.renderScreen(40))
		// }
		// fmt.Println("\nStep", i, "sand:", c.sandStable, c.renderScreen(30))
		// fmt.Printf("Boundaries x:[%d,%d] y:[%d,%d]\n", c.minX, c.maxX, c.minY, c.maxY)
		goOn := c.simStep(part)
		if !goOn {
			break
		}
	}

	return fmt.Sprint(c.sandStable)
}

func (d *day) Run1() string {
	return d.runSim(1)
}

func (d *day) Run2() string {
	return d.runSim(2)
}

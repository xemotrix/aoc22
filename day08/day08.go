package day08

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	. "aoc/evilgo"
)

//go:embed input.txt
var input string

type day struct {
	input       string
	parsedInput forest
}

func BuildDay08() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type forest [][]int

func (d *day) parse() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	strMat := Map(FCurry(strings.Split, ""), lines)
	d.parsedInput = Map(Curry(Map[string, int], Expect(strconv.Atoi)), strMat)
}

func transpose[T any](f [][]T) [][]T {
	w := len(f[0])
	h := len(f)
	newForest := make([][]T, w)
	for i := 0; i < w; i++ {
		newForest[i] = make([]T, h)
	}
	for i, l := range f {
		for j, num := range l {
			newForest[j][i] = num
		}
	}

	return newForest
}

func reverse[T any](f [][]T) [][]T {
	w := len(f[0])
	h := len(f)
	newMat := make([][]T, h)

	for i, l := range f {
		newMat[i] = make([]T, w)
		for j, k := 0, len(l)-1; j <= k; j, k = j+1, k-1 {
			newMat[i][j], newMat[i][k] = f[i][k], f[i][j]
		}
	}
	return newMat
}

func checkLineVisibility(s []int) []bool {
	res := make([]bool, len(s))
	max := -1
	for i, ele := range s {
		if ele > max {
			res[i] = true
			max = ele
		}
	}
	return res
}

func printMat[T any](m [][]T) {
	for i := range m {
		fmt.Println(m[i])
	}
}

func (d *day) Run1() string {

	type matFunc[T any] func(in [][]T) [][]T
	viewFuncs := []struct {
		forInt  matFunc[int]
		forBool matFunc[bool]
	}{
		{
			Identity[[][]int],
			Identity[[][]bool],
		},
		{
			reverse[int],
			reverse[bool],
		},
		{
			transpose[int],
			transpose[bool],
		},
		{
			Pipe(transpose[int], reverse[int]),
			Pipe(reverse[bool], transpose[bool]),
		},
	}

	visiMaps := make([][][]bool, 4)

	checkVisibility := Curry(Map[[]int, []bool], checkLineVisibility)
	for i, vf := range viewFuncs {
		newMat := vf.forInt(d.parsedInput)
		visi := checkVisibility(newMat)
		visiMaps[i] = vf.forBool(visi)
	}

	sumReduce := Curry(Reduce[int], func(a, b int) int { return a + b })

	scoreLine := Pipe(
		Curry(
			Map[bool, int],
			func(b bool) int {
				if b {
					return 1
				}
				return 0
			},
		),
		sumReduce,
	)

	matOr := func(a [][]bool, b [][]bool) [][]bool {
		for i := range a {
			for j := range a[i] {
				a[i][j] = a[i][j] || b[i][j]
			}
		}
		return a
	}

	res := sumReduce(Map(scoreLine, Reduce(matOr, visiMaps)))
	return fmt.Sprint(res)
}

func (d *day) calcScenicScore(y, x int) int {
	w := len(d.parsedInput[0])
	h := len(d.parsedInput)

	checkVisiScore := func(tree int, s []int) (visi int) {
		for _, t := range s {
			visi++
			if t >= tree {
				break
			}
		}
		return
	}

	mulReduce := Curry(Reduce[int], func(a, b int) int { return a * b })

	treeHeight := d.parsedInput[y][x]

	res := mulReduce([]int{
		checkVisiScore(treeHeight, d.parsedInput[y][x+1:]),
		checkVisiScore(treeHeight, reverse(d.parsedInput)[y][w-x:]),
		checkVisiScore(treeHeight, transpose(d.parsedInput)[x][y+1:]),
		checkVisiScore(treeHeight, reverse(transpose(d.parsedInput))[x][h-y:]),
	})

	return res
}

func (d *day) Run2() string {
	max := 0
	for i := range d.parsedInput {
		for j := range d.parsedInput[i] {
			if score := d.calcScenicScore(i, j); score > max {
				max = score
			}
		}
	}
	return fmt.Sprint(max)
}

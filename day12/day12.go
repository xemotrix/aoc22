package day12

import (
	. "aoc/evilgo"
	_ "embed"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

//go:embed input.txt
var input string

type day struct {
	input string
	graph Graph[int]
	// startUUID   string
	// endUUID     string
	parsedInput [][]string
}

func BuildDay12() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d
}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

func (d *day) buildGraph(part int) {
	g := BuildGraph[int]()

	type matFunc[T any] func(in [][]T) [][]T
	transformList := []matFunc[string]{
		Identity[[][]string],
		Reverse[string],
		Transpose[string],
		Pipe(Transpose[string], Reverse[string]),
	}

	getHeightDiff := func(a, b string) int {
		aw := []rune(a)[0]
		bw := []rune(b)[0]

		if a[:1] == "S" {
			aw = 'a'
		}
		if a[:1] == "E" {
			aw = 'z'
		}
		if b[:1] == "S" {
			bw = 'a'
		}
		if b[:1] == "E" {
			bw = 'z'
		}
		return int(bw - aw)
	}

	buildRowEdges := func(s []string) any {
		for i := 1; i < len(s); i++ {
			if part == 1 {
				w := getHeightDiff(s[i-1], s[i])
				if w <= 1 {
					g.AddEdge(s[i-1], s[i], 1)
				}
			} else {
				w := getHeightDiff(s[i], s[i-1])
				if w <= 1 {
					g.AddEdge(s[i-1], s[i], 1)
				}
			}
		}
		return nil
	}

	for _, f := range transformList {
		Map(buildRowEdges, f(d.parsedInput))
	}
	d.graph = g
}

func (d *day) parse() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	mat := Map(FCurry(strings.Split, ""), lines)
	mat = Map(
		Curry(Map[string, string],
			func(s string) string {
				if s != "S" && s != "E" {
					return s + "-" + uuid.New().String()
				}
				return s
			},
		),
		mat,
	)
	d.parsedInput = mat
}

func (d *day) dijkstra(start, end string) string {

	dist := map[*GraphNode[int]]int{}
	prev := map[*GraphNode[int]]*GraphNode[int]{}

	q := BuildPQueue(func(a, b *GraphNode[int]) bool {
		return dist[a] < dist[b]
	})

	for _, vertex := range d.graph.Nodes {
		dist[vertex] = 999999999
		prev[vertex] = nil
		q.Push(vertex)
	}

	dist[d.graph.Nodes[start]] = 0
	q.Update(d.graph.Nodes[start])

	var visited []string
	// dijkstra
	for !q.Empty() {
		u := q.Pop()
		visited = append(visited, u.NodeName)
		for _, edge := range d.graph.Edges[u.NodeName] {
			v := edge.To
			alt := dist[u] + edge.Weight
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				q.Update(v)
			}
		}
	}

	// backtracking
	currentNode := d.graph.Nodes[end]
	path := []string{currentNode.NodeName}
	for {
		// d.render(start, end, path...)
		p := prev[currentNode]
		if p == nil || p.NodeName == start {
			break
		}
		path = append(path, p.NodeName)
		currentNode = p
	}

	d.render(start, end, path...)
	return fmt.Sprint(len(path))
}

const (
	Reset   = "\033[0m"
	GreenBG = "\u001b[46m"
	RedBG   = "\u001b[41m"
	BlackFG = "\u001b[30m"
)

func (d *day) render(start, end string, highlights ...string) {
	for _, l := range d.parsedInput {
		fmt.Println()
		for _, ele := range l {
			eleChr := ele[:1]
			if ele == start || ele == end {
				eleChr = RedBG + BlackFG + eleChr + Reset
			} else {
				for _, h := range highlights {
					if h == ele {
						eleChr = GreenBG + BlackFG + eleChr + Reset
					}
				}
			}
			fmt.Print(eleChr)
		}
	}
	fmt.Println()
}

func (d *day) Run1() string {
	d.buildGraph(1)
	return d.dijkstra("S", "E")
}

func (d *day) Run2() string {
	d.parse()
	unifyAs := func(s string) string {
		if s[:1] == "a" {
			return "a-lol"
		}
		return s
	}
	d.parsedInput = Map(Curry(Map[string, string], unifyAs), d.parsedInput)
	d.buildGraph(2)

	return fmt.Sprint(d.dijkstra("E", "a-lol"))
}

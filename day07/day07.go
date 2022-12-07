package day07

import (
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
	parsedInput []cmd
}

func BuildDay07() *day {
	d := day{
		input: input,
	}
	d.parse()
	return &d

}

func (d *day) Run() (string, string) {
	return d.Run1(), d.Run2()
}

type cmd struct {
	command string
	output  []string
}

func (d *day) parse() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i][0] == '$' {
			c := cmd{
				command: lines[i][2:],
			}
			i++
			for i < len(lines) && lines[i][0] != '$' {
				c.output = append(c.output, lines[i])
				i++
			}
			d.parsedInput = append(d.parsedInput, c)
			i--
		}
	}
}

type node struct {
	name     string
	nodeType string
	size     int
	children []*node
	parent   *node
}

const (
	nodeTypeDir  = "DIR"
	nodeTypeFile = "FILE"
)

var (
	reCmdCd = regexp.MustCompile(`^cd (.*)$`)
	reCmdLs = regexp.MustCompile(`^ls$`)
)

func (d *day) buildFileTree() node {
	currentNode := &node{
		name:     "/",
		nodeType: nodeTypeDir,
		size:     0,
		parent:   nil,
	}
	parentNode := currentNode

	for i := 1; i < len(d.parsedInput); i++ {
		cmd := d.parsedInput[i]
		switch cmdStr := cmd.command; {
		case reCmdCd.MatchString(cmdStr):
			if cmdStr[3:] == ".." {
				currentNode = currentNode.parent
			} else {
				currentNode = &node{
					name:     cmdStr[3:],
					nodeType: nodeTypeDir,
					size:     0,
					parent:   currentNode,
				}
				currentNode.parent.children = append(currentNode.parent.children, currentNode)
			}
		case reCmdLs.MatchString(cmdStr):
			for _, chi := range cmd.output {
				if chi[:3] == "dir" {
					continue
				}
				parts := strings.Split(chi, " ")
				size, _ := strconv.Atoi(parts[0])
				currentNode.children = append(currentNode.children, &node{
					name:     parts[1],
					nodeType: nodeTypeFile,
					size:     size,
					parent:   currentNode,
				})
			}

		default:
			panic(fmt.Sprintf("Error parsing, unkown command: %v", cmdStr))
		}
	}
	return *parentNode
}

func printTree(n *node, depth int) {
	fmt.Printf("\n%s> %s (%v): %v",
		strings.Repeat("\t", depth),
		n.name,
		n.nodeType,
		n.size,
	)
	for _, c := range n.children {
		printTree(c, depth+1)
	}
}

func calcSizes(n *node) int {
	if n.size != 0 {
		return n.size
	}
	sum := 0
	for _, c := range n.children {
		sum += calcSizes(c)
	}
	n.size = sum
	return sum
}

func findDirs(n *node) []*node {
	found := []*node{}
	if n.nodeType == nodeTypeDir && n.size < 100000 {
		found = append(found, n)
	}
	for _, c := range n.children {
		found = append(found, findDirs(c)...)
	}
	return found
}

func (d *day) Run1() string {
	root := d.buildFileTree()
	calcSizes(&root)

	foundDirs := findDirs(&root)

	res := 0
	for _, d := range foundDirs {
		res += d.size
	}
	return fmt.Sprint(res)
}

func findDirs2(n *node, size int) []*node {
	found := []*node{}
	if n.nodeType == nodeTypeDir && n.size > size {
		found = append(found, n)
	}
	for _, c := range n.children {
		found = append(found, findDirs2(c, size)...)
	}
	return found
}

func (d *day) Run2() string {
	root := d.buildFileTree()
	calcSizes(&root)

	unused := 30000000 - (70000000 - root.size)

	foundDirs := findDirs2(&root, unused)

	minDir := root
	for _, d := range foundDirs {
		if d.size < minDir.size {
			minDir = *d
		}
	}

	return fmt.Sprint(minDir.size)
}

package main

import (
	"aoc/day10"
	"fmt"
)

func main() {
	var d Day = day10.BuildDay10()
	res1, res2 := d.Run()

	fmt.Printf("Result part 1: %v\nResult part 2: %v\n", res1, res2)
}

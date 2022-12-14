package main

import (
	"aoc/day14"
	"fmt"
)

func main() {
	var d Day = day14.BuildDay14()
	res1, res2 := d.Run()

	fmt.Printf("Result part 1: %v\nResult part 2: %v\n", res1, res2)
}

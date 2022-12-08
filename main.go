package main

import (
	"aoc/day08"
	"fmt"
)

func main() {
	var d Day = day08.BuildDay08()
	res1, res2 := d.Run()

	fmt.Printf("Result part 1: %v\nResult part 2: %v\n", res1, res2)
}

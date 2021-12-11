package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	depths := getDepths(filename)
	fmt.Println(getIncreases(depths, 1))
	fmt.Println(getIncreases(depths, 3))
}

func getIncreases(depths []int, windowSize int) int {
	increases := 0
	for i, after := range depths[windowSize:] {
		before := depths[i]
		if after > before {
			increases++
		}
	}
	return increases
}

func getDepths(filename string) []int {
	lines := aoc.GetInputLines(filename)
	return aoc.ParseInts(lines)
}

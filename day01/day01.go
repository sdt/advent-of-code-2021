package main

import (
	"advent-of-code/common"
	"fmt"
)

func main() {
	filename := common.GetFilename()
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
	lines := common.GetInputLines(filename)
	return common.ParseInts(lines)
}

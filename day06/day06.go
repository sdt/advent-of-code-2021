package main

import (
	"advent-of-code/common"
	"fmt"
	"strings"
)

var memo [500]int

func main() {
	filename := common.GetFilename()
	fmt.Println(part1(filename))
	fmt.Println(part2(filename))
}

func part1(filename string) int {
	startCycles := getInput(filename)
	total := 0
	for _, start := range startCycles {
		total += totalFish(80 - start - 1)
	}
	return total
}

func part2(filename string) int {
	startCycles := getInput(filename)
	total := 0
	for _, start := range startCycles {
		total += totalFish(256 - start - 1)
	}
	return total
}

func totalFish(days int) int {
	if days < 0 {
		return 1
	}
	if count := memo[days]; count > 0 {
		return count
	}
	count := 1

	for n := days; n >= 0; n -= 7 {
		count += totalFish(n - 8 - 1)
	}

	memo[days] = count
	return count
}

func getInput(filename string) []int {
	lines := common.GetInputLines(filename)
	words := strings.Split(lines[0], ",")
	return common.ParseInts(words)
}

package main

import (
	"advent-of-code/common"
	"fmt"
	"math"
	"strings"
)

type FuelFunction func(int) int

func main() {
	filename := common.GetFilename()
	positions := getInput(filename)

	fmt.Println(part1(positions))
	fmt.Println(part2(positions))
}

func part1(positions []int) int {
	return getLeastFuel(positions, func(dist int) int {
		return abs(dist)
	})
}

func part2(positions []int) int {
	return getLeastFuel(positions, func(dist int) int {
		dist = abs(dist)
		return (dist*dist + dist) / 2
	})
}

func getLeastFuel(positions []int, getFuel FuelFunction) int {
	min := math.MaxInt
	max := math.MinInt

	for _, position := range positions {
		if position < min {
			min = position
		}
		if position > max {
			max = position
		}
	}

	leastTotal := math.MaxInt

	for dest := min; dest <= max; dest++ {
		total := 0
		for _, position := range positions {
			total += getFuel(position - dest)
		}
		if total < leastTotal {
			leastTotal = total
		}
	}
	return leastTotal
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getInput(filename string) []int {
	lines := common.GetInputLines(filename)
	words := strings.Split(lines[0], ",")
	return common.ParseInts(words)
}

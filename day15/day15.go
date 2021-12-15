package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
)

type Cavern struct {
	rows, cols int
	riskLevel []int
}

type Position struct {
	row, col int
}

func main() {
	filename := aoc.GetFilename()
	cavern := parseCavern(aoc.GetInputLines(filename))

	fmt.Println(part1(&cavern))
}

func part1(cavern *Cavern) int {
	lowest := make([]int, cavern.rows * cavern.cols)

	//printGrid(cavern.rows, cavern.cols, cavern.riskLevel)

	// initialise the lowest grid with maxint, with a one square border all
	// around containing zeros
	max := math.MaxInt

	rows := cavern.rows - 1
	cols := cavern.cols - 1
	start := cavern.cols + 1
	i := start
	for row := 1; row < rows; row++ {
		for col := 1; col < cols; col++ {
			lowest[i] = max
			i++
		}
		i += 2
	}

	walkLowestRiskLevel(start, 0, lowest, cavern)
	end := (cavern.rows - 1) * cavern.cols - 2

	//printGrid(cavern.rows, cavern.cols, lowest)

	return lowest[end] - cavern.riskLevel[start]
}

func walkLowestRiskLevel(pos, riskLevel int, lowest []int, cavern *Cavern) {
	lowestRiskLevel := lowest[pos]
	if lowestRiskLevel == 0 {
		return // off the map
	}

	riskLevel += cavern.riskLevel[pos]
	if riskLevel >= lowestRiskLevel {
		return // this is no better than previous
	}

	// We've found a better path to this square
	lowest[pos] = riskLevel

	h := 1
	v := cavern.cols

	walkLowestRiskLevel(pos + h, riskLevel, lowest, cavern) // right
	walkLowestRiskLevel(pos + v, riskLevel, lowest, cavern) // down
	walkLowestRiskLevel(pos - h, riskLevel, lowest, cavern) // left
	walkLowestRiskLevel(pos - v, riskLevel, lowest, cavern) // up
}

func parseCavern(lines []string) Cavern {
	// Leave a one square border around each side
	rows := len(lines) + 2
	cols := len(lines[0]) + 2
	riskLevel := make([]int, rows * cols)

	i := cols + 1 // {1,1}
	for _, line := range lines {
		for _, digit := range line {
			riskLevel[i] = int(digit - '0')
			i++
		}
		i += 2
	}

	return Cavern{rows:rows, cols:cols, riskLevel:riskLevel}
}

func printGrid(rows, cols int, cell []int) {
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			fmt.Printf("%3d ", cell[i])
			i++
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

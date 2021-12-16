package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
)

type Cavern struct {
	rows, cols int
	riskLevel []int	// this is (rows+2) * (cols+2) - sentinel border all around
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
	// initialise the lowest grid with maxint, with a one square border all
	// around containing zeros
	lowest := cavern.makeGrid(math.MaxInt)

	//printGrid(cavern.rows + 2, cavern.cols + 2, cavern.riskLevel)
	//printGrid(cavern.rows + 2, cavern.cols + 2, lowest)

	w := cavern.cols + 2
	h := cavern.rows + 2

	start := w + 1
	end := w * (h-1) - 2
	walkLowestRiskLevel(start, 0, lowest, cavern)

	//printGrid(cavern.rows + 2, cavern.cols + 2, lowest)

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
	v := cavern.cols + 2

	walkLowestRiskLevel(pos + h, riskLevel, lowest, cavern) // right
	walkLowestRiskLevel(pos + v, riskLevel, lowest, cavern) // down
	walkLowestRiskLevel(pos - h, riskLevel, lowest, cavern) // left
	walkLowestRiskLevel(pos - v, riskLevel, lowest, cavern) // up
}

func (c *Cavern) makeGrid(value int) []int {
	grid := make([]int, len(c.riskLevel))

	for row := 0; row < c.rows; row++ {
		for col := 0; col < c.cols; col++ {
			index := c.index(row, col)
			grid[index] = value
		}
	}

	return grid
}

func (c *Cavern) index(row, col int) int {
	// Take the borders into account here
	return (row + 1) * (c.cols + 2) + col + 1
}

func parseCavern(lines []string) Cavern {
	// Leave a one square border around each side
	rows := len(lines)
	cols := len(lines[0])
	riskLevel := make([]int, (rows + 2) * (cols + 2))

	cavern := Cavern{rows:rows, cols:cols, riskLevel:riskLevel}

	for row, line := range lines {
		for col, digit := range line {
			index := cavern.index(row, col)
			riskLevel[index] = int(digit - '0')
		}
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

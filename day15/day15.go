package main

import (
	"advent-of-code/aoc"
	"fmt"
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

	bigCavern := cavern.expand(5)
	fmt.Println(part1(&bigCavern))
}

func part1(cavern *Cavern) int {
	lowest := cavern.makeGrid(9999)

	//printGrid(cavern.rows + 2, cavern.cols + 2, cavern.riskLevel)
	//printGrid(cavern.rows + 2, cavern.cols + 2, lowest)

	posStack := make([]int, cavern.rows * cavern.cols)
	stackTop := 0
	w := cavern.cols + 2
	start := w + 1

	lowest[start] = 0
	posStack[stackTop] = start
	stackTop++

	hDelta := 1
	vDelta := cavern.cols + 2

	for stackTop > 0 {
		stackTop--
		pos := posStack[stackTop]
		pathRiskLevel := lowest[pos]

		p := pos + hDelta
		if nextRiskLevel := cavern.riskLevel[p]; nextRiskLevel > 0 {
			if riskLevel := pathRiskLevel + nextRiskLevel; riskLevel < lowest[p] {
				lowest[p] = riskLevel
				posStack[stackTop] = p
				stackTop++
			}
		}
		p = pos + vDelta
		if nextRiskLevel := cavern.riskLevel[p]; nextRiskLevel > 0 {
			if riskLevel := pathRiskLevel + nextRiskLevel; riskLevel < lowest[p] {
				lowest[p] = riskLevel
				posStack[stackTop] = p
				stackTop++
			}
		}
		p = pos - hDelta
		if nextRiskLevel := cavern.riskLevel[p]; nextRiskLevel > 0 {
			if riskLevel := pathRiskLevel + nextRiskLevel; riskLevel < lowest[p] {
				lowest[p] = riskLevel
				posStack[stackTop] = p
				stackTop++
			}
		}
		p = pos - vDelta
		if nextRiskLevel := cavern.riskLevel[p]; nextRiskLevel > 0 {
			if riskLevel := pathRiskLevel + nextRiskLevel; riskLevel < lowest[p] {
				lowest[p] = riskLevel
				posStack[stackTop] = p
				stackTop++
			}
		}
	}

	h := cavern.rows + 2
	end := w * (h-1) - 2
	return lowest[end]
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

func (c *Cavern) expand(factor int) Cavern {
	rows := c.rows * factor
	cols := c.cols * factor
	riskLevel := make([]int, (rows + 2) * (cols + 2))
	newCavern := Cavern{rows:rows, cols:cols, riskLevel:riskLevel}

	for row := 0; row < c.rows; row++ {
		for col := 0; col < c.cols; col++ {
			value := c.riskLevel[c.index(row, col)];
			for drow := 0; drow < factor; drow++ {
				for dcol := 0; dcol < factor; dcol++ {
					dvalue := value + drow + dcol
					if dvalue > 9 {
						dvalue -= 9
					}
					index := newCavern.index(drow * c.rows + row, dcol * c.cols + col)
					newCavern.riskLevel[index] = dvalue
				}
			}
		}
	}

	return newCavern
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
			fmt.Printf("%d ", cell[i])
			i++
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

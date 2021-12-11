package main

import (
	"advent-of-code/common"
	"fmt"
)

type Grid struct {
	rows, cols int
	energy     []int
}

type Point struct {
	row, col int
}

func main() {
	filename := common.GetFilename()
	grid := parseGrid(filename)

	fmt.Println(part1(grid, 100))
}

func part1(g Grid, steps int) int {
	//g.print()
	flashes := 0
	for step := 0; step < steps; step++ {
		flashes += g.doStep()
		//g.print()
	}
	return flashes
}

func (g* Grid) doStep() int {
	candidates := make([]Point, 0)

	// First part of step, all octopuses increase energy by one
	i := 0
	p := Point{row:0, col:0}
	for ; p.row < g.rows; p.row++ {
		for p.col = 0; p.col < g.cols; p.col++ {
			g.energy[i]++
			if g.energy[i] > 9 {
				candidates = append(candidates, p)
				//fmt.Println("Adding candidate(1): ", p)
			}
			i++
		}
	}

	flashing := 0
	isFlashing := make([]bool, g.rows * g.cols)
	for {
		remaining := len(candidates)
		if remaining == 0 {
			break
		}
		candidate := candidates[remaining-1]
		candidates = candidates[0:remaining-1]
		index, _ := g.index(&candidate)

		if isFlashing[index] {
			continue
		}
		isFlashing[index] = true
		flashing++

		for drow := -1; drow <= 1; drow++ {
			for dcol := -1; dcol <= 1; dcol++ {
				if drow == 0 && dcol == 0 {
					continue
				}

				neighbour := Point{ row: candidate.row + drow, col: candidate.col + dcol }
				if index, onGrid := g.index(&neighbour); onGrid {
					if !isFlashing[index] {
						g.energy[index]++
						if g.energy[index] > 9 {
							candidates = append(candidates, neighbour)
							//fmt.Println("Adding candidate(2): ", neighbour)
						}
					}
				}
			}
		}
	}

	i = 0
	p = Point{row:0, col:0}
	for ; p.row < g.rows; p.row++ {
		for p.col = 0; p.col < g.cols; p.col++ {
			if g.energy[i] > 9 {
				g.energy[i] = 0
			}
			i++
		}
	}

	return flashing
}

func (g* Grid) index(p* Point) (int, bool) {
	if p.row < 0 || p.col < 0 || p.row >= g.rows || p.col >= g.cols {
		return -1, false
	}
	return g.cols * p.row + p.col, true
}

func (g* Grid) print() {
	i := 0
	for row := 0; row < g.rows; row++ {
		for col := 0; col < g.cols; col++ {
			if g.energy[i] == 0 {
				fmt.Printf("*")
			} else {
				fmt.Printf("%d", g.energy[i])
			}
			i++
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func parseGrid(filename string) Grid {
	lines := common.GetInputLines(filename)
	rows := len(lines)
	cols := len(lines[0])
	grid := Grid{rows: rows, cols: cols, energy: make([]int, rows*cols)}

	index := 0
	for _, line := range lines {
		for _, digit := range line {
			grid.energy[index] = int(digit - '0')
			index++
		}
	}

	return grid
}

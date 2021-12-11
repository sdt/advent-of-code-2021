package main

import (
	"advent-of-code/aoc"
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
	filename := aoc.GetFilename()

	fmt.Println(part1(parseGrid(filename), 100))
	fmt.Println(part2(parseGrid(filename)))
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

func part2(g Grid) int {
	all := g.rows * g.cols
	for step := 0; ; step++ {
		flashes := g.doStep()
		if flashes == all {
			return step + 1
		}
	}
}

func (g *Grid) doStep() int {
	flashing := 0

	p := Point{row: 0, col: 0}
	for ; p.row < g.rows; p.row++ {
		for p.col = 0; p.col < g.cols; p.col++ {
			flashing += g.increaseEnergy(&p)
		}
	}

	i := 0
	p = Point{row: 0, col: 0}
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

func (g *Grid) increaseEnergy(p* Point) int {
	index, onGrid := g.index(p)
	if !onGrid {
		return 0
	}

	g.energy[index]++
	if g.energy[index] != 10 {
		return 0
	}

	total := 1
	for drow := -1; drow <= 1; drow++ {
		for dcol := -1; dcol <= 1; dcol++ {
			if drow == 0 && dcol == 0 {
				continue
			}

			neighbour := Point{row: p.row + drow, col: p.col + dcol}
			total += g.increaseEnergy(&neighbour)
		}
	}
	return total
}


func (g *Grid) index(p *Point) (int, bool) {
	if p.row < 0 || p.col < 0 || p.row >= g.rows || p.col >= g.cols {
		return -1, false
	}
	return g.cols*p.row + p.col, true
}

func (g *Grid) print() {
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
	lines := aoc.GetInputLines(filename)
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

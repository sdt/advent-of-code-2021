package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
)

type Cavern struct {
	rows, cols int
	riskLevel  []int // this is (rows+2) * (cols+2) - sentinel border all around
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
	lowest := cavern.makeGrid(math.MaxInt)

	//printGrid(cavern.rows + 2, cavern.cols + 2, cavern.riskLevel)
	//printGrid(cavern.rows + 2, cavern.cols + 2, lowest)

	w := cavern.cols + 2
	h := cavern.rows + 2
	start := w + 1
	end := w*(h-1) - 2

	pathQueue := MakePathQueue()
	pathQueue.Insert(MakePath(start, 0))

	lowest[start] = 0

	hDelta := 1
	vDelta := cavern.cols + 2
	deltas := [...]int{hDelta, vDelta, -hDelta, -vDelta}

	for {
		path, ok := pathQueue.Remove()
		if !ok {
			panic("All out of bubblegum!")
		}

		for _, delta := range deltas {
			p := path.pos + delta
			if nextRiskLevel := cavern.riskLevel[p]; nextRiskLevel > 0 {
				if riskLevel := path.cost + nextRiskLevel; riskLevel < lowest[p] {
					if p == end {
						return riskLevel
					}
					lowest[p] = riskLevel
					pathQueue.Insert(MakePath(p, riskLevel))
				}
			}
		}
	}

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
	return (row+1)*(c.cols+2) + col + 1
}

func (c *Cavern) expand(factor int) Cavern {
	rows := c.rows * factor
	cols := c.cols * factor
	riskLevel := make([]int, (rows+2)*(cols+2))
	newCavern := Cavern{rows: rows, cols: cols, riskLevel: riskLevel}

	for row := 0; row < c.rows; row++ {
		for col := 0; col < c.cols; col++ {
			value := c.riskLevel[c.index(row, col)]
			for drow := 0; drow < factor; drow++ {
				for dcol := 0; dcol < factor; dcol++ {
					dvalue := value + drow + dcol
					if dvalue > 9 {
						dvalue -= 9
					}
					index := newCavern.index(drow*c.rows+row, dcol*c.cols+col)
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
	riskLevel := make([]int, (rows+2)*(cols+2))

	cavern := Cavern{rows: rows, cols: cols, riskLevel: riskLevel}

	for row, line := range lines {
		for col, digit := range line {
			index := cavern.index(row, col)
			riskLevel[index] = int(digit - '0')
		}
	}

	return Cavern{rows: rows, cols: cols, riskLevel: riskLevel}
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

//------------------------------------------------------------------------------

type Path struct {
	pos, cost int
}

func MakePath(pos, cost int) Path {
	return Path{pos: pos, cost: cost}
}

type PathQueue struct {
	queue []Path
}

func MakePathQueue() PathQueue {
	return PathQueue{queue: make([]Path, 0)}
}

func (this *PathQueue) Insert(path Path) {
	this.queue = append(this.queue, path)

	// Upheap
	for child := len(this.queue) - 1; child > 0; {
		parent := (child - 1) / 2
		if this.queue[parent].cost <= path.cost {
			return
		}

		this.queue[child] = this.queue[parent]
		this.queue[parent] = path
		child = parent
	}
}

func (this *PathQueue) Remove() (Path, bool) {
	if len(this.queue) == 0 {
		return Path{}, false
	}

	ret := this.queue[0]
	size := len(this.queue) - 1

	if size > 0 {
		this.queue[0] = this.queue[size]

		// Downheap
		parent := 0
		for {
			lchild := parent*2 + 1
			if lchild >= size {
				break
			}
			rchild := parent*2 + 2

			var child int
			if rchild >= size {
				child = lchild
			} else if this.queue[lchild].cost <= this.queue[rchild].cost {
				child = lchild
			} else {
				child = rchild
			}

			if this.queue[parent].cost < this.queue[child].cost {
				break
			}

			this.queue[parent], this.queue[child] = this.queue[child], this.queue[parent]
			parent = child
		}
	}

	this.queue = this.queue[0:size]
	return ret, true
}

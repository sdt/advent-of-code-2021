package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Cavern struct {
	rows, cols int
	cell []int
}

type Position struct {
	row, col int
}

type Path struct {
	seen map[Position]bool
	riskLevel int
	pos Position
}

type PathQueue struct {
	heap []Path
}

func main() {
	filename := aoc.GetFilename()
	cavern := parseCavern(aoc.GetInputLines(filename))

	fmt.Println(part1(&cavern))
}

func part1(cavern *Cavern) int {
	queue := makePathQueue()
	queue.AddPath(makePath())

	for {
		path := queue.BestPath()

		if cavern.AtGoal(&path.pos) {
			return path.riskLevel
		}

		if next, ok := path.Extend(0, 1, cavern); ok {
			queue.AddPath(*next)
		}
		if next, ok := path.Extend(1, 0, cavern); ok {
			queue.AddPath(*next)
		}
		if next, ok := path.Extend(0, -1, cavern); ok {
			queue.AddPath(*next)
		}
		if next, ok := path.Extend(-1, 0, cavern); ok {
			queue.AddPath(*next)
		}
	}

	return 0
}

func parseCavern(lines []string) Cavern {
	rows := len(lines)
	cols := len(lines[0])
	cell := make([]int, rows * cols)

	i := 0
	for _, line := range lines {
		for _, digit := range line {
			cell[i] = int(digit - '0')
			i++
		}
	}

	return Cavern{rows:rows, cols:cols, cell:cell}
}

func makePath() Path {
	seen := make(map[Position]bool)
	riskLevel := 0
	pos := Position{row:0, col:0}
	seen[pos] = true

	return Path{seen:seen, riskLevel:riskLevel, pos:pos}
}

func makePathQueue() PathQueue {
	return PathQueue{heap:make([]Path, 0)}
}

func (c *Cavern) InCavern(p *Position) bool {
	return p.row >= 0 && p.col >= 0 && p.row < c.rows && p.col < c.cols
}

func (c *Cavern) AtGoal(p* Position) bool {
	return (p.row + 1 == c.rows) && (p.col + 1 == c.cols)
}

func (c *Cavern) GetRiskLevel(p *Position) int {
	return c.cell[p.row * c.cols + p.col]
}

func (q *PathQueue) AddPath(p Path) {
	q.heap = append(q.heap, p)

	// up heap
	child := len(q.heap) - 1
	for child > 0 {
		parent := (child - 1) / 2

		if q.heap[parent].riskLevel < q.heap[child].riskLevel {
			return
		}

		q.heap[parent], q.heap[child] = q.heap[child], q.heap[parent]
		child = parent
	}
}

func (q *PathQueue) BestPath() Path {
	best := q.heap[0]

	size := len(q.heap) - 1

	if size > 0 {
		q.heap[0] = q.heap[size]

		// down heap
		parent := 0
		for {
			lchild := parent * 2 + 1
			if lchild >= size {
				break
			}
			rchild := parent * 2 + 2
			var child int
			if rchild >= size {
				child = lchild
			} else if q.heap[lchild].riskLevel <= q.heap[rchild].riskLevel {
				child = lchild
			} else {
				child = rchild
			}

			if q.heap[parent].riskLevel < q.heap[child].riskLevel {
				break
			}

			q.heap[parent], q.heap[child] = q.heap[child], q.heap[parent]
			parent = child
		}
	}

	q.heap = q.heap[0:size]
	return best
}

func (p *Path) Extend(x, y int, cavern *Cavern) (*Path, bool) {
	pos := Position{p.pos.row + y, p.pos.col + x}
	if !cavern.InCavern(&pos) || p.seen[pos] {
		return nil, false
	}

	if p.seen[pos] {
		return nil, false
	}

	seen := make(map[Position]bool, len(p.seen) + 1)
	for k,v := range p.seen {
		seen[k] = v
	}
	seen[pos] = true

	ret := Path{ seen:seen, riskLevel:p.riskLevel + cavern.GetRiskLevel(&pos), pos:pos }

	return &ret, true
}

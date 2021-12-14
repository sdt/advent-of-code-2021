package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

type Point struct {
	p [2]int
}

type Fold struct {
	dir int
	pos int
}

func main() {
	points, folds := getInput(aoc.GetFilename())

	fmt.Println(part1(points, folds))
	part2(points, folds)
}

func part1(points []Point, folds []Fold) int {
	unique := make(map[Point]bool)

	for _, in := range points {
		out := folds[0].Apply(in)
		unique[out] = true
	}
	return len(unique)
}

func part2(points []Point, folds []Fold) {
	dot := make(map[Point]bool)
	max := makePoint(0, 0)

	for _, p := range points {
		for _, f := range folds {
			p = f.Apply(p)
		}

		dot[p] = true
		for i := 0; i < 2; i++ {
			if p.p[i] > max.p[i] {
				max.p[i] = p.p[i]
			}
		}
	}

	for y := 0; y <= max.p[1]; y++ {
		for x := 0; x <= max.p[0]; x++ {
			p := makePoint(x, y)
			if dot[p] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func makePoint(x, y int) Point {
	return Point{p: [2]int{x, y}}
}

func parsePoint(line string) Point {
	// 2,15
	words := strings.Split(line, ",")
	coords := aoc.ParseInts(words)
	return makePoint(coords[0], coords[1])
}

func parseFold(line string) Fold {
	// fold along x=3
	words := strings.Split(line, " ")
	parts := strings.Split(words[2], "=")

	dir := 0
	if parts[0] == "y" {
		dir = 1
	}

	pos := aoc.ParseInt(parts[1])
	return Fold{dir: dir, pos: pos}
}

func (this *Fold) Apply(p Point) Point {
	if p.p[this.dir] > this.pos {
		p.p[this.dir] = 2*this.pos - p.p[this.dir]
	}
	return p
}

func getInput(filename string) ([]Point, []Fold) {
	lines := aoc.GetInputLines(filename)
	points := make([]Point, 0)
	folds := make([]Fold, 0)

	for _, line := range lines {
		if strings.Contains(line, ",") {
			points = append(points, parsePoint(line))
		} else if strings.Contains(line, "=") {
			folds = append(folds, parseFold(line))
		}
	}

	return points, folds
}

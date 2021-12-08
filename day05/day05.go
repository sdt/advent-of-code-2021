package main

import (
	"advent-of-code/common"
	"fmt"
	"regexp"
)

var LineRegex = regexp.MustCompile("^(\\d+),(\\d+) -> (\\d+),(\\d+)$")

type Point struct {
	x, y int
}

type Line struct {
	p [2]Point
}

type Ocean map[Point]int

func main() {
	filename := common.GetFilename()
	lines := parse(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []Line) int {
	ocean := make(Ocean)

	for _, line := range lines {
		line.trace(ocean, false)
	}

	total := 0
	for _, count := range ocean {
		if count > 1 {
			total++
		}
	}
	//fmt.Println(ocean)

	return total
}

func part2(lines []Line) int {
	ocean := make(Ocean)

	for _, line := range lines {
		line.trace(ocean, true)
	}

	total := 0
	for _, count := range ocean {
		if count > 1 {
			total++
		}
	}
	//fmt.Println(ocean)

	return total
}

func (line *Line) trace(ocean Ocean, withDiag bool) {
	dx, nx := signMag(line.p[1].x - line.p[0].x)
	dy, ny := signMag(line.p[1].y - line.p[0].y)
	p := line.p[0]

	if dy == 0 {
		for i := 0; i <= nx; i++ {
			ocean[p]++
			p.x += dx
		}
	} else if dx == 0 {
		for i := 0; i <= ny; i++ {
			ocean[p]++
			p.y += dy
		}
	} else if withDiag {
		for i := 0; i <= nx; i++ {
			ocean[p]++
			p.x += dx
			p.y += dy
		}
	}
}

func signMag(value int) (int, int) {
	if value > 0 {
		return 1, value
	}
	if value < 0 {
		return -1, -value
	}
	return 0, 0
}

func parse(filename string) []Line {
	input := common.GetInputLines(filename)
	lines := make([]Line, len(input))

	for i, line := range input {
		lines[i] = parseLine(line)
	}
	return lines
}

func parseLine(line string) Line {
	matches := LineRegex.FindStringSubmatch(line)
	numbers := make([]int, 4)
	for i := 0; i < 4; i++ {
		numbers[i] = common.ParseInt(matches[i+1])
	}
	return Line{p: [2]Point{Point{x: numbers[0], y: numbers[1]}, Point{x: numbers[2], y: numbers[3]}}}
}

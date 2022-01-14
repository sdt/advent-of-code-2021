package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

type Coord int64

type Point struct {
	x, y, z Coord
}

type Dist uint64

func abs(x Coord) Dist {
	if x < 0 {
		x = -x
	}
	return Dist(x)
}

func (this Coord) Distance(that Coord) Dist {
	return abs(this - that)
}

func (this Point) ManhattanDistance(that Point) Dist {
	return this.x.Distance(that.x) +
	       this.y.Distance(that.y) +
	       this.z.Distance(that.z)
}

func (this Point) ManhattanMagnitude() Dist {
	return abs(this.x) + abs(this.y) + abs(this.z)
}

type Scanner struct {
	pos Point
	r Dist
}

func (this Scanner) Contains(p Point) bool {
	return this.pos.ManhattanDistance(p) <= this.r
}

func (this Scanner) Corners() []Point {
	r := Coord(this.r)
	r2 := r / 3
	return []Point{
		Point{this.pos.x - r, this.pos.y, this.pos.z},
		Point{this.pos.x + r, this.pos.y, this.pos.z},
		Point{this.pos.x, this.pos.y - r, this.pos.z},
		Point{this.pos.x, this.pos.y + r, this.pos.z},
		Point{this.pos.x, this.pos.y, this.pos.z - r},
		Point{this.pos.x, this.pos.y, this.pos.z + r},

		Point{this.pos.x + r2, this.pos.y + r2, this.pos.z + r2},
		Point{this.pos.x - r2, this.pos.y + r2, this.pos.z + r2},
		Point{this.pos.x + r2, this.pos.y - r2, this.pos.z + r2},
		Point{this.pos.x - r2, this.pos.y - r2, this.pos.z + r2},
		Point{this.pos.x + r2, this.pos.y + r2, this.pos.z - r2},
		Point{this.pos.x - r2, this.pos.y + r2, this.pos.z - r2},
		Point{this.pos.x + r2, this.pos.y - r2, this.pos.z - r2},
		Point{this.pos.x - r2, this.pos.y - r2, this.pos.z - r2},
	}
}

var ScannerRegex = regexp.MustCompile("^pos=<(-?\\d+),(-?\\d+),(-?\\d+)>, r=(\\d+)$")

func NewScanner(line string) Scanner {
	// pos=<50,50,50>, r=200
	matches := ScannerRegex.FindStringSubmatch(line)

	v := aoc.ParseInts(matches[1:])

	pos := Point{Coord(v[0]),Coord(v[1]),Coord(v[2])}
	r := Dist(v[3])

	return Scanner{pos, r}
}

func ParseInput(filename string) []Scanner {
	lines := aoc.GetInputLines(filename)

	scanners := make([]Scanner, 0, len(lines))
	for _, line := range lines {
		scanners = append(scanners, NewScanner(line))
	}

	return scanners
}

func countScanners(scanners []Scanner, point Point) int {
	count := 0
	for _, scanner := range scanners {
		if scanner.Contains(point) {
			count++
		}
	}
	return count
}

func main() {
	scanners := ParseInput(aoc.GetFilename())

	/*
	max := 0

	hist := make(map[Point]int)

	for _, scanner := range scanners {
		for _, corner := range scanner.Corners() {
			if _, found := hist[corner]; found {
				continue
			}
			count := countScanners(scanners, corner)
			if count >= max {
				max = count
				fmt.Printf("Point %v touches %d scanners NEW MAX\n", corner, count)
			} else {
				fmt.Printf("Point %v touches %d scanners\n", corner, count)
			}
			hist[corner] = count
		}
	}

	for point, count := range hist {
		if count == max {
			fmt.Printf("Point %v touches %d scanners: dist=%d\n", point, count, point.ManhattanMagnitude())
		}
	}

	*/

	// 95540990 is too low		15434918 32581102 47524970
	// 96142944 is too high		23128222 28049091 44965631

	start := Point{23128222, 28049091, 44965631}
	dx := Coord(1)
	for {
		newPos := start
		newPos.x -= dx
		newCount := countScanners(scanners, newPos)
		if newCount < 876 {
			fmt.Printf("%v drops to %d\n", newPos, newCount)
			newPos.x += 1
			fmt.Println(newPos.ManhattanMagnitude())
			break
		}
		dx++	// 95540990 too low (already seen this)
	}
	dy := Coord(1)
	for {
		newPos := start
		newPos.y -= dy
		newCount := countScanners(scanners, newPos)
		if newCount < 876 {
			fmt.Printf("%v drops to %d\n", newPos, newCount)
			newPos.y += 1
			fmt.Println(newPos.ManhattanMagnitude())
			break
		}
		dy++	// 95540990 too low (already seen this)
	}
	dz := Coord(1)
	for {
		newPos := start
		newPos.z -= dz
		newCount := countScanners(scanners, newPos)
		if newCount < 876 {
			fmt.Printf("%v drops to %d\n", newPos, newCount)
			newPos.z += 1
			fmt.Println(newPos.ManhattanMagnitude())
			break
		}
		dz++	// 96142944 too high (already seen this)
	}
	dxyz := Coord(1)
	for {
		newPos := start
		newPos.x -= dxyz
		newPos.y -= dxyz
		newPos.z -= dxyz
		newCount := countScanners(scanners, newPos)
		if newCount < 876 {
			fmt.Printf("%v drops to %d\n", newPos, newCount)
			newPos.x += 1
			newPos.y += 1
			newPos.z += 1
			fmt.Println(newPos.ManhattanMagnitude())
			break
		}
		dxyz++	// 95540991 too low
	}
}

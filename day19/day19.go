package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	scanners := ParseInput(aoc.GetFilename())
	transforms := CreateScannerTransforms(scanners)

	fmt.Println(part1(scanners, transforms))
	fmt.Println(part2(transforms))
}

//------------------------------------------------------------------------------

func part1(scanners []Scanner, transforms []Transform) int {
	type UniquePoints map[Point]bool
	unique := make(UniquePoints)

	for i, scanner := range scanners {
		transform := transforms[i]

		for _, from := range scanner {
			to := transform(from)
			unique[to] = true
		}
	}

	return len(unique)
}

func part2(transforms []Transform) int {

	origin := MakePoint(0, 0, 0)
	pos := make([]Point, len(transforms))
	for i, transform := range transforms {
		pos[i] = transform(origin)
	}

	best := 0
	for _, p0 := range pos {
		for _, p1 := range pos {
			dist := p0.ManhattanDistance(p1)
			if dist > best {
				best = dist
			}
		}
	}
	return best
}

//------------------------------------------------------------------------------

type Scanner []Point

var LineRegex = regexp.MustCompile("scanner (\\d+)")

func ParseInput(filename string) []Scanner {
	lines := aoc.GetInputLines(filename)

	scanner := make([]Scanner, 0)
	var index int

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if matches := LineRegex.FindStringSubmatch(line); len(matches) == 2 {
			index = aoc.ParseInt(matches[1])
			scanner = append(scanner, make([]Point, 0))
		} else {
			scanner[index] = append(scanner[index], ParsePoint(line))
		}
	}

	return scanner
}

func CompareScanners(s0, s1 Scanner) (Transform, bool) {
	type Hist map[Point]int

	for _, rot := range rotations {
		hist := make(Hist)
		for _, p1 := range s1 {
			rp1 := rot(p1)
			for _, p0 := range s0 {
				p := p0.Sub(rp1)
				if count, found := hist[p]; found {
					if count == 11 {
						return rot.MakeTransform(p), true
					}
					hist[p] = count + 1
				} else {
					hist[p] = 1
				}
			}
		}
	}
	return nil, false
}

func CreateScannerTransforms(scanners []Scanner) []Transform {
	scannerTransform := make([]Transform, len(scanners))
	scannerTransform[0] = Identity

	unmapped := make([]int, len(scanners)-1)
	for i := 0; i < len(scanners)-1; i++ {
		unmapped[i] = i + 1
	}

	candidates := []int{0}

	for len(unmapped) > 0 {
		//fmt.Printf("Unmapped=%d candidates=%d\n", len(unmapped), len(candidates))
		from := candidates[0]
		candidates = candidates[1:]

		stillUnmapped := make([]int, 0, len(unmapped))

		for _, to := range unmapped {
			if transform, found := CompareScanners(scanners[from], scanners[to]); found {
				candidates = append(candidates, to)
				//fmt.Printf("Scanner %d -> %d\n", from, to)
				scannerTransform[to] = scannerTransform[from].Extend(transform)
			} else {
				stillUnmapped = append(stillUnmapped, to)
			}
		}

		unmapped = stillUnmapped
	}

	return scannerTransform
}

//------------------------------------------------------------------------------

type Point struct {
	x, y, z int
}

func MakePoint(x, y, z int) Point {
	return Point{x: x, y: y, z: z}
}

func ParsePoint(line string) Point {
	words := strings.Split(line, ",")
	numbers := aoc.ParseInts(words)
	return MakePoint(numbers[0], numbers[1], numbers[2])
}

func (p Point) Add(q Point) Point {
	return MakePoint(p.x+q.x, p.y+q.y, p.z+q.z)
}

func (p Point) Sub(q Point) Point {
	return MakePoint(p.x-q.x, p.y-q.y, p.z-q.z)
}

func (p Point) ManhattanDistance(q Point) int {
	return abs(p.x-q.x) + abs(p.y-q.y) + abs(p.z-q.z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Rotation func(Point) Point
type Transform func(Point) Point

func Identity(p Point) Point {
	return p
}

func (rot Rotation) MakeTransform(offset Point) Transform {
	return func(p Point) Point {
		return rot(p).Add(offset)
	}
}

func (t1 Transform) Extend(t2 Transform) Transform {
	return func(p Point) Point {
		return t1(t2(p))
	}
}

var rotations [24]Rotation = [...]Rotation{
	func(p Point) Point { return MakePoint(+p.x, +p.y, +p.z) },
	func(p Point) Point { return MakePoint(+p.y, -p.x, +p.z) },
	func(p Point) Point { return MakePoint(-p.x, -p.y, +p.z) },
	func(p Point) Point { return MakePoint(-p.y, +p.x, +p.z) },

	func(p Point) Point { return MakePoint(-p.x, +p.y, -p.z) },
	func(p Point) Point { return MakePoint(-p.y, -p.x, -p.z) },
	func(p Point) Point { return MakePoint(+p.x, -p.y, -p.z) },
	func(p Point) Point { return MakePoint(+p.y, +p.x, -p.z) },

	func(p Point) Point { return MakePoint(-p.y, -p.z, +p.x) },
	func(p Point) Point { return MakePoint(+p.z, -p.y, +p.x) },
	func(p Point) Point { return MakePoint(+p.y, +p.z, +p.x) },
	func(p Point) Point { return MakePoint(-p.z, +p.y, +p.x) },

	func(p Point) Point { return MakePoint(-p.y, +p.z, -p.x) },
	func(p Point) Point { return MakePoint(-p.z, -p.y, -p.x) },
	func(p Point) Point { return MakePoint(+p.y, -p.z, -p.x) },
	func(p Point) Point { return MakePoint(+p.z, +p.y, -p.x) },

	func(p Point) Point { return MakePoint(-p.x, +p.z, +p.y) },
	func(p Point) Point { return MakePoint(-p.z, -p.x, +p.y) },
	func(p Point) Point { return MakePoint(+p.x, -p.z, +p.y) },
	func(p Point) Point { return MakePoint(+p.z, +p.x, +p.y) },

	func(p Point) Point { return MakePoint(+p.x, +p.z, -p.y) },
	func(p Point) Point { return MakePoint(-p.z, +p.x, -p.y) },
	func(p Point) Point { return MakePoint(-p.x, -p.z, -p.y) },
	func(p Point) Point { return MakePoint(+p.z, -p.x, -p.y) },
}

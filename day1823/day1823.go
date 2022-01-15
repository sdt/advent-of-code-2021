package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

//------------------------------------------------------------------------------

type Int int64

func (this Int) Abs() Int {
	if this < 0 {
		return -this
	}
	return this
}

//------------------------------------------------------------------------------

type Vec3 struct {
	x, y, z Int
}

func (this Vec3) Sub(that Vec3) Vec3 {
	return Vec3{this.x - that.x, this.y - that.y, this.z - that.z}
}

func (this Vec3) Add(that Vec3) Vec3 {
	return Vec3{this.x + that.x, this.y + that.y, this.z + that.z}
}

func (this Vec3) Dot(that Vec3) Int {
	return this.x * that.x + this.y * that.y + this.z * that.z
}

func (this Vec3) ManhattanMagnitude() Int {
	return this.x.Abs() + this.y.Abs() + this.z.Abs()
}

func (this Vec3) ManhattanDistance(that Vec3) Int {
	return this.Sub(that).ManhattanMagnitude()
}

func (this Vec3) Orientations() []Vec3 {
	x, y, z := this.x, this.y, this.z

	return []Vec3{
		Vec3{+x, +y, +z},
		Vec3{+x, +y, -z},
		Vec3{+x, -y, +z},
		Vec3{+x, -y, -z},
		Vec3{-x, +y, +z},
		Vec3{-x, +y, -z},
		Vec3{-x, -y, +z},
		Vec3{-x, -y, -z},
	}
}

//------------------------------------------------------------------------------

type Cube struct {
	pos Vec3
	size Int
}

func (this *Cube) Corners() []Vec3 {
	x0, x1 := this.pos.x, this.pos.x + this.size - 1
	y0, y1 := this.pos.y, this.pos.y + this.size - 1
	z0, z1 := this.pos.z, this.pos.z + this.size - 1
	return []Vec3{
		Vec3{x0, y0, z0},
		Vec3{x0, y0, z1},
		Vec3{x0, y1, z0},
		Vec3{x0, y1, z1},
		Vec3{x1, y0, z0},
		Vec3{x1, y0, z1},
		Vec3{x1, y1, z0},
		Vec3{x1, y1, z1},
	}
}

func (this *Cube) IsFullyInside(scanner *Scanner) bool {
	for _, corner := range this.Corners() {
		if !scanner.Contains(corner) {
			return false
		}
	}
	return true
}

func (this *Cube) IsFullyOutside(scanner *Scanner) bool {
	// Scanner is made of eight planes. For a cube to be fully distinct from
	// a scanner, all eight corners of the cube must be outside one of those
	// planes.

	// Transform the scanner to the origin, and the cube relative to that.
	corners := this.Corners()
	for i := 0; i < len(corners); i++ {
		corners[i] = corners[i].Sub(scanner.pos)
	}

	unit := Vec3{1, 1, 1}
	planes := unit.Orientations()

	for _, plane := range planes {
		inside := false
		for _, corner := range corners {
			if plane.Dot(corner) <= scanner.r {
				inside = true
				break
			}
		}
		if !inside {
			return true
		}
	}
	return false
}

func (this *Cube) Split() []Cube {
	if this.size == 1 {
		panic("Trying to split unit cube")
	}

	size := this.size / 2
	positions := (&Cube{this.pos, size + 1}).Corners()
	cubes := make([]Cube, 0, 8)

	for _, pos := range positions {
		cubes = append(cubes, Cube{pos, size})
	}

	return cubes
}

//------------------------------------------------------------------------------

type Scanner struct {
	pos Vec3
	r Int
}

func (this *Scanner) Contains(p Vec3) bool {
	return this.pos.ManhattanDistance(p) <= this.r
}

func (this *Scanner) Corners() []Vec3 {
	r := this.r
	return []Vec3{
		Vec3{this.pos.x - r, this.pos.y, this.pos.z},
		Vec3{this.pos.x + r, this.pos.y, this.pos.z},
		Vec3{this.pos.x, this.pos.y - r, this.pos.z},
		Vec3{this.pos.x, this.pos.y + r, this.pos.z},
		Vec3{this.pos.x, this.pos.y, this.pos.z - r},
		Vec3{this.pos.x, this.pos.y, this.pos.z + r},
	}
}

var ScannerRegex = regexp.MustCompile("^pos=<(-?\\d+),(-?\\d+),(-?\\d+)>, r=(\\d+)$")

func NewScanner(line string) Scanner {
	// pos=<50,50,50>, r=200
	matches := ScannerRegex.FindStringSubmatch(line)

	v := aoc.ParseInts(matches[1:])

	pos := Vec3{Int(v[0]),Int(v[1]),Int(v[2])}
	r := Int(v[3])

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

func countScanners(scanners []Scanner, point Vec3) int {
	count := 0
	for _, scanner := range scanners {
		if scanner.Contains(point) {
			count++
		}
	}
	return count
}

//------------------------------------------------------------------------------

func makeInitialCube(scanners []Scanner) Cube {
	size := Int(1)

	for _, scanner := range scanners {
		corners := scanner.Corners()
		for _, corner := range corners {
			for c := corner.x.Abs(); c > size; size *= 2 { }
			for c := corner.y.Abs(); c > size; size *= 2 { }
			for c := corner.z.Abs(); c > size; size *= 2 { }
		}
	}

	return Cube{Vec3{-size, -size, -size}, size*2}
}

//------------------------------------------------------------------------------

func testCube(cube *Cube, scanner *Scanner) {
	if cube.IsFullyInside(scanner) {
		fmt.Println("****", cube, "is fully inside", scanner)
		return
	}

	if cube.IsFullyOutside(scanner) {
		//fmt.Println(cube, "is fully outside", scanner)
		return
	}

	subCubes := cube.Split()
	for _, subCube := range subCubes {
		testCube(&subCube, scanner)
	}
}

//------------------------------------------------------------------------------

func main() {
	scanners := ParseInput(aoc.GetFilename())

	cube := makeInitialCube(scanners)
	fmt.Println(cube)

	testCube(&cube, &scanners[4])
}

//------------------------------------------------------------------------------

type Heap interface {
	IsHigherPriority(parent, child int) bool
	IsValid(i int) bool
	Swap(i, j int)
}

func upheap(heap Heap, child int) {
	for child > 0 {
		parent := (child - 1) / 2

		if heap.IsHigherPriority(parent, child) {
			return
		}

		heap.Swap(parent, child)
		child = parent
	}
}

func downheap(heap Heap, parent int) {
	for {
		lchild := parent * 2 + 1
		if !heap.IsValid(lchild) {
			return
		}

		rchild := parent * 2 + 2
		var child int
		if !heap.IsValid(rchild) || heap.IsHigherPriority(lchild, rchild) {
			child = lchild
		} else {
			child = rchild
		}

		if heap.IsHigherPriority(parent, child) {
			return
		}

		heap.Swap(parent, child)
		parent = child
	}
}

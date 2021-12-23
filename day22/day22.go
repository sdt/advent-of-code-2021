package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
	"sort"
)

func main() {
	steps := ParseInput(aoc.GetFilename())

	fmt.Println(part1(MakeCube(50), steps))
	fmt.Println(part2(steps))
}

func part1(bounds Cuboid, steps []Step) int {
	cuboids := make([]Cuboid, 0)

	for _, step := range steps {
		if step.cuboid.Overlaps(&bounds) {
			cuboids = SplitCuboids(cuboids, step.cuboid)
			if step.isOn {
				cuboids = append(cuboids, step.cuboid)
			}
		}
	}

	total := 0
	for _, cuboid := range cuboids {
		total += cuboid.Size()
	}
	return total
}

func part2(steps []Step) int {
	cuboids := make([]Cuboid, 0)

	for _, step := range steps {
		cuboids = SplitCuboids(cuboids, step.cuboid)
		if step.isOn {
			cuboids = append(cuboids, step.cuboid)
		}
	}

	total := 0
	for _, cuboid := range cuboids {
		total += cuboid.Size()
	}
	return total
}

func SplitCuboids(in []Cuboid, splitter Cuboid) []Cuboid {
	out := make([]Cuboid, 0)

	// Apply the splitter turn-off cuboid to each input fragments in turn.
	// If the splitter cuboid overlaps the input cuboid, split the input cuboid,
	// otherwise keep the entire input cuboid. Leave the splitter untouched.

	for _, inputCuboid := range in {
		if splitter.Overlaps(&inputCuboid) {
			fragments := splitter.Split(&inputCuboid)
			out = append(out, fragments...)
		} else {
			out = append(out, inputCuboid)
		}
	}

	return out
}

type Point struct {
	v [3]int // x y z
}

func MakePoint(x, y, z int) Point {
	return Point{v: [...]int{x, y, z}}
}

type Cuboid struct {
	min, max Point
}

type Step struct {
	cuboid Cuboid
	isOn   bool
}

func (this Step) GoString() string {
	var onOff string
	if this.isOn {
		onOff = "on"
	} else {
		onOff = "off"
	}

	return fmt.Sprintf("Turn %s %#v", onOff, this.cuboid)
}

func MakeCube(size int) Cuboid {
	min, max := -size, size+1
	return MakeCuboid(min, max, min, max, min, max)
}

func MakeCuboid(xmin, xmax, ymin, ymax, zmin, zmax int) Cuboid {
	return Cuboid{min: MakePoint(xmin, ymin, zmin),
		max: MakePoint(xmax, ymax, zmax)}
}

func (this *Cuboid) IsValid() bool {
	for i := 0; i < 3; i++ {
		if this.min.v[i] >= this.max.v[i] {
			return false
		}
	}
	return true
}

func (this *Cuboid) Overlaps(that *Cuboid) bool {
	for i := 0; i < 3; i++ {
		if this.min.v[i] >= that.max.v[i] {
			return false
		}
		if that.min.v[i] >= this.max.v[i] {
			return false
		}
	}

	return true
}

func (this *Cuboid) Split(that *Cuboid) []Cuboid {
	var v [3][4]int

	for i := 0; i < 3; i++ {
		v[i][0] = this.min.v[i]
		v[i][1] = this.max.v[i]
		v[i][2] = that.min.v[i]
		v[i][3] = that.max.v[i]

		sort.Sort(sort.IntSlice(v[i][0:]))
	}
	//fmt.Println(corner)

	fragments := make([]Cuboid, 0)
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			for z := 0; z < 3; z++ {
				fragment := MakeCuboid(
					v[0][x], v[0][x+1],
					v[1][y], v[1][y+1],
					v[2][z], v[2][z+1],
				)

				if !fragment.IsValid() {
					continue // skip zero-volume fragments
				}

				// Keep only the fragments that touch that, but not this
				if !fragment.Overlaps(this) && fragment.Overlaps(that) {
					fragments = append(fragments, fragment)
				}
			}
		}
	}

	return fragments

}

func (this *Cuboid) Combine(that *Cuboid, doAdd bool) []Cuboid {
	var corner [3][4]int

	for i := 0; i < 3; i++ {
		corner[i][0] = this.min.v[i]
		corner[i][1] = this.max.v[i]
		corner[i][2] = that.min.v[i]
		corner[i][3] = that.max.v[i]

		sort.Sort(sort.IntSlice(corner[i][0:]))
	}
	//fmt.Println(corner)

	subCuboids := make([]Cuboid, 0)
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			for z := 0; z < 3; z++ {
				subCuboid := MakeCuboid(
					corner[0][x], corner[0][x+1],
					corner[1][y], corner[1][y+1],
					corner[2][z], corner[2][z+1],
				)

				if !subCuboid.IsValid() {
					continue
				}

				if doAdd {
					if subCuboid.Overlaps(this) || subCuboid.Overlaps(that) {
						subCuboids = append(subCuboids, subCuboid)
					}
				} else {
					if subCuboid.Overlaps(this) && !subCuboid.Overlaps(that) {
						subCuboids = append(subCuboids, subCuboid)
					}
				}
			}
		}
	}

	return subCuboids
}

func (this *Cuboid) Size() int {
	size := 1
	for i := 0; i < 3; i++ {
		size *= this.max.v[i] - this.min.v[i]
	}
	return size
}

func (this Cuboid) GoString() string {
	return fmt.Sprintf("x=%d..%d,y=%d..%d,z=%d..%d",
		this.min.v[0], this.max.v[0]-1,
		this.min.v[1], this.max.v[1]-1,
		this.min.v[2], this.max.v[2]-1)
}

func ParseInput(filename string) []Step {
	lines := aoc.GetInputLines(filename)

	steps := make([]Step, len(lines))

	for i, line := range lines {
		steps[i] = ParseStep(line)
	}

	return steps
}

var StepRegex = regexp.MustCompile("^(on|off) x=(-?\\d+)\\.\\.(-?\\d+),y=(-?\\d+)\\.\\.(-?\\d+),z=(-?\\d+)\\.\\.(-?\\d+)$")

func ParseStep(line string) Step {
	matches := StepRegex.FindStringSubmatch(line)
	isOn := (matches[1] == "on")

	v := make([]int, 6)
	for i := 0; i < 6; i++ {
		v[i] = aoc.ParseInt(matches[i+2])
	}

	return Step{cuboid: MakeCuboid(v[0], v[1]+1, v[2], v[3]+1, v[4], v[5]+1), isOn: isOn}
}

package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

func main() {
	steps := ParseInput(aoc.GetFilename())

	fmt.Println(part1(MakeCube(50), steps))
}

func part1(bounds Cuboid, steps []Step) int {
	smallSteps := make([]Step, 0)
	for _, step := range steps {
		if step.cuboid.Overlaps(&bounds) {
			smallSteps = append(smallSteps, step)
		}
	}

	world := make(map[Point]bool)

	for _, step := range smallSteps {
		c := &step.cuboid
		for z, zmax := c.min.p[2], c.max.p[2]; z <= zmax; z++ {
			for y, ymax := c.min.p[1], c.max.p[1]; y <= ymax; y++ {
				for x, xmax := c.min.p[0], c.max.p[0]; x <= xmax; x++ {
					p := Point{p: [...]int{x, y, z}}
					if step.isOn {
						world[p] = true
					} else {
						delete(world, p)
					}
				}
			}
		}
	}

	return len(world)
}

type Point struct {
	p [3]int	// x y z
}

type Cuboid struct {
	min, max Point
}

type Step struct {
	cuboid Cuboid
	isOn bool
}

func MakeCube(size int) Cuboid {
	return Cuboid{ min: Point{p: [...]int{-size, -size, -size}},
				   max: Point{p: [...]int{+size, +size, +size}} };
}

func (this *Cuboid) Overlaps(that *Cuboid) bool {
	for i := 0; i < 3; i++ {
		if this.min.p[i] > that.max.p[i] {
			return false
		}
		if that.min.p[i] > this.max.p[i] {
			return false
		}
	}

	return true
}

func ParseInput(filename string) []Step {
	lines := aoc.GetInputLines(filename);

	steps := make([]Step, len(lines))

	for i, line := range lines {
		steps[i] = ParseStep(line)
	}

	return steps
}

var StepRegex = regexp.MustCompile("^(on|off) x=(-?\\d+)\\.\\.(-?\\d+),y=(-?\\d+)\\.\\.(-?\\d+),z=(-?\\d+)\\.\\.(-?\\d+)$");

func ParseStep(line string) Step {
	matches := StepRegex.FindStringSubmatch(line)
	isOn := (matches[1] == "on");

	values := make([]int, 6)
	for i := 0; i < 6; i++ {
		values[i] = aoc.ParseInt(matches[i + 2])
	}

	return Step{
		cuboid: Cuboid{
			min: Point{ p: [...]int{ values[0], values[2], values[4] } },
			max: Point{ p: [...]int{ values[1], values[3], values[5] } } }, isOn: isOn }
}

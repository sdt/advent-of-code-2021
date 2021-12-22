package main

import (
	"advent-of-code/aoc"
	"fmt"
	"regexp"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)
	target := ParseTarget(lines[0])

	fmt.Println(part1(&target))
	fmt.Println(part2(&target))
}

//------------------------------------------------------------------------------

func part1(t *Target) int {
	mindx := MinimumDX(t)

	highest := 0

	for dx := mindx; dx < mindx + 10; dx++ {
		hits := 0
		for dy := -450; dy < 450; dy++ {
			p := MakeProjectile(dx, dy)
			hit, maxy := p.FireAt(t)
			//fmt.Printf("%d,%d -> %v,%d\n", dx, dy, hit, maxy)
			if hit {
				hits++
				if maxy > highest {
					highest = maxy
				}
			}
		}
		if hits == 0 {
			return highest
		}
	}

	return highest
}

//------------------------------------------------------------------------------

func part2(t *Target) int {
	mindx := MinimumDX(t)

	totalHits := 0

	for dx := mindx-2; dx < 69+2; dx++ {
		hits := 0
		for dy := -1550; dy < 1550; dy++ {
			p := MakeProjectile(dx, dy)
			hit, _ := p.FireAt(t)
			//fmt.Printf("%d,%d -> %v,%d\n", dx, dy, hit, maxy)
			if hit {
				hits++
			}
		}
		totalHits += hits
	}

	return totalHits
}

//------------------------------------------------------------------------------

func MinimumDX(t *Target) int {
	dx := 1
	for dist := t.xmin; dist > 0; dist -= dx {
		dx++
	}
	return dx
}

//------------------------------------------------------------------------------

type Projectile struct {
	x, y int
	dx, dy int
}

func MakeProjectile(dx, dy int) Projectile {
	return Projectile{x: 0, y: 0, dx: dx, dy: dy}
}

func (this *Projectile) Update() {
	this.x += this.dx
	this.y += this.dy
	this.dx -= sign(this.dx)
	this.dy -= 1
}

type Classification int

const (
	OnTarget Classification = 0
	Above = 1
	Below = 2
	Left = 4
	Right = 8
)

type Result int

const (
	Hit Result = iota
	Miss
)

func (this *Projectile) FireAt(t *Target) (bool, int) {
	c := this.Classify(t)
	h := 0

	for {
		//fmt.Printf("%15v %04b\n", this, this.Classify(t))
		if c == 0 {
			return true, h
		} else if (c & (Below|Right)) != 0 {
			return false, h
		}
		this.Update()
		c = this.Classify(t)
		if this.y > h {
			h = this.y
		}
	}
}

func (this *Projectile) Classify(t *Target) Classification {
	classification := OnTarget

	if this.x > t.xmax {
		classification |= Right
	} else if this.x < t.xmin {
		classification |= Left
	}

	if this.y > t.ymax {
		classification |= Above
	} else if this.y < t.ymin {
		classification |= Below
	}

	return classification
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

//------------------------------------------------------------------------------

var LineRegex = regexp.MustCompile("x=(-?\\d+)\\.\\.(-?\\d+), y=(-?\\d+)\\.\\.(-?\\d+)$")

type Target struct {
	xmin, xmax, ymin, ymax int
}

func ParseTarget(line string) Target {
	matches := LineRegex.FindStringSubmatch(line)
	p := aoc.ParseInts(matches[1:])

	return Target{xmin: p[0], xmax: p[1], ymin: p[2], ymax: p[3]}
}

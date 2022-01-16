package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

type Item uint8

const (
	Empty    Item = '.'
	East          = '>'
	South         = 'v'
	Leaving       = 'o'
	Entering      = '*'
)

type Seafloor struct {
	w, h int
	pos  []Item
}

func NewSeafloor(filename string) Seafloor {
	lines := aoc.GetInputLines(filename)

	w := len(lines[0])
	h := len(lines)
	pos := make([]Item, w*h)

	i := 0
	for _, line := range lines {
		for _, data := range line {
			pos[i] = Item(data)
			i++
		}
	}

	return Seafloor{w, h, pos}
}

func (this Seafloor) String() string {
	var b strings.Builder

	i := 0
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			b.WriteRune(rune(this.pos[i]))
			i++
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (this *Seafloor) Step() int {

	moves := 0

	// Eastbound
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			p := &this.pos[this.index(x, y)]
			q := &this.pos[this.index(x+1, y)]

			if *p == East && *q == Empty {
				*p, *q = Leaving, Entering
				moves++
			}
		}
		for x := 0; x < this.w; x++ {
			p := &this.pos[this.index(x, y)]
			if *p == Leaving {
				*p = Empty
			} else if *p == Entering {
				*p = East
			}
		}
	}

	//fmt.Println(this)

	// Southbound
	for x := 0; x < this.w; x++ {
		for y := 0; y < this.h; y++ {
			p := &this.pos[this.index(x, y)]
			q := &this.pos[this.index(x, y+1)]

			if *p == South && *q == Empty {
				*p, *q = Leaving, Entering
				moves++
			}
		}
		for y := 0; y < this.h; y++ {
			p := &this.pos[this.index(x, y)]
			if *p == Leaving {
				*p = Empty
			} else if *p == Entering {
				*p = South
			}
		}
	}

	return moves
}

func (this *Seafloor) index(x, y int) int {
	x %= this.w
	y %= this.h
	return y*this.w + x
}

func main() {
	fmt.Println(part1(aoc.GetFilename()))
}

func part1(filename string) int {
	seafloor := NewSeafloor(filename)

	steps := 0
	for {
		moves := seafloor.Step()
		steps++
		if moves == 0 {
			return steps
		}
	}
}

package main

import (
	"advent-of-code/common"
	"fmt"
	"strings"
)

func main() {
	filename := common.GetFilename()
	lines := common.GetInputLines(filename)

	p1 := newPart1()
	p2 := newPart2()

	for _, line := range lines {
		words := strings.Split(line, " ")
		count := common.ParseInt(words[1])

		p1.update(words[0], count)
		p2.update(words[0], count)
	}

	fmt.Println(p1.value())
	fmt.Println(p2.value())
}

type part1 struct {
	depth int
	pos   int
}

func newPart1() part1 {
	return part1{depth: 0, pos: 0}
}

func (p *part1) update(command string, arg int) {
	switch command {
	case "forward":
		p.pos += arg
	case "down":
		p.depth += arg
	case "up":
		p.depth -= arg
	}
}

func (p *part1) value() int {
	return p.pos * p.depth
}

type part2 struct {
	depth int
	pos   int
	aim   int
}

func newPart2() part2 {
	return part2{depth: 0, pos: 0, aim: 0}
}

func (p *part2) update(command string, arg int) {
	switch command {
	case "forward":
		p.pos += arg
		p.depth += p.aim * arg
	case "down":
		p.aim += arg
	case "up":
		p.aim -= arg
	}
}

func (p *part2) value() int {
	return p.pos * p.depth
}

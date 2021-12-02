package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	p1 := newPart1()
	p2 := newPart2()

	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		count, err := strconv.Atoi(words[1])
		if err != nil {
			log.Fatal(err)
		}

		p1.update(words[0], count)
		p2.update(words[0], count)
	}

	fmt.Println(p1.value())
	fmt.Println(p2.value())
}

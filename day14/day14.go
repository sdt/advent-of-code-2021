package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

type Rules map[string]string

func main() {
	start, rules := getInput(aoc.GetFilename())

	fmt.Println(part1(start, rules))
}

func part1(start string, rules Rules) int {
	polymer := start
	//fmt.Println(polymer)
	for i := 0; i < 10; i++ {
		polymer = rules.Apply(polymer)
		//fmt.Println(polymer)
	}

	return score(polymer)
}

func score(polymer string) int {
	hist := make([]int, 26)
	for _, letter := range polymer {
		hist[int(letter - 'A')]++
	}
	least, most := len(polymer), 0

	for _, count := range hist {
		if count > most {
			most = count
		} else if count > 0 && count < least {
			least = count
		}
	}
	return most - least
}

func getInput(filename string) (string, Rules) {
	lines := aoc.GetInputLines(filename)
	start := lines[0]
	rules := make(Rules)

	for _, line := range lines[2:] {
		parts := strings.Split(line, " -> ")
		from := parts[0]
		to := parts[1]
		rules[from] = to + from[1:]
	}

	return start, rules
}

func (r Rules) Apply(in string) string {
	n := len(in) - 1
	out := in[0:1]
	for i := 0; i < n; i++ {
		from := in[i:i+2]
		out += r[from]
	}
	return out
}

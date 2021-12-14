package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
	"strings"
)

type Rules map[string]byte
type Hist map[byte]int
type Signature struct {
	lhs, rhs byte
	depth    int
}
type Memo map[Signature]Hist

var memo Memo = make(Memo)

func main() {
	start, rules := getInput(aoc.GetFilename())

	fmt.Println(part1(start, rules))
	fmt.Println(part2(start, rules))
}

func part1(start string, rules Rules) int {
	return findScore(start, rules, 10)
}

func part2(start string, rules Rules) int {
	return findScore(start, rules, 40)
}

func findScore(start string, rules Rules, depth int) int {
	hist := make(Hist)
	for _, letter := range start {
		hist[byte(letter)]++
	}

	n := len(start) - 1
	for i := 0; i < n; i++ {
		addHist(hist, recurse(start[i], start[i+1], rules, depth))
	}

	least, most := math.MaxInt, 0
	for _, count := range hist {
		if count > most {
			most = count
		} else if count > 0 && count < least {
			least = count
		}
	}

	return most - least
}

func recurse(lhs, rhs byte, rules Rules, depth int) Hist {
	signature := Signature{lhs: lhs, rhs: rhs, depth: depth}
	if answer, found := memo[signature]; found {
		return answer
	}

	hist := make(Hist)
	mid := rules[string([]byte{lhs, rhs})]
	hist[mid] = 1

	if depth--; depth > 0 {
		addHist(hist, recurse(lhs, mid, rules, depth))
		addHist(hist, recurse(mid, rhs, rules, depth))
	}

	memo[signature] = hist
	return hist
}

func addHist(to, from Hist) {
	for key, value := range from {
		to[key] += value
	}
}

func getInput(filename string) (string, Rules) {
	lines := aoc.GetInputLines(filename)
	start := lines[0]
	rules := make(Rules)

	for _, line := range lines[2:] {
		parts := strings.Split(line, " -> ")
		from := parts[0]
		to := parts[1][0]
		rules[from] = to
	}

	return start, rules
}

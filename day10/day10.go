package main

import (
	"advent-of-code/aoc"
	"fmt"
	"sort"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		stack := make([]rune, 0)
		count := 0
		for _, symbol := range line {
			if isOpen(symbol) {
				stack = append(stack, symbol)
				count++
				continue
			}

			if count == 0 {
				break // incomplete
			}

			match := stack[count-1]
			if !isMatch(match, symbol) {
				total += syntaxScore(symbol)
				break // corrupted
			}
			stack = stack[:count-1]
			count--
		}
	}
	return total
}

func part2(lines []string) int {
	scores := make([]int, 0)
	for _, line := range lines {
		stack := make([]rune, 0)
		count := 0
		corrupted := false
		for _, symbol := range line {
			if isOpen(symbol) {
				stack = append(stack, symbol)
				count++
				continue
			}

			if count == 0 {
				corrupted = true
				break // incomplete
			}

			match := stack[count-1]
			if !isMatch(match, symbol) {
				corrupted = true
				break // corrupted
			}
			stack = stack[:count-1]
			count--
		}

		if corrupted {
			continue
		}

		score := 0
		for i := count - 1; i >= 0; i-- {
			score = score*5 + autocompleteScore(stack[i])
		}
		scores = append(scores, score)
	}
	sort.Sort(sort.IntSlice(scores))
	return scores[len(scores)/2]
}

func isOpen(symbol rune) bool {
	switch symbol {
	case '{', '<', '[', '(':
		return true

	default:
		return false
	}
}

func isMatch(a, b rune) bool {
	return (a == '{' && b == '}') ||
		(a == '<' && b == '>') ||
		(a == '[' && b == ']') ||
		(a == '(' && b == ')')
}

func syntaxScore(symbol rune) int {
	switch symbol {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	default:
		return 0
	}
}

func autocompleteScore(symbol rune) int {
	switch symbol {
	case '(':
		return 1
	case '[':
		return 2
	case '{':
		return 3
	case '<':
		return 4
	default:
		return 0
	}
}

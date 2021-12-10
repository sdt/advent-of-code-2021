package main

import (
	"advent-of-code/common"
	"fmt"
)

func main() {
	filename := common.GetFilename()
	lines := common.GetInputLines(filename)

	fmt.Println(part1(lines))
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

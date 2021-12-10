package main

import (
	"advent-of-code/common"
	"fmt"
	"log"
	"strings"
)

const SegmentCount = 7

var DigitSegments = [...]string{
	"ABCEFG",	// 0
	"CF",		// 1
	"ACDEG",	// 2
	"ACDFG",	// 3
	"BCDF",		// 4
	"ABDFG",	// 5
	"ABDEFG",	// 6
	"ACF",		// 7
	"ABCDEFG",	// 8
	"ABCDFG",	// 9
}

func main() {
	filename := common.GetFilename()
	lines := common.GetInputLines(filename)

	makeSignatureMap()

	fmt.Println(part1(lines))
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		_, digits := parseLine(line)
		for _, digit := range digits {
			length := len(digit)
			if length == 2 || length == 3 || length == 4 || length == 7 {
				total++
			}
		}
	}

	return total
}

func makeSignatureMap() map[string]int {
	var counts [SegmentCount][SegmentCount+1]int
	for _, segments := range DigitSegments {
		length := len(segments)
		for _, segment := range segments {
			counts[segment -'A'][length]++
		}
	}

	signatureMap := make(map[string]int)
	for i, count := range counts {
		signature := makeSignature(count[0:])
		signatureMap[signature] = i
	}

	return signatureMap
}

func makeSignature(counts []int) string {
	var out strings.Builder
	for _, count := range counts {
		out.WriteRune(rune('0' + count))
	}
	return out.String()
}

func parseLine(line string) ([]string, []string) {
	words := strings.Split(line, " ")
	if len(words) != 15 {
		log.Fatal("Expecting 15 words in input, got ", len(words))
	}

	signals := words[0:10]	// ten signal patterns
	digits := words[11:15]	// four output digits

	return signals, digits
}

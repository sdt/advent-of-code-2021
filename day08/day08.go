package main

import (
	"advent-of-code/aoc"
	"fmt"
	"log"
	"sort"
	"strings"
)

const SegmentCount = 7

var DigitSegments = [...]string{
	"ABCEFG",  // 0
	"CF",      // 1
	"ACDEG",   // 2
	"ACDFG",   // 3
	"BCDF",    // 4
	"ABDFG",   // 5
	"ABDEFG",  // 6
	"ACF",     // 7
	"ABCDEFG", // 8
	"ABCDFG",  // 9
}

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
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

func part2(lines []string) int {
	signatureMap := makeSignatureMap(DigitSegments[0:])
	signalMap := make(map[rune]rune)

	digitMap := make(map[string]int)
	for digit, segments := range DigitSegments {
		digitMap[segments] = digit
	}

	total := 0
	for _, line := range lines {
		signals, digits := parseLine(line)

		signatures := makeSignatures(signals, 'a')
		for i, signature := range signatures {
			signalMap[rune('a'+i)] = signatureMap[signature]
		}

		value := 0
		for _, digit := range digits {
			segments := mapDigits(digit, signalMap)
			value = value*10 + digitMap[segments]
		}
		total += value
	}

	return total
}

func mapDigits(digits string, signalMap map[rune]rune) string {
	segments := make([]int, len(digits))

	for i, digit := range digits {
		segments[i] = int(signalMap[digit])
	}

	sort.Sort(sort.IntSlice(segments))
	var digit strings.Builder
	for _, segment := range segments {
		digit.WriteRune(rune(segment))
	}

	return digit.String()
}

func makeSignatureMap(signals []string) map[string]rune {
	signatureMap := make(map[string]rune)
	for i, signature := range makeSignatures(signals, 'A') {
		signatureMap[signature] = rune('A' + i)
	}
	return signatureMap
}

func makeSignatures(signals []string, base rune) []string {
	var counts [SegmentCount][SegmentCount + 1]int
	for _, segments := range signals {
		length := len(segments)
		for _, segment := range segments {
			counts[segment-base][length]++
		}
	}

	signatures := make([]string, SegmentCount)
	for i, count := range counts {
		signatures[i] = makeSignature(count[0:])
	}

	return signatures

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

	signals := words[0:10] // ten signal patterns
	digits := words[11:15] // four output digits

	return signals, digits
}

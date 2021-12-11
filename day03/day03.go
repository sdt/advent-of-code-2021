package main

import (
	"advent-of-code/aoc"
	"fmt"
	"log"
	"strconv"
)

type MostCommonBits int

const (
	Same MostCommonBits = iota
	MoreZeros
	MoreOnes
)

func main() {
	filename := aoc.GetFilename()
	reports, width := getReports(filename)

	fmt.Println(part1(reports, width))
	fmt.Println(part2(reports, width))
}

func part1(reports []int, width int) int {
	gamma := 0
	epsilon := 0
	for i := 0; i < width; i++ {
		if getMostCommonBits(reports, i) == MoreZeros {
			epsilon |= (1 << i)
		} else {
			gamma |= (1 << i)
		}
	}

	return gamma * epsilon
}

func part2(reports []int, width int) int {
	oxygens := getRatings(reports, width, 1)
	co2s := getRatings(reports, width, 0)

	return oxygens * co2s
}

func getRatings(reports []int, width, match int) int {
	ratings := reports
	for bit := 0; bit < width; bit++ {
		ratings = filter(ratings, width-bit-1, match)
		if len(ratings) == 1 {
			return ratings[0]
		}
		if len(ratings) == 0 {
			log.Fatal("No more matching reports")
		}
	}
	return 0
}

func filter(in []int, bit, defaultMatch int) []int {
	var match int
	switch getMostCommonBits(in, bit) {
	case MoreOnes:
		match = defaultMatch
	case MoreZeros:
		match = 1 ^ defaultMatch
	default:
		match = defaultMatch
	}
	match <<= bit
	mask := 1 << bit

	out := make([]int, 0)
	for _, value := range in {
		if value&mask == match {
			out = append(out, value)
		}
	}
	return out
}

func getMostCommonBits(reports []int, bit int) MostCommonBits {
	sense := 0
	mask := 1 << bit

	for _, report := range reports {
		if (report & mask) != 0 {
			sense++
		} else {
			sense--
		}
	}

	if sense > 0 {
		return MoreOnes
	} else if sense < 0 {
		return MoreZeros
	} else {
		return Same
	}
}

func getReports(filename string) ([]int, int) {
	lines := aoc.GetInputLines(filename)
	reports := make([]int, 0)

	for _, line := range lines {
		report, err := strconv.ParseInt(line, 2, 32)
		aoc.CheckErr(err)
		reports = append(reports, int(report))
	}

	return reports, len(lines[0])
}

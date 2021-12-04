package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	reports, width := getReports(os.Args[1])
	gamma := getMostCommonBits(reports, width)

	fmt.Println(part1(gamma, width))
	//fmt.Println(getIncreases(depths, 3))
}

func part1(gamma int, width int) int {
	mask := (1 << width) - 1
	epsilon := ^gamma & mask
	return gamma * epsilon
}

func getMostCommonBits(reports []int, width int) int {
	gamma := 0
	for bit := 0; bit < width; bit++ {
		mask := 1 << bit
		sense := 0
		for _, report := range reports {
			if report & mask != 0 {
				sense++
			} else {
				sense--
			}
		}

		if sense >= 0 {
			gamma |= mask
		}
	}

	return gamma
}

func getReports(filename string) ([]int, int) {
	lines := getInput(filename)
	reports := make([]int, 0)

	for _, line := range lines {
		report, err := strconv.ParseInt(line, 2, 32)
		if err != nil {
			log.Fatal(err)
		}

		reports = append(reports, int(report))
	}

	return reports, len(lines[0])
}

func getInput(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

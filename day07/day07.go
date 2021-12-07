package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type FuelFunction func(int) int

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	positions := getInput(os.Args[1])

	fmt.Println(part1(positions))
	fmt.Println(part2(positions))
}

func part1(positions []int) int {
	return getLeastFuel(positions, func(dist int) int {
		return abs(dist)
	})
}

func part2(positions []int) int {
	return getLeastFuel(positions, func(dist int) int {
		dist = abs(dist)
		return (dist*dist + dist) / 2
	})
}

func getLeastFuel(positions []int, getFuel FuelFunction) int {
	min := math.MaxInt
	max := math.MinInt

	for _, position := range positions {
		if position < min {
			min = position
		}
		if position > max {
			max = position
		}
	}

	leastTotal := math.MaxInt

	for dest := min; dest <= max; dest++ {
		total := 0
		for _, position := range positions {
			total += getFuel(position - dest)
		}
		if total < leastTotal {
			leastTotal = total
		}
	}
	return leastTotal
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getInput(filename string) []int {
	words := strings.Split(getInputLines(filename)[0], ",")

	positions := make([]int, len(words))
	for i, word := range words {
		value, err := strconv.Atoi(word)
		if err != nil {
			log.Fatal(err)
		}
		positions[i] = value
	}
	return positions
}

func getInputLines(filename string) []string {
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

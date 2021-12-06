package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var memo map[int]int = make(map[int]int)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	fmt.Println(part1(os.Args[1]))
	fmt.Println(part2(os.Args[1]))
}

func part1(filename string) int {
	startCycles := getInput(filename)
	total := 0
	for _, start := range startCycles {
		total += totalFish(80 - start - 1)
	}
	return total
}

func part2(filename string) int {
	startCycles := getInput(filename)
	total := 0
	for _, start := range startCycles {
		total += totalFish(256 - start - 1)
	}
	return total
}

func totalFish(days int) int {
	if count, found := memo[days]; found {
		return count
	}
	count := 1

	for n := days; n >= 0; n -= 7 {
		count += totalFish(n - 8 - 1)
	}

	memo[days] = count
	return count
}

func getInput(filename string) []int {
	words := strings.Split(getInputLines(filename)[0], ",")

	startCycles := make([]int, len(words))
	for i, word := range words {
		value, err := strconv.Atoi(word)
		if err != nil {
			log.Fatal(err)
		}
		startCycles[i] = value
	}
	return startCycles
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

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
	part1(os.Args[1])
	part2(os.Args[1])
}

func part1(filename string) {
	fmt.Println(getIncreases(filename, 1))
}

func part2(filename string) {
	fmt.Println(getIncreases(filename, 3))
}

func getIncreases(filename string, windowSize int) int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	window := make([]int, windowSize)
	increases := 0

	for scanner.Scan() {
		curr, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		count++
		if count > windowSize {
			prev := window[count % windowSize]
			if curr > prev {
				increases++
			}
		}
		window[count % windowSize] = curr
	}

	return increases
}

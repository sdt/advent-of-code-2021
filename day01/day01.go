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
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	started := false
	increases := 0
	var prev int

	for scanner.Scan() {
		curr, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if started {
			if curr > prev {
				increases++
			}
		} else {
			started = true
		}
		prev = curr
	}

	fmt.Println(increases)
}

func part2(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const windowSize = 3
	count := 0
	var window [windowSize]int
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

	fmt.Println(increases)
}

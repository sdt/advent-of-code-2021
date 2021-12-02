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

	depths, err := getDepths(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(getIncreases(depths, 1))
	fmt.Println(getIncreases(depths, 3))
}

func getDepths(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	depths := make([]int, 0)
	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		depths = append(depths, depth)
	}
	return depths, nil
}

func getIncreases(depths []int, windowSize int) int {
	increases := 0
	for i, after := range depths[windowSize:] {
		before := depths[i]
		if after > before {
			increases++
		}
	}
	return increases
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	reports, err := getReports(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(part1(reports))
	//fmt.Println(getIncreases(depths, 3))
}

func part1(reports []string) int {
	bits := len(reports[0])
	senses := make([]int, bits)
	for _, report := range reports {
		for i, bit := range report {
			if bit == '0' {
				senses[i]--
			} else {
				senses[i]++
			}
		}
	}

	gamma := 0
	epsilon := 0
	for i, sense := range senses {
		value := 1 << (bits - i - 1)
		if sense > 0 {
			gamma |= value
		} else {
			epsilon |= value
		}
	}

	return gamma * epsilon
}

func getReports(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reports := make([]string, 0)
	for scanner.Scan() {
		reports = append(reports, scanner.Text())
	}
	return reports, nil
}

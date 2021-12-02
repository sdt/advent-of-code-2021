package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	fmt.Println(part1(os.Args[1]))
}

func part1(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	depth := 0
	pos := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		count, err := strconv.Atoi(words[1])
		if err != nil {
			log.Fatal(err)
		}

		switch words[0] {
		case "forward":
			pos += count
		case "down":
			depth += count
		case "up":
			depth -= count
		}
	}
	return depth * pos
}

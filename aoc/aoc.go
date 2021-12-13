package aoc

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func GetFilename() string {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	return os.Args[1]
}

func GetInputLines(filename string) []string {
	file, err := os.Open(filename)
	CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func ParseInts(in []string) []int {
	out := make([]int, len(in))
	for i, line := range in {
		out[i] = ParseInt(line)
	}
	return out
}

func ParseInt(in string) int {
	out, err := strconv.Atoi(in)
	CheckErr(err)
	return out
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const Size = 5
var WhiteSpace = regexp.MustCompile(" +")
const Mark = 0x80000000

type Board struct {
	cells [Size][Size]int
	score int
}
type Draw []int
type Bingo struct {
	draw   Draw
	boards []*Board
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: ", os.Args[0], " input-file")
	}

	bingo := parseBingo(getInputLines(os.Args[1]))

	fmt.Println(part1(bingo))
}

func part1(bingo Bingo) int {
	for _, drawn := range bingo.draw {
		//fmt.Printf("Playing %d\n", drawn)
		for _, board := range bingo.boards {
			if board.mark(drawn) {
				//fmt.Printf("Board %d has a %d: %v\n", i, drawn, board)
				if board.isComplete() {
					return drawn * board.incompleteCellSum()
				}
			}
		}
	}
	return 0
}

func parseBingo(lines []string) Bingo {
	numBoards := (len(lines) - 1) / (Size + 1)
	bingo := Bingo{draw: parseDraw(lines[0]), boards: make([]*Board, numBoards)}

	for i := 0; i < numBoards; i++ {
		firstLine := i * (Size+1) + 2
		lastLine := firstLine + Size
		bingo.boards[i] = parseBoard(lines[firstLine:lastLine])
	}

	return bingo
}

func parseBoard(lines []string) *Board {
	board := Board{}
	for row, line := range lines {
		line = strings.TrimSpace(line)
		for col, word := range WhiteSpace.Split(line, -1) {
			value, err := strconv.Atoi(word)
			if err != nil {
				log.Fatal(err)
			}
			board.cells[row][col] = value
		}
	}
	return &board
}

func (board *Board) incompleteCellSum() int {
	sum := 0
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			value := board.cells[row][col]
			if value & Mark == 0 {
				sum += value
			}
		}
	}
	return sum
}

func (board *Board) isComplete() bool {
	for i := 0; i < Size; i++ {
		if board.isRowComplete(i) || board.isColComplete(i) {
			return true
		}
	}
	return false
}

func (board *Board) isRowComplete(row int) bool {
	for col := 0; col < Size; col++ {
		value := board.cells[row][col]
		if value & Mark == 0 {
			return false
		}
	}
	return true
}

func (board *Board) isColComplete(col int) bool {
	for row := 0; row < Size; row++ {
		value := board.cells[row][col]
		if value & Mark == 0 {
			return false
		}
	}
	return true
}

func (board *Board) mark(number int) bool {
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			if board.cells[row][col] == number {
				board.cells[row][col] |= Mark
				return true
			}
		}
	}
	return false
}

func parseDraw(line string) Draw {
	draw := make([]int, 0)
	for _, number := range strings.Split(line, ",") {
		value, err := strconv.Atoi(number)
		if err != nil {
			log.Fatal(err)
		}
		draw = append(draw, value)
	}
	return draw
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

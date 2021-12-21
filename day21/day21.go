package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

func main() {
	startPositions := parseInput(aoc.GetFilename())

	fmt.Println(part1(startPositions))
}

//------------------------------------------------------------------------------

func part1(startPositions [2]int) int {
	player := [2]Player{
		MakePlayer(startPositions[0]),
		MakePlayer(startPositions[1]),
	}

	dice := MakeDice(100)

	for next := 0; ; next ^= 1 {
		player[next].TakeTurn(&dice)
		if player[next].score >= 1000 {
			loser := next ^ 1
			return player[loser].score * dice.rolls
		}
	}
}

//------------------------------------------------------------------------------

type Player struct {
	position int
	score int
}

func MakePlayer(position int) Player {
	return Player{position: position, score: 0}
}

func (this *Player) TakeTurn(dice *Dice) {
	roll := 0
	for i := 0; i < 3; i++ {
		roll += dice.Roll()
	}

	this.position = ((this.position + roll - 1) % 10) + 1
	this.score += this.position
}

//------------------------------------------------------------------------------

type Dice struct {
	sides int
	current int
	rolls int
}

func MakeDice(sides int) Dice {
	return Dice{sides: sides, current: 1, rolls: 0}
}

func (this *Dice) Roll() int {
	ret := this.current
	if this.current++; this.current > this.sides {
		this.current = 1
	}
	this.rolls++
	return ret
}

//------------------------------------------------------------------------------

func parseInput(filename string) [2]int {
	lines := aoc.GetInputLines(filename)

	var numbers [2]int
	for i := range numbers {
		parts := strings.Split(lines[i], ": ")
		numbers[i] = aoc.ParseInt(parts[1])
	}

	return numbers
}

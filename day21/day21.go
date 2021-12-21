package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

func main() {
	startPositions := parseInput(aoc.GetFilename())

	fmt.Println(part1(startPositions))
	fmt.Println(part2(startPositions))
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

func part2(startPositions [2]int) uint64 {
	game := Game2{player: [...]Player2{
		MakePlayer2(uint8(startPositions[0])),
		MakePlayer2(uint8(startPositions[1])),
	}}

	rolls := MakeRolls()

	multiverse := make(Multiverse)
	multiverse[game] = 1

	scores := Scores{}

	for player := 0; len(multiverse) > 0; player ^= 1 {
		multiverse = TakeTurn(player, multiverse, rolls, &scores)
	}

	if scores.wins[0] > scores.wins[1] {
		return scores.wins[0]
	} else {
		return scores.wins[1]
	}
}

type Roll struct {
	score, count uint8
}

type Multiverse map[Game2]uint64

type Wins [2]uint64

type Game2 struct {
	player [2]Player2
}

type Player2 struct {
	position uint8
	score uint8
}

type Scores struct {
	wins [2]uint64
}

func MakePlayer2(position uint8) Player2 {
	return Player2{position: position, score: 0}
}


func MakeRolls() []Roll {
	var hist [10]uint8

	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				hist[ i + j + k ]++
			}
		}
	}

	rolls := make([]Roll, 0)
	for score, count := range hist {
		if count > 0 {
			rolls = append(rolls, Roll{ uint8(score), count })
		}
	}
	return rolls
}

func TakeTurn(whichPlayer int, from Multiverse, rolls []Roll, scores *Scores) Multiverse {

	to := make(Multiverse)

	for game, gameCount := range from {
		for _, roll := range rolls {
			player := game.player[whichPlayer].Move(roll.score)

			totalCount := gameCount * uint64(roll.count)

			if player.score >= 21 {
				scores.wins[whichPlayer] += totalCount
			} else {
				nextGame := game.NextGame(whichPlayer, player)
				if count, found := to[nextGame]; found {
					totalCount += count
				}
				to[nextGame] = totalCount
			}
		}
	}

	return to
}

func (this* Game2) NextGame(whichPlayer int, player Player2) Game2 {
	if whichPlayer == 0 {
		return Game2{ player: [...]Player2{ player, this.player[1] } }
	} else {
		return Game2{ player: [...]Player2{ this.player[0], player } }
	}
}

func (this* Player2) Move(roll uint8) Player2 {
	position := ((this.position + roll - 1) % 10) + 1
	return Player2{ position: position, score: this.score + position }
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

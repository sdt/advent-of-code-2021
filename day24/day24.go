package main

import (
	"fmt"
)

type Word int64
type Params struct {
	a, b, c Word
}

var params = [...]Params{
	Params{ 1, 12,  4},
	Params{ 1, 11, 10},
	Params{ 1, 14, 12},
	Params{26, -6, 14},
	Params{ 1, 15,  6},
	Params{ 1, 12, 16},
	Params{26, -9,  1},
	Params{ 1, 14,  7},
	Params{ 1, 14,  8},
	Params{26, -5, 11},
	Params{26, -9,  8},
	Params{26, -5,  3},
	Params{26, -2,  1},
	Params{26, -7,  8},
}

func main() {
	// 91398299697996
	// 41171183141291
	fmt.Println(part1())
	fmt.Println(part2())
}

func part1() string {
	return solve(func (i int) Word {
		return Word(9 - i)
	})
}

func part2() string {
	return solve(func (i int) Word {
		return Word(i + 1)
	})
}


func solve(f func(int)Word) string {
	solution := make([]byte, len(params))

	inputs := make([]Word, 9)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = f(i)
	}

	recurse(solution, inputs, 0, 0)

	return string(solution)
}

func recurse(solution[] byte, inputs[] Word, depth int, z Word) bool {
	if depth == len(params) {
		return true
	}

	p := &params[depth]

	for _, input := range inputs {
		nextZ, ok := NativeChunk(input, p, z)
		if (ok || p.a == 1) && recurse(solution, inputs, depth+1, nextZ) {
			solution[depth] = '0' + byte(input)
			return true
		}
	}
	return false
}


// a: line  5
// b: line  6
// c: line 16

func NativeChunk(input Word, p *Params, z Word) (Word, bool) {
	if z%26 != input-p.b {
		return (z/p.a)*26 + input + p.c, false
	}
	return z / p.a, true
}

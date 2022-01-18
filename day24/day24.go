package main

import (
	"fmt"
	"strconv"
)

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
	lowest, highest := solve()
	fmt.Println(highest)
	fmt.Println(lowest)
}

func solve() (string, string) {
	partials := []Partial{Partial{}}
	from := 0
	for {
		size := chunkSize(from)
		if size == 0 {
			break
		}

		next := []Partial{}
		for _, partial := range partials {
			next = append(next, TestPrefixes(partial, size)...)
		}

		//fmt.Printf("From=%d size=%d remaining=%d\n", from, size, len(next))
		from += size
		partials = next
	}
	return partials[0].input, partials[len(partials)-1].input
}

func chunkSize(from int) int {
	if from == len(params) {
		return 0
	}

	i := 0
	for params[from+i].a == 1 {
		i++
	}
	return i + 1
}

type Word int64

// a: line  5
// b: line  6
// c: line 16

func pow(base, exponent int) int {
	ret := 1
	for i := 0; i < exponent; i++ {
		ret *= base
	}
	return ret
}

type Partial struct {
	input string
	z     Word
}

func TestPrefixes(partial Partial, more int) []Partial {
	count := pow(9, more)
	out := make([]Partial, 0)

	offset := len(partial.input)

	for i := 0; i < count; i++ {
		input := convert(Word(i), more)

		keep := false
		z := partial.z
		for j := 0; j < more; j++ {
			k := offset + j
			z, keep = NativeChunk(input[j], &params[k], z)
		}
		if keep {
			out = append(out, Partial{input: partial.input + input, z: z})
			//fmt.Println(input, z)
		}
	}
	return out
}

func NativeChunk(input byte, p *Params, z Word) (Word, bool) {
	w := Word(input - '0')
	if z%26 != w-p.b {
		return (z/p.a)*26 + w + p.c, false
	}
	return z / p.a, true
}

func convert(value Word, width int) string {
	format := fmt.Sprintf("%%0%ds", width)
	digits := fmt.Sprintf(format, strconv.FormatInt(int64(value), 9))
	number := make([]byte, len(digits))
	for i, digit := range digits {
		number[i] = byte(digit) + 1
	}
	return string(number)
}

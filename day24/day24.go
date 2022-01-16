package main

import (
	"fmt"
	"strconv"
)

var a = [...]Word{1, 1, 1, 26, 1, 1, 26, 1, 1, 26, 26, 26, 26, 26}
var b = [...]Word{12, 11, 14, -6, 15, 12, -9, 14, 14, -5, -9, -5, -2, -7}
var c = [...]Word{4, 10, 12, 14, 6, 16, 1, 7, 8, 11, 8, 3, 1, 8}

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
	if from == len(a) {
		return 0
	}

	i := 0
	for a[from+i] == 1 {
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
		input := partial.input + convert(Word(i), more)

		keep := false
		z := partial.z
		for j := 0; j < more; j++ {
			k := offset + j
			z, keep = NativeChunk(input[k], a[k], b[k], c[k], z)
		}
		if keep {
			out = append(out, Partial{input: input, z: z})
			//fmt.Println(input, z)
		}
	}
	return out
}

func NativeChunk(input byte, a, b, c, z Word) (Word, bool) {
	w := Word(input - '0')
	if z%26 != w-b {
		return (z/a)*26 + w + c, false
	}
	return z / a, true
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

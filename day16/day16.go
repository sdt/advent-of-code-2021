package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	bitstream := getInput(filename)

	fmt.Println(part1(&bitstream))
	fmt.Println(part2(&bitstream))
}

func part1(bitstream *Bitstream) int {
	bitstream.Reset()

	part1 := uint64(0)
	depth := 0
	parseBitstream(bitstream, &ParseConfig{
		onLiteral: func(version, value uint64) {
			part1 += version
			//fmt.Printf("Literal: version=%d value=%d depth=%d\n", version, value, depth)
		},
		onBeginOperator: func(version, typeId uint64) {
			part1 += version
			//fmt.Printf("Operator: version=%d type=%d depth=%d->%d\n", version, typeId, depth, depth+1)
			depth++
		},
		onEndOperator: func() {
			//fmt.Printf("End operator: depth=%d->%d\n", depth, depth-1)
			depth--
		}})
	return int(part1)
}

type Opcode int64

const (
	Sum Opcode = iota
	Product
	Minimum
	Maximum
	Literal
	GreaterThan
	LessThan
	EqualTo
)

func part2(bitstream *Bitstream) int {
	bitstream.Reset()

	stack := MakeOperatorStack(Literal)

	parseBitstream(bitstream, &ParseConfig{
		onLiteral: func(_, value uint64) {
			stack.Top().AddValue(value)
		},
		onBeginOperator: func(_, typeId uint64) {
			stack.Push(Opcode(typeId))
		},
		onEndOperator: func() {
			value := stack.Top().Evaluate()
			stack.Pop()
			stack.Top().AddValue(value)
		}})
	return int(stack.Top().Evaluate())
}

func getInput(filename string) Bitstream {
	lines := aoc.GetInputLines(filename)
	return MakeBitstream(lines[0])
}

//------------------------------------------------------------------------------

type OnLiteral func(version, value uint64)
type OnBeginOperator func(version, typeId uint64)
type OnEndOperator func()

type ParseConfig struct {
	onLiteral       OnLiteral
	onBeginOperator OnBeginOperator
	onEndOperator   OnEndOperator
}

func parseBitstream(bitstream *Bitstream, config *ParseConfig) {
	for !bitstream.IsEmpty() {
		parsePacket(bitstream, config)
	}
}

func parsePacket(bitstream *Bitstream, config *ParseConfig) {
	version := bitstream.ReadBits(3)
	typeId := bitstream.ReadBits(3)

	if typeId == uint64(Literal) {
		config.onLiteral(version, parseLiteral(bitstream))
	} else {
		lengthTypeId := bitstream.ReadBits(1)
		config.onBeginOperator(version, typeId)
		if lengthTypeId == 0 {
			parseOperator0(bitstream, config)
		} else {
			parseOperator1(bitstream, config)
		}
		config.onEndOperator()
	}
}

func parseLiteral(bitstream *Bitstream) uint64 {
	ret := uint64(0)

	for {
		chunk := bitstream.ReadBits(5)
		ret = (ret << 4) | (chunk & 0xf)
		if (chunk & 0x10) == 0 {
			return ret
		}
	}
}

func parseOperator0(bitstream *Bitstream, config *ParseConfig) {
	bits := bitstream.ReadBits(15)
	endPosition := bitstream.BitPosition() + bits

	for {
		parsePacket(bitstream, config)
		if bitstream.BitPosition() >= endPosition {
			return
		}
	}
}

func parseOperator1(bitstream *Bitstream, config *ParseConfig) {
	packets := int(bitstream.ReadBits(11))

	for i := 0; i < packets; i++ {
		parsePacket(bitstream, config)
	}
}

//------------------------------------------------------------------------------

type Bitstream struct {
	data       []uint8
	byteCursor int
	bitCursor  int
}

func (this *Bitstream) Reset() {
	this.byteCursor = 0
	this.bitCursor = 0
}

func (this *Bitstream) BitPosition() uint64 {
	return uint64(this.byteCursor*8 + this.bitCursor)
}

// Is there less than a full byte remaining
func (this *Bitstream) IsEmpty() bool {
	bytes := len(this.data)
	return (this.byteCursor >= bytes) || ((this.byteCursor == bytes-1) && (this.bitCursor > 0))
}

func (this *Bitstream) ByteAlign() {
	if this.bitCursor != 0 {
		this.bitCursor = 0
		this.byteCursor++
	}
}

func (this *Bitstream) ReadBit() uint8 {
	shift := 8 - this.bitCursor - 1
	ret := (this.data[this.byteCursor] & (1 << shift)) >> shift
	if this.bitCursor++; this.bitCursor >= 8 {
		this.bitCursor = 0
		this.byteCursor++
	}
	return ret
}

func (this *Bitstream) ReadBits(n int) uint64 {
	ret := uint64(0)
	for i := 0; i < n; i++ {
		ret = (ret << 1) | uint64(this.ReadBit())
	}
	return ret
}

func MakeBitstream(hexchars string) Bitstream {
	data := make([]uint8, len(hexchars)/2)

	var nybble uint8
	for i, hexchar := range hexchars {
		if i&1 == 0 {
			nybble = parseNybble(hexchar)
		} else {
			data[i>>1] = (nybble << 4) | parseNybble(hexchar)
		}
	}

	return Bitstream{data: data, byteCursor: 0, bitCursor: 0}
}

func parseNybble(hexchar rune) uint8 {
	if hexchar >= '0' && hexchar <= '9' {
		return uint8(hexchar - '0')
	} else {
		return 10 + uint8(hexchar-'A')
	}
}

//------------------------------------------------------------------------------

type Operator struct {
	opcode Opcode
	values []uint64
}

func MakeOperator(opcode Opcode) Operator {
	return Operator{opcode: opcode, values: make([]uint64, 0)}
}

func (this *Operator) AddValue(value uint64) {
	this.values = append(this.values, value)
}

func (this *Operator) Evaluate() uint64 {
	value := this.values[0]

	switch this.opcode {
	case Sum:
		for _, x := range this.values[1:] {
			value += x
		}
	case Product:
		for _, x := range this.values[1:] {
			value *= x
		}
	case Minimum:
		for _, x := range this.values[1:] {
			if x < value {
				value = x
			}
		}
	case Maximum:
		for _, x := range this.values[1:] {
			if x > value {
				value = x
			}
		}
	case Literal:
		value = this.values[len(this.values)-1]
	case GreaterThan:
		if value > this.values[1] {
			value = 1
		} else {
			value = 0
		}
	case LessThan:
		if value < this.values[1] {
			value = 1
		} else {
			value = 0
		}
	case EqualTo:
		if value == this.values[1] {
			value = 1
		} else {
			value = 0
		}
	}

	return value
}

//------------------------------------------------------------------------------

type OperatorStack struct {
	stack []Operator
}

func MakeOperatorStack(defaultOpcode Opcode) OperatorStack {
	stack := []Operator{MakeOperator(defaultOpcode)}
	return OperatorStack{stack: stack}
}

func (this *OperatorStack) Push(opcode Opcode) {
	this.stack = append(this.stack, MakeOperator(opcode))
}

func (this *OperatorStack) Top() *Operator {
	return &this.stack[len(this.stack)-1]
}

func (this *OperatorStack) Pop() {
	this.stack = this.stack[0 : len(this.stack)-1]
}

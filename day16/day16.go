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

func part2(bitstream *Bitstream) int {
	bitstream.Reset()

	stack := make([]Operator, 0)
	stack = append(stack, MakeIdentity())

	parseBitstream(bitstream, &ParseConfig{
		onLiteral: func(_, value uint64) {
			stack[len(stack)-1].addParam(value)
		},
		onBeginOperator: func(_, typeId uint64) {
			var operator Operator
			switch typeId {
			case 0:
				operator = MakeSum()
			case 1:
				operator = MakeProduct()
			case 2:
				operator = MakeMinimum()
			case 3:
				operator = MakeMaximum()
			case 4:
				operator = MakeIdentity()
			case 5:
				operator = MakeGreaterThan()
			case 6:
				operator = MakeLessThan()
			case 7:
				operator = MakeEqualTo()
			}
			stack = append(stack, operator)
		},
		onEndOperator: func() {
			value := stack[len(stack)-1].evaluate()
			stack = stack[0 : len(stack)-1]
			stack[len(stack)-1].addParam(value)
		}})
	return int(stack[len(stack)-1].evaluate())
}

func indent(depth int) {
	fmt.Printf("\n")
	for i := 0; i < depth; i++ {
		fmt.Printf("  ")
	}
}

func getInput(filename string) Bitstream {
	lines := aoc.GetInputLines(filename)
	return MakeBitstream(lines[0])
}

//------------------------------------------------------------------------------

const Literal = 4

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

	if typeId == Literal {
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
	//fmt.Printf("Type 0 - %d bits\n", bits)
	endPosition := bitstream.BitPosition() + bits

	for {
		parsePacket(bitstream, config)
		if bitstream.BitPosition() >= endPosition {
			return
		}
	}
	//fmt.Printf("End type 0 - %d bits\n", bits)
}

func parseOperator1(bitstream *Bitstream, config *ParseConfig) {
	packets := int(bitstream.ReadBits(11))
	//fmt.Printf("Type 1 - %d packets\n", packets)

	for i := 0; i < packets; i++ {
		parsePacket(bitstream, config)
	}

	//fmt.Printf("End type 1 - %d packets\n", packets)
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

type Operator interface {
	addParam(param uint64)
	evaluate() uint64
}

//------------------------------------------------------------------------------

type Identity struct {
	value uint64
}

func MakeIdentity() *Identity {
	return &Identity{value: 0}
}

func (this *Identity) addParam(param uint64) {
	this.value = param
}

func (this *Identity) evaluate() uint64 {
	return this.value
}

//------------------------------------------------------------------------------

type Sum struct {
	value uint64
}

func MakeSum() *Sum {
	return &Sum{value: 0}
}

func (this *Sum) addParam(param uint64) {
	this.value += param
}

func (this *Sum) evaluate() uint64 {
	return this.value
}

//------------------------------------------------------------------------------

type Product struct {
	value uint64
}

func MakeProduct() *Product {
	return &Product{value: 1}
}

func (this *Product) addParam(param uint64) {
	this.value *= param
}

func (this *Product) evaluate() uint64 {
	return this.value
}

//------------------------------------------------------------------------------

type Minimum struct {
	value uint64
}

func MakeMinimum() *Minimum {
	return &Minimum{value: 0xffffffffffffffff}
}

func (this *Minimum) addParam(param uint64) {
	if param < this.value {
		this.value = param
	}
}

func (this *Minimum) evaluate() uint64 {
	return this.value
}

//------------------------------------------------------------------------------

type Maximum struct {
	value uint64
}

func MakeMaximum() *Maximum {
	return &Maximum{value: 0}
}

func (this *Maximum) addParam(param uint64) {
	if param > this.value {
		this.value = param
	}
}

func (this *Maximum) evaluate() uint64 {
	return this.value
}

//------------------------------------------------------------------------------

type LessThan struct {
	value []uint64
}

func MakeLessThan() *LessThan {
	return &LessThan{value: make([]uint64, 0)}
}

func (this *LessThan) addParam(param uint64) {
	this.value = append(this.value, param)
}

func (this *LessThan) evaluate() uint64 {
	return fromBool(this.value[0] < this.value[1])
}

//------------------------------------------------------------------------------

type GreaterThan struct {
	value []uint64
}

func MakeGreaterThan() *GreaterThan {
	return &GreaterThan{value: make([]uint64, 0)}
}

func (this *GreaterThan) addParam(param uint64) {
	this.value = append(this.value, param)
}

func (this *GreaterThan) evaluate() uint64 {
	return fromBool(this.value[0] > this.value[1])
}

//------------------------------------------------------------------------------

type EqualTo struct {
	value []uint64
}

func MakeEqualTo() *EqualTo {
	return &EqualTo{value: make([]uint64, 0)}
}

func (this *EqualTo) addParam(param uint64) {
	this.value = append(this.value, param)
}

func (this *EqualTo) evaluate() uint64 {
	return fromBool(this.value[0] == this.value[1])
}

func fromBool(x bool) uint64 {
	if x {
		return 1
	}
	return 0
}

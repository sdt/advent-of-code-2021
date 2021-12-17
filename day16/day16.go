package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	bitstream := getInput(filename)

	part1 := uint64(0)
	depth := 0
	parseBitstream(&bitstream, &ParseConfig{
		onLiteral:func (version, value uint64) {
			part1 += version
			//fmt.Printf("Literal: version=%d value=%d depth=%d\n", version, value, depth)
		},
		onBeginOperator:func (version, typeId uint64) {
			part1 += version
			//fmt.Printf("Operator: version=%d type=%d depth=%d->%d\n", version, typeId, depth, depth+1)
			depth++
		},
		onEndOperator:func () {
			//fmt.Printf("End operator: depth=%d->%d\n", depth, depth-1)
			depth--
		}})
	fmt.Println(part1)
}

func getInput(filename string) Bitstream {
	lines := aoc.GetInputLines(filename)
	return MakeBitstream(lines[0])
}

//------------------------------------------------------------------------------

const Literal = 4

type OnLiteral func (version, value uint64)
type OnBeginOperator func (version, typeId uint64)
type OnEndOperator func ()

type ParseConfig struct {
	onLiteral 		OnLiteral
	onBeginOperator OnBeginOperator
	onEndOperator	OnEndOperator
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
		ret = (ret << 4) | (chunk & 0xf);
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
	data []uint8
	byteCursor int
	bitCursor int
}

func (this *Bitstream) BitPosition() uint64 {
	return uint64(this.byteCursor * 8 + this.bitCursor)
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
		if i & 1 == 0 {
			nybble = parseNybble(hexchar)
		} else {
			data[i >> 1] = (nybble << 4) | parseNybble(hexchar)
		}
	}

	return Bitstream{data:data, byteCursor:0, bitCursor:0}
}

func parseNybble(hexchar rune) uint8 {
	if hexchar >= '0' && hexchar <= '9' {
		return uint8(hexchar - '0')
	} else {
		return 10 + uint8(hexchar - 'A')
	}
}

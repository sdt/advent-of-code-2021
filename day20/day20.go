package main

import (
	"advent-of-code/aoc"
	"fmt"
	"math"
)

func main() {
	image, enhancer := ParseInput(aoc.GetFilename())

	fmt.Println(part1(image, enhancer))
	fmt.Println(part2(image, enhancer))
}

//------------------------------------------------------------------------------

func part1(image Image, enhancer Enhancer) int {
	image = image.EnhanceTwice(enhancer, 2)

	return len(image.pixmap)
}

func part2(image Image, enhancer Enhancer) int {
	for i := 0; i < 25; i++ {
		image = image.EnhanceTwice(enhancer, 2)
	}

	return len(image.pixmap)
}

//------------------------------------------------------------------------------

type Pixel byte

const (
	Off Pixel = '.'
	On        = '#'
)

type Coord struct {
	x, y int
}

type Pixmap map[Coord]bool

func (this *Pixmap) SetPixel(x, y int, pixel Pixel) {
	if pixel == On {
		(*this)[Coord{x: x, y: y}] = true
	}
}

func (this *Pixmap) GetPixel(x, y int) Pixel {
	if _, found := (*this)[Coord{x: x, y: y}]; found {
		return On
	}
	return Off
}

func (this *Pixmap) EnhancePixel(x, y int, enhancer Enhancer) Pixel {
	address := 0

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			address = address << 1
			if this.GetPixel(x+dx, y+dy) == On {
				address |= 1
			}
		}
	}

	return enhancer[address]
}

func (this *Pixmap) EnhancePixelTwice(x, y int, enhancer Enhancer) Pixel {
	pixmap := make(Pixmap)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			pixmap.SetPixel(x+dx, y+dy,
				this.EnhancePixel(x+dx, y+dy, enhancer))
		}
	}

	return pixmap.EnhancePixel(x, y, enhancer)
}

func (this Pixmap) Extents() (Coord, Coord) {
	min := Coord{x: math.MaxInt, y: math.MaxInt}
	max := Coord{x: math.MinInt, y: math.MinInt}

	for coord, _ := range this {
		if coord.x < min.x {
			min.x = coord.x
		} else if coord.x > max.x {
			max.x = coord.x
		}

		if coord.y < min.y {
			min.y = coord.y
		} else if coord.y > max.y {
			max.y = coord.y
		}
	}

	return min, max
}

type Image struct {
	pixmap   Pixmap
	min, max Coord
}

type Enhancer []Pixel

func (this Image) Enhance(enhancer Enhancer, extra int) Image {
	pixmap := make(Pixmap)

	for y := this.min.y - extra; y <= this.max.y+extra; y++ {
		for x := this.min.x - extra; x <= this.max.x+extra; x++ {
			pixmap.SetPixel(x, y, this.pixmap.EnhancePixel(x, y, enhancer))
		}
	}
	min, max := pixmap.Extents()
	return Image{pixmap: pixmap, min: min, max: max}
}

func (this Image) EnhanceTwice(enhancer Enhancer, extra int) Image {
	pixmap := make(Pixmap)

	for y := this.min.y - extra; y <= this.max.y+extra; y++ {
		for x := this.min.x - extra; x <= this.max.x+extra; x++ {
			pixmap.SetPixel(x, y, this.pixmap.EnhancePixelTwice(x, y, enhancer))
		}
	}
	min, max := pixmap.Extents()
	return Image{pixmap: pixmap, min: min, max: max}
}

func (this Image) Print(msg string) {
	fmt.Println(msg)
	for y := this.min.y - 1; y <= this.max.y+1; y++ {
		for x := this.min.x - 1; x <= this.max.x+1; x++ {
			fmt.Printf("%c", this.pixmap.GetPixel(x, y))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func ParseInput(filename string) (Image, Enhancer) {
	lines := aoc.GetInputLines(aoc.GetFilename())

	enhancer := ParseEnhancer(lines[0])
	image := ParseImage(lines[2:])

	return image, enhancer
}

func ParseImage(lines []string) Image {
	pixmap := make(Pixmap)
	for y, line := range lines {
		for x, char := range line {
			pixmap.SetPixel(x, y, Pixel(char))
		}
	}
	min, max := pixmap.Extents()

	return Image{pixmap: pixmap, min: min, max: max}
}

func ParseEnhancer(line string) Enhancer {
	enhancer := make([]Pixel, len(line))
	for i, char := range line {
		enhancer[i] = Pixel(char)
	}
	return enhancer
}

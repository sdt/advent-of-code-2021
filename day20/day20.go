package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	image, enhancer := ParseInput(aoc.GetFilename())

	fmt.Println(part1(image, enhancer))
	fmt.Println(part2(image, enhancer))
}

//------------------------------------------------------------------------------

func part1(image Image, enhancer Enhancer) int {

	for i := 0; i < 2; i++ {
		image = image.Enhance(enhancer)
	}

	return image.LitPixels()
}

func part2(image Image, enhancer Enhancer) int {

	for i := 0; i < 50; i++ {
		image = image.Enhance(enhancer)
	}

	return image.LitPixels()
}

//------------------------------------------------------------------------------

type Pixel byte

const (
	Off Pixel = '.'
	On        = '#'
)

func (this Pixel) Flip() Pixel {
	if this == Off {
		return On
	} else {
		return Off
	}
}

type Image struct {
	pixel        []Pixel
	w, h         int
	defaultPixel Pixel
}

type Enhancer []Pixel

func (this Image) Print(msg string) {
	fmt.Println(msg)
	i := 0
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			fmt.Printf("%c", this.pixel[i])
			i++
		}
		fmt.Print("\n")
	}
	fmt.Println(this.defaultPixel, this.LitPixels())
	fmt.Print("\n")
}

func (this Image) GetPixel(x, y int) Pixel {
	if (x < 0) || (x >= this.w) ||
	   (y < 0) || (y >= this.h) {
	   return this.defaultPixel
	}

	offset := y * this.w + x
	return this.pixel[offset]
}

func (this Image) EnhancePixel(x, y int, enhancer Enhancer) Pixel {
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

func (this Image) Enhance(enhancer Enhancer) Image {
	w := this.w + 2
	h := this.h + 2

	pixels := make([]Pixel, w * h)
	i := 0

	for y := -1; y <= this.h; y++ {
		for x := -1; x <= this.w; x++ {
			pixels[i] = this.EnhancePixel(x, y, enhancer)
			i++
		}
	}

	return Image{pixels, w, h, this.defaultPixel.Flip()}
}

func (this Image) LitPixels() int {
	count := 0

	for _, pixel := range this.pixel {
		if pixel == On {
			count++
		}
	}

	return count
}

func ParseInput(filename string) (Image, Enhancer) {
	lines := aoc.GetInputLines(aoc.GetFilename())

	enhancer := ParseEnhancer(lines[0])
	image := ParseImage(lines[2:])

	return image, enhancer
}

func ParseImage(lines []string) Image {
	w := len(lines[0])
	h := len(lines)

	pixels := make([]Pixel, w * h)
	i := 0

	for _, line := range lines {
		for _, pixel := range line {
			pixels[i] = Pixel(pixel)
			i++
		}
	}

	return Image{pixels, w, h, Off}
}

func ParseEnhancer(line string) Enhancer {
	enhancer := make([]Pixel, len(line))
	for i, char := range line {
		enhancer[i] = Pixel(char)
	}
	return enhancer
}

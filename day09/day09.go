package main

import (
	"advent-of-code/common"
	"fmt"
)

type HeightMap struct {
	rows, cols int
	height []int
}

func main() {
	filename := common.GetFilename()
	heightMap := parseHeightMap(filename)

	fmt.Println(part1(&heightMap))
}

func part1(h *HeightMap) int {
	total := 0
	for row := 0; row < h.rows; row++ {
		for col := 0; col < h.cols; col++ {
			total += h.RiskLevel(row, col)
		}
	}
	return total
}

func parseHeightMap(filename string) HeightMap {
	lines := common.GetInputLines(filename)
	rows := len(lines)
	cols := len(lines[0])
	heightMap := HeightMap{ rows:rows, cols:cols, height: make([]int, rows * cols) }

	index := 0
	for _, line := range(lines) {
		for _, digit := range(line) {
			heightMap.height[index] = int(digit - '0')
			index++
		}
	}

	return heightMap
}

func (h *HeightMap) Height(row, col int) int {
	if row < 0 || col < 0 || row >= h.rows || col >= h.cols {
		return 9
	}

	index := row * h.cols + col
	return h.height[index]
}

func (h *HeightMap) IsLowPoint(row, col int) bool {
	height := h.Height(row, col)
	return height < h.Height(row - 1, col) && height < h.Height(row + 1, col) && height < h.Height(row, col - 1) && height < h.Height(row, col + 1)
}

func (h *HeightMap) RiskLevel(row, col int) int {
	if !h.IsLowPoint(row, col) {
		return 0
	}

	return h.Height(row, col) + 1
}

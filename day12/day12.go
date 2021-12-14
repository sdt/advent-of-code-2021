package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strings"
)

const (
	Start     = 0
	End       = 1
	Double    = 2
	FirstCave = 3
)

type Cave int32

type Caves struct {
	to      [][]Cave
	isLarge Cave
}

func MakeCaves(count int, isLarge Cave) Caves {
	caves := Caves{to: make([][]Cave, count), isLarge: isLarge}
	for i := 0; i < count; i++ {
		caves.to[i] = make([]Cave, 0)
	}
	return caves
}

type CaveMap struct {
	mapping map[string]Cave
	next    Cave
}

func MakeCaveMap() CaveMap {
	cavemap := CaveMap{mapping: make(map[string]Cave), next: FirstCave}
	cavemap.mapping["start"] = Start
	cavemap.mapping["end"] = End
	return cavemap
}

type Path struct {
	seen int
	from Cave
}

func MakePath() Path {
	return Path{seen: 0, from: Start}
}

func main() {
	caves := parseCaves(aoc.GetFilename())

	fmt.Println(part1(&caves))
	fmt.Println(part2(&caves))
}

func part1(caves *Caves) int {
	return walk1(MakePath(), caves)
}

func walk1(path Path, caves *Caves) int {
	paths := 0

	for _, to := range caves.To(path.from) {
		if to == End {
			paths++
		} else if caves.IsLarge(to) || !path.Seen(to) {
			next := path.Extend(to)
			paths += walk1(next, caves)
		}
	}
	return paths
}

func part2(caves *Caves) int {
	return walk2(MakePath(), caves)
}

func walk2(path Path, caves *Caves) int {
	paths := 0

	for _, to := range caves.To(path.from) {
		if to == End {
			paths++
		} else if caves.IsLarge(to) || !path.Seen(to) {
			next := path.Extend(to)
			paths += walk2(next, caves)
		} else if !path.Seen(Double) {
			next := path.Extend(to)
			next.SeenDouble()
			paths += walk2(next, caves)
		}
	}
	return paths
}

func parseCaves(filename string) Caves {
	cavemap := MakeCaveMap()
	lines := aoc.GetInputLines(filename)
	isLarge := Cave(0)
	for _, line := range lines {
		for _, end := range strings.Split(line, "-") {
			value := cavemap.AssignValue(end)
			if strings.ToUpper(end) == end { // ffs...
				isLarge |= 1 << value
			}
		}
	}

	caves := MakeCaves(int(cavemap.next), isLarge)
	for _, line := range lines {
		ends := strings.Split(line, "-")
		a := cavemap.Lookup(ends[0])
		b := cavemap.Lookup(ends[1])
		caves.AddPath(a, b)
		caves.AddPath(b, a)
	}
	return caves
}

func (this *Caves) AddPath(from, to Cave) {
	if (from != End) && (to != Start) {
		this.to[from] = append(this.to[from], to)
	}
}

func (this *Caves) To(from Cave) []Cave {
	return this.to[from]
}

func (this *Caves) IsLarge(cave Cave) bool {
	return (this.isLarge & (1 << cave)) != 0
}

func (this *CaveMap) AssignValue(cavename string) Cave {
	value, found := this.mapping[cavename]
	if found {
		return value
	}
	value = this.next
	this.mapping[cavename] = value
	this.next++
	return value
}

func (this *CaveMap) Lookup(cave string) Cave {
	if value, found := this.mapping[cave]; found {
		return value
	}
	return -1
}

func (this *Path) Extend(cave Cave) Path {
	return Path{seen: this.seen | 1<<cave, from: cave}
}

func (this *Path) SeenDouble() {
	this.seen |= 1 << Double
}

func (this *Path) Seen(cave Cave) bool {
	return (this.seen & (1 << cave)) != 0
}

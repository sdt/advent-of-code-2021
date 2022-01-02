package main

import (
	"advent-of-code/aoc"
	"fmt"
)

const Width = 13
const XMax = Width-1
const Height = 5
const YMax = Height-1

const HallY = 1
const Home1Y = 2
const Home0Y = 3

type Coord int8
type Index uint8
type Amphipod uint8

type Point struct {
	index Index
}

func MakePoint(x, y Coord) Point {
	return Point{index: Index(y << 4 | x)}
}

func (point Point) x() Coord {
	return Coord(point.index & 0xf)
}

func (point Point) y() Coord {
	return Coord((point.index >> 4) & 0xf)
}

func (point Point) move(dir Point) Point {
	return MakePoint(point.x() + dir.x(), point.y() + dir.y())
}

func (point Point) up() Point {
	return Point{index: point.index - 16}
}

func (point Point) down() Point {
	return Point{index: point.index + 16}
}

func (point Point) left() Point {
	return Point{index: point.index - 1}
}

func (point Point) right() Point {
	return Point{index: point.index + 1}
}

func (point Point) String() string {
	return fmt.Sprintf("%d,%d", point.x(), point.y())
}

const (
	A Amphipod = iota
	B
	C
	D
)

/*
	Cell indexes
#############					*---> +X
#0..........#	 0 .. 10		|
###B#C#B#D###	11 .. 14		|
  #A#D#C#A#		15 .. 18		V +Y
  #########
   2 4 6 8
               111
   X 0123456789012
Y 0  #############
  1  #...........#
  2  ###B#C#B#D###
  3    #A#D#C#A#
  4    #########
*/

type State struct {
	point [8]Point
}

type EnergyLevel int

type StateTracker struct {
	best map[State]EnergyLevel
	queue []*Path
}

func MakeStateTracker() StateTracker {
	return StateTracker{best: make(map[State]EnergyLevel)}
}

type Path struct {
	state State
	energy EnergyLevel
}

func (this *Path) Extend(index int, pos Point, distance int) Path {
	state := State{}
	for i, value := range this.state.point {
		state.point[i] = value
	}
	state.point[index] = pos
	state.Normalise()

	amphipod, _ := MakeAmphipod(index)
	energy := this.energy + EnergyLevel(distance * energyFactor[amphipod])

	return Path{state: state, energy: energy}
}

func (this *StateTracker) IsHigherPriority(i, j int) bool {
	return this.queue[i].energy < this.queue[j].energy
}

func (this *StateTracker) IsValid(i int) bool {
	return (i >= 0) && (i < len(this.queue))
}

func (this *StateTracker) Swap(i, j int) {
	this.queue[i], this.queue[j] = this.queue[j], this.queue[i]
}

func (this *StateTracker) Add(state State, energy EnergyLevel) {
	if existingEnergy, found := this.best[state]; found {
		if existingEnergy < energy {
			return	// we already have a better path to this state
		}
	}

	this.best[state] = energy
	this.queue = append(this.queue, &Path{state, energy})
	upheap(this, len(this.queue)-1)
}

func (this *StateTracker) Next() (path *Path, found bool) {
	for len(this.queue) > 0 {
		path := this.queue[0]
		size := len(this.queue)
		this.queue[0] = this.queue[size-1]
		this.queue = this.queue[0:size-1]
		downheap(this, 0)

		if path.energy == this.best[path.state] {
			return path, true
		}

		// A better path to this state has already been tried, keep looking
	}
	return nil, false
}

/*
#############
#ab.c.d.e.fg#
###.#.#.#.###
  #.#.#.#.#
  #########
*/
var outDestinations = []Point{
	MakePoint(1,  HallY),	// a
	MakePoint(2,  HallY),	// b
	MakePoint(4,  HallY),	// c
	MakePoint(6,  HallY),	// d
	MakePoint(8,  HallY),	// e
	MakePoint(10, HallY),	// f
	MakePoint(11, HallY),	// f
}

/*
#############
#...........#
###a#c#e#g###
  #b#d#f#h#
  #########
*/
var home = [][]Point {
	[]Point{ MakePoint(3, Home0Y), MakePoint(3, Home1Y) },	// b a
	[]Point{ MakePoint(5, Home0Y), MakePoint(5, Home1Y) },	// d c
	[]Point{ MakePoint(7, Home0Y), MakePoint(7, Home1Y) },	// f e
	[]Point{ MakePoint(9, Home0Y), MakePoint(9, Home1Y) },	// h g
}

var energyFactor = []int{ 1, 10, 100, 1000 }

type Burrow struct {
	cell [Height * Width]byte
}

func (this *State) MakeBurrow() Burrow {
	burrow := Burrow{}

	for i, point := range this.point {
		burrow.cell[point.index] = byte('A') + byte(i >> 1)
	}

	return burrow
}

func (this *State) IsHome(index int) bool {
	amphipod, which := MakeAmphipod(index)

	//fmt.Printf("Testing %c[%d] %v\n", 'A' + amphipod, which, this.point[index])

	if this.point[index].x() != home[amphipod][which].x() {
		return false
	}

	if this.point[index].y() == Home0Y {
		return true
	}

	if this.point[index].y() != Home1Y {
		return false
	}

	return this.IsHome(index ^ 1)
}

func (this *Burrow) CanReach(from, to Point) (int, bool) {
	steps := 0

	// First scan up to the hallway
	for from.y() > HallY {
		from = from.up()
		steps++
		if !this.IsClear(from) {
			return steps, false
		}
	}

	// Next scan across to the matching column
	for from.x() < to.x() {
		from = from.right()
		steps++
		if !this.IsClear(from) {
			return steps, false
		}
	}

	for from.x() > to.x() {
		from = from.left()
		steps++
		if !this.IsClear(from) {
			return steps, false
		}
	}

	// Finally scan down
	for from.y() < to.y() {
		from = from.down()
		steps++
		if !this.IsClear(from) {
			return steps, false
		}
	}

	return steps, true
}

func (this *Burrow) IsClear(p Point) bool {
	return this.cell[p.index] == 0
}

func sign(x Coord) Coord {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

/*
type Cell struct {
	pos 		Point
	cellType 	CellType
	neighbours 	[]CellIndex
}

type CellPosTable map[Point]CellIndex
var cellPosTable CellPosTable = make(CellPosTable)

var cellTable []Cell = []Cell{
	MakeCell( 1, 1, Hall,     1),			//  0
	MakeCell( 2, 1, Hall,     0,  2),		//  1
	MakeCell( 3, 1, Doorway,  1,  3, 11),	//  2
	MakeCell( 4, 1, Hall,     2,  4),		//  3
	MakeCell( 5, 1, Doorway,  3,  5, 12),	//  4
	MakeCell( 6, 1, Hall,     4,  6),		//  5
	MakeCell( 7, 1, Doorway,  5,  7, 13),	//  6
	MakeCell( 8, 1, Hall,     6,  8),		//  7
	MakeCell( 9, 1, Doorway,  7,  9, 14),	//  8
	MakeCell(10, 1, Hall,     8, 10),		//  9
	MakeCell(11, 1, Hall,     9),			// 10

	MakeCell( 3, 2, HomeA0,   2, 15),		// 11
	MakeCell( 5, 2, HomeB0,   4, 16),		// 12
	MakeCell( 7, 2, HomeC0,   6, 17),		// 13
	MakeCell( 9, 2, HomeD0,   8, 18),		// 14

	MakeCell( 3, 3, HomeA1,   11),			// 15
	MakeCell( 5, 3, HomeB1,   12),			// 16
	MakeCell( 7, 3, HomeC1,   13),			// 17
	MakeCell( 9, 3, HomeD1,   14),			// 18
}

func MakeCell(x, y int, cellType CellType, neighbours ...CellIndex) Cell {
	return Cell{pos: Point{x: x, y: y}, cellType: cellType, neighbours: neighbours}
}
*/

func main() {
	state := ParseMaze(aoc.GetFilename())
	tracker := MakeStateTracker()

	tracker.Add(state, 0)

	for {
		path, found := tracker.Next();
		if !found {
			panic("No more states")
		}
		//path.state.Print()
		//fmt.Printf("Energy = %d\n\n", path.energy)

		burrow := path.state.MakeBurrow()

		homeCount := 0
		for i, pos := range path.state.point {
			amphipod, _ := MakeAmphipod(i)

			if path.state.IsHome(i) {
				homeCount++
				continue // already home, no need to go anywhere
			}

			destinations := []Point{}

			if burrow.IsClear(home[amphipod][0]) {
				destinations = append(destinations, home[amphipod][0])
			} else {
				destinations = append(destinations, home[amphipod][1])
			}

			if pos.y() != HallY {
				// Outbound
				destinations = append(destinations, outDestinations...)
			}

			for _, to := range destinations {
				if distance, canReach := burrow.CanReach(pos, to); canReach {
					//fmt.Printf("%c can move from %v to %v (energy=%d)\n", 'A' + amphipod, pos, to, distance * energyFactor[amphipod])

					next := path.Extend(i, to, distance)
					tracker.Add(next.state, next.energy)
				}
			}
		}

		if homeCount == 8 {
			fmt.Println("Solved!")
			fmt.Println("Energy=", path.energy)
			return
		}
	}
}

func MakeAmphipod(index int) (Amphipod, int) {
	return Amphipod(index / 2), index & 1
}

func ParseMaze(filename string) State {
	lines := aoc.GetInputLines(filename)

	state := State{}
	var count [4]int

	for y, line := range lines {
		for x, char := range line {
			if char >= 'A' && char <= 'D' {
				amphipod := int(A) + int(char - 'A')
				index := amphipod * 2 + count[amphipod]
				count[amphipod]++

				state.point[index] = MakePoint(Coord(x), Coord(y))
			}
		}
	}

	state.Normalise()

	return state
}

// For each pair of AA BB CC DD, the lhs should have a lower index than rhs
// Just so that 0123 matches 1032
func (this *State) Normalise() {
	for i := 0; i < 8; i += 2 {
		if this.point[i].index > this.point[i+1].index {
			this.point[i], this.point[i+1] = this.point[i+1], this.point[i]
		}
	}
}

func (this *State) Print() {
	for y := Coord(0); y < Height; y++ {
		for x := Coord(0); x < Width; x++ {
			point := MakePoint(x, y)
			fmt.Printf("%c", getMazePoint(this, point))
		}
		fmt.Print("\n")
	}
}

var baseMaze [5]string = [...]string{
	"#############",
	"#...........#",
	"###.#.#.#.###",
  	"  #.#.#.#.#  ",
  	"  #########  ",
}

func getMazePoint(state *State, point Point) uint8 {
	for i := 0; i < len(state.point); i++ {
		if state.point[i] == point {
			return uint8('A' + (i >> 1))
		}
	}

	return baseMaze[point.y()][point.x()]
}

type Heap interface {
	IsHigherPriority(parent, child int) bool
	IsValid(i int) bool
	Swap(i, j int)
}

func upheap(heap Heap, child int) {
	for child > 0 {
		parent := (child - 1) / 2

		if heap.IsHigherPriority(parent, child) {
			return
		}

		heap.Swap(parent, child)
		child = parent
	}
}

func downheap(heap Heap, parent int) {
	for {
		lchild := parent * 2 + 1
		if !heap.IsValid(lchild) {
			return
		}

		rchild := parent * 2 + 2
		var child int
		if !heap.IsValid(rchild) || heap.IsHigherPriority(lchild, rchild) {
			child = lchild
		} else {
			child = rchild
		}

		if heap.IsHigherPriority(parent, child) {
			return
		}

		heap.Swap(parent, child)
		parent = child
	}
}

package main

import (
	"advent-of-code/aoc"
	"fmt"
)

type Display byte

type Amphipod uint8	// 5 values -> 3 bits wide (use 4 to match RoomState)

var energyFactor []Energy = []Energy{ 0, 1, 10, 100, 1000 }

func (this Amphipod) energy(distance int) Energy {
	if this == Empty {
		panic("EMPTY ENERGY?")
	}
	return Energy(distance) * energyFactor[this]
}

func (this Amphipod) roomIndex() int {
	return int(this) - 1
}

func (this Amphipod) hallwayPos() int {
	return this.roomIndex() * 2 + 2
}

const (
	Empty Amphipod = iota
	Amber
	Bronze
	Copper
	Desert
)

func (amphipod Amphipod) Display() Display {
	if amphipod == Empty {
		return '.'
	}
	return 'A' + Display(amphipod) - 1
}

func (display Display) Amphipod() Amphipod {
	return Amphipod(display - 'A' + 1)
}

type Energy uint32

//------------------------------------------------------------------------------

// Say our room is the A room, and starts like this: #BACD
// There's only nine states this room can be in:
//
//	            CanLeave CanEnter Final Distance
//	0: #BACD.	T        -        -     1
//  1: #BAC..	T        -        -     2
//  2: #BA...   T        -        -     3
//  3: #B....   T        -        -     4
//  4: #.....   -        T        -     4
//  5: #A....   -        T        -     3
//  6: #AA...   -        T        -     2
//  7: #AAA..   -        T        -     1
//  8: #AAAA.   -        -        T     -
//
// For a room that starts partially complete, we can build the states as if the
// leading matches weren't there. For the #AABC. case, this works out the same
// as if we had just #BC. The states and the distances match.
//
//	            CanLeave CanEnter Final Distance
//  0: #AABC.   T        -        -     1
//  1: #AAB..   T        -        -     2
//  2: #AA...   -        T        -     2
//  3: #AAA..   -        T        -     1
//  4: #AAAA.   -        -        T     -
//
// Final == !( CanLeave || CanEnter )

// offset freeSlots before after total
// 0      4			4      4     9
// 1      3         3      3     7
// 2      2         2      2     5
// 3      1         1      1     3
// 4      0         0      0     1

type RoomState uint8	// 9 of these -> 4 bits

type Room struct {
	roomType  Amphipod		// which is the home type for this room
	amphipod  []Amphipod	// #BACD. -> DCAB - index matches state
	freeSlots int			// how many amphipods start in the home position
}

func NewRoom(roomType Amphipod, amphipod []Amphipod) Room {
	var freeSlots int
	for freeSlots = len(amphipod); freeSlots > 0; freeSlots-- {
		if amphipod[freeSlots - 1] != roomType {
			break
		}
	}

	return Room{roomType, amphipod, freeSlots}
}

func (this Room) CanLeave(state RoomState) bool {
	return int(state) < this.freeSlots
}

func (this Room) CanEnter(state RoomState) bool {
	return int(state) >= this.freeSlots && int(state) < (this.freeSlots * 2)
}

func (this Room) IsFinal(state RoomState) bool {
	return int(state) == this.freeSlots * 2
}

func (this Room) FinalState() RoomState {
	return RoomState(this.freeSlots * 2)
}

func (this Room) LeaveDistance(state RoomState) int {
	// Assuming can leave. ie state in [0, freeSlots)
	return int(state) + 1
}

func (this Room) EnterDistance(state RoomState) int {
	// Assuming can enter. ie state in [freeSlots, freeSlots*2)
	return this.freeSlots * 2 - int(state)
}

func (this Room) Amphipod(state RoomState) Amphipod {
	// Assuming can leave. ie state in [0, freeSlots]
	return this.amphipod[state]
}

func (this Burrow) String(burrowState BurrowState) string {
	s := ""
	for i, roomType := range this.hallway.roomType {
		if roomType == Empty {
			amphipod := burrowState.GetAmphipod(i)
			s += fmt.Sprintf("%c", amphipod.Display())
		} else {
			roomState := burrowState.GetRoomState(i)
			s += fmt.Sprintf("%d", roomState)
		}
	}
	return s
}

func (this Burrow) Draw(burrowState BurrowState) {
	for i, roomType := range this.hallway.roomType {
		if roomType == Empty {
			amphipod := burrowState.GetAmphipod(i)
			fmt.Printf("%c", amphipod.Display())
		} else {
			roomState := burrowState.GetRoomState(i)
			fmt.Printf("%d", roomState)
		}
	}
	fmt.Print("\n")
	/*
	for pos, roomType := range this.hallway.roomType {
		if roomType == Empty {
			amphipod := burrowState.GetAmphipod(pos)
			fmt.Printf("%c", amphipod.Display())
		} else {
			fmt.Print(".")
		}
	}
	fmt.Print("\n")

	for row := len(this.room[0].amphipod) - 1; row >= 0; row-- {
		for pos, roomType := range this.hallway.roomType {
			if roomType == Empty {
				fmt.Print("#")
			} else {
				roomIndex := roomType.roomIndex()
				roomState := burrowState.GetRoomState(pos)
				fmt.Printf("%c", this.room[roomIndex].Draw(roomState, row))
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
	*/
}


// #3210
// #DCBA freeslots = 4
// State:   012345678 orig	new
// Row 0/3: A.......A
// Row 1/2: BB.....AA S < F-2		S > F+2
// Row 2/1: CCC...AAA S < F-X		S > F+X
// Row 3/0: DDDD.AAAA S < F-X			S > F
//
// #AADC
// State: 01234
// Row 0: AAAAA
// Row 1: AAAAA
// Row 2: DD.AA
// Row 3: C...A

func (this Room) Draw(state RoomState, row int) Display {
	xrow := len(this.amphipod) - row - 1

	used := len(this.amphipod) - this.freeSlots - 1
	if int(state) < used - row {
		return this.amphipod[xrow].Display()
	}
	if int(state) > used + row {
		return this.roomType.Display()
	}
	return '.'
}

//------------------------------------------------------------------------------

// Hallway has 7 valid positions, and four rooms.
//  01x2x3x4x56
// Each position contains an amphipod (or empty), which means 7x3=21 bits
// Each rooms needs 4 bits, so that's 4x4 + 7x3 = 16 + 21 = 37 bits of state

/*
    Pos		Type
	  0		hall
	  1		hall
	  2		room A
	  3		hall
	  4		room B
	  5		hall
	  6		room C
	  7		hall
	  8		room D
	  9		hall
	 10		hall
*/


type Hallway struct {
	roomType  []Amphipod
	roomIndex []int
}

func NewHallway() Hallway {
	// Squares default to Empty - ie no Roomtype
	roomType := make([]Amphipod, 11)

	// Rooms are at ..2.4.6.8..
	for i := 0; i < 4; i++ {
		roomType[i * 2 + 2] = Amber + Amphipod(i)
	}

	roomIndex, hallIndex := 0, 0
	index := make([]int, 11)
	for i, roomType := range roomType {
		if roomType == Empty {
			index[i] = hallIndex
			hallIndex++
		} else {
			index[i] = roomIndex
			roomIndex++
		}
	}

	return Hallway{roomType, index}
}

//------------------------------------------------------------------------------

/*
	Offset	Len		Descr
	0		4		Hallway 0 state
	4		4		Hallway 1 state
	8		4		Room 0 state
	12		4		Hallway 2 state
	16		4		Room 1 state
	20		4		Hallway 3 state
	24		4		Room 2 state
	28		4		Hallway 4 state
	32		4		Room 3 state		// offset = room * 4
	36		4		Hallway 5 state
	40		4		Hallway 6 state		// offset = hallway * 3 + 16

	44 bits total
*/

type BurrowState uint64

func (burrowState BurrowState) ClearAmpipod(position int) BurrowState {
	shift := position * 4
	mask := BurrowState(0xf) << shift
	return burrowState & ^mask
}

func (burrowState BurrowState) SetAmpipod(position int, amphipod Amphipod) BurrowState {
	shift := position * 4
	mask := BurrowState(0xf) << shift
	return (burrowState & ^mask) | BurrowState(amphipod) << shift
}

func (burrowState BurrowState) NextRoomState(position int) BurrowState {
	shift := position * 4
	incr := BurrowState(1) << shift
	return burrowState + incr
}

func (burrowState BurrowState) GetRoomState(position int) RoomState {
	shift := position * 4
	return RoomState((burrowState >> shift) & 0xf)
}

func (burrowState BurrowState) SetRoomState(position int, roomState RoomState) BurrowState {
	shift := position * 4
	mask := BurrowState(0xf) << shift
	return (burrowState & ^mask) | (BurrowState(roomState) << shift)
}

func (burrowState BurrowState) GetAmphipod(position int) Amphipod {
	shift := position * 4
	return Amphipod((burrowState >> shift) & 0xf)
}

//------------------------------------------------------------------------------

type Burrow struct {
	room       []Room
	hallway    Hallway
	finalState BurrowState
}

func NewBurrow(filename string, isPart2 bool) Burrow {
	amphipod := make([][]Amphipod, 4)
	for i, _ := range amphipod {
		amphipod[i] = []Amphipod{}
	}

	parseLine := func (line string) []byte {
		value := make([]byte, 4)
		for i := 0; i < 4; i++ {
			value[i] = line[i * 2 + 3]
		}
		return value
	}

	lines := aoc.GetInputLines(filename)
	rows := [][]byte{}

	rows = append(rows, parseLine(lines[2]))
	if isPart2 {
		rows = append(rows, []byte{ 'D', 'C', 'B', 'A' })
		rows = append(rows, []byte{ 'D', 'B', 'A', 'C' })
	}
	rows = append(rows, parseLine(lines[3]))

	for _, row := range rows {
		for i, value := range row {
			amphipod[i] = append(amphipod[i], Display(value).Amphipod())
		}
	}

	room := []Room{}
	for roomType, occupants := range amphipod {
		room = append(room, NewRoom(Amphipod(roomType+1), occupants))
	}

	hallway := NewHallway()
	finalState := BurrowState(0)

	for pos, roomType := range hallway.roomType {
		if roomType != Empty {
			roomIndex := roomType - 1
			roomState := room[roomIndex].FinalState()
			finalState = finalState.SetRoomState(pos, roomState)
		}
	}

	return Burrow{room, hallway, finalState}
}

func (this Burrow) Solve() Energy {
	solution := NewBurrowSolution()

	attempts := 0

	/*
	this.Draw(0x00000000000)
	this.Draw(0x00101010100)
	this.Draw(0x00202020200)
	return Energy(0)
	*/

	for {
		path, found := solution.Next()

		if !found {
			//fmt.Printf("No solution found after %d attempts\n", attempts)
			panic("whoops")
		}
		attempts++

		//fmt.Printf("\nAttempt %d: trying %d %s\n", attempts, path.energy, this.String(path.burrowState))
		//this.Draw(path.burrowState)

		if this.IsSolution(path.burrowState) {
			fmt.Printf("Solution: %d after %d attempts\n", path.energy, attempts)
			return path.energy
		}

		for pos, roomType := range this.hallway.roomType {
			//display := amphipodToDisplay(roomType)
			if roomType == Empty {
				// This is a regular hallway square
				amphipod := path.burrowState.GetAmphipod(pos)
				if amphipod == Empty {
					//fmt.Printf("Pos %d is empty hallway\n", pos)
					continue // nothing in this square
				}
				//fmt.Printf("Pos %d is hallway with a %c\n", pos, amphipod.Display())
				this.TryMoveFromHallway(&solution, path, pos, amphipod)
			} else {
				// This is a room
				//fmt.Printf("Pos %d is a %c room\n", pos, roomType.Display())
				this.TryMoveFromRoom(&solution, path, pos, roomType)
			}
		}
	}
}

func (this Burrow) IsSolution(burrowState BurrowState) bool {
	return burrowState == this.finalState
}

func (this Burrow) TryMoveFromHallway(burrowSolution *BurrowSolution, path *Path, pos int, amphipod Amphipod) {
	// We are looking at an amphipod in the hallway. It can only move home, and
	// it can only do so if it's home room CanEnter, and if the route is clear.

	roomIndex := amphipod.roomIndex()
	roomPos := amphipod.hallwayPos()
	roomState := path.burrowState.GetRoomState(roomPos)

	if !this.room[roomIndex].CanEnter(roomState) {
		//fmt.Printf("Cannot enter room %c at %d\n", amphipod.Display(), amphipod.hallwayPos())
		return
	}

	// The room is suitable to enter. Can we get there?
	if dist, canMove := this.CanMove(pos, roomPos, path.burrowState); canMove {
		burrowState := path.burrowState.ClearAmpipod(pos)
		burrowState = burrowState.NextRoomState(roomPos)

		dist += this.room[roomIndex].EnterDistance(roomState)

		//fmt.Printf("Move %c at %d into room at %d: %s -> %s", amphipod.Display(), pos, roomPos, this.String(path.burrowState), this.String(burrowState))
		burrowSolution.Add(burrowState, path, amphipod.energy(dist))
	}
}

func (this Burrow) TryMoveFromRoom(burrowSolution *BurrowSolution, path *Path, pos int, fromRoomType Amphipod) {

	fromRoomState := path.burrowState.GetRoomState(pos)
	fromRoomIndex := this.hallway.roomIndex[pos]
	if !this.room[fromRoomIndex].CanLeave(fromRoomState) {
		return
	}

	amphipod := this.room[fromRoomIndex].Amphipod(fromRoomState)

	burrowState := path.burrowState.NextRoomState(pos)

	toRoomPos   := amphipod.hallwayPos()
	toRoomIndex := this.hallway.roomIndex[pos]
	toRoomState := path.burrowState.GetRoomState(toRoomPos)

	dist := this.room[fromRoomIndex].LeaveDistance(fromRoomState)

	// Scan to the left.
	for leftPos := pos - 1; leftPos >= 0; leftPos-- {
		roomType := this.hallway.roomType[leftPos]
		if roomType == Empty {
			if path.burrowState.GetAmphipod(leftPos) != Empty {
				//fmt.Printf("Left scan %c stops at %d\n", amphipod.Display(), leftPos)
				break // hallway is blocked
			}
			// We are in a hallway and we can stop here
			next := burrowState.SetAmpipod(leftPos, amphipod)
			//fmt.Printf("Left scan %c from %d lands at hallway pos %d: %s -> %s", amphipod.Display(), pos, leftPos, this.String(path.burrowState), this.String(next))
			burrowSolution.Add(next, path, amphipod.energy(dist + pos - leftPos))
		} else if roomType == amphipod {
			// This is our home room
			toRoomState = path.burrowState.GetRoomState(leftPos)
			toRoomIndex = amphipod.roomIndex()

			if !this.room[toRoomIndex].CanEnter(toRoomState) {
				//fmt.Printf("Left scan %c cannot enter room at %d\n", amphipod.Display(), leftPos)
				continue
			}

			next := burrowState.NextRoomState(leftPos)
			//fmt.Printf("Left scan %c from %d lands at in room %d: %s -> %s", amphipod.Display(), pos, leftPos, this.String(path.burrowState), this.String(next))
			burrowSolution.Add(next, path, amphipod.energy(dist + pos - leftPos + this.room[toRoomIndex].EnterDistance(toRoomState)))
		}
	}

	// Scan to the right.
	for rightPos := pos + 1; rightPos < len(this.hallway.roomType); rightPos++ {
		roomType := this.hallway.roomType[rightPos]
		if roomType == Empty {
			if path.burrowState.GetAmphipod(rightPos) != Empty {
				//fmt.Printf("Right scan %c stops at %d\n", amphipod.Display(), rightPos)
				break // hallway is blocked
			}
			// We are in a hallway and we can stop here
			next := burrowState.SetAmpipod(rightPos, amphipod)
			//fmt.Printf("Right scan %c from %d lands at hallway pos %d: %s -> %s", amphipod.Display(), pos, rightPos, this.String(path.burrowState), this.String(next))
			burrowSolution.Add(next, path, amphipod.energy(dist + rightPos - pos))
		} else if roomType == amphipod {
			// This is our home room
			toRoomState = path.burrowState.GetRoomState(rightPos)
			toRoomIndex = amphipod.roomIndex()

			if !this.room[toRoomIndex].CanEnter(toRoomState) {
				//fmt.Printf("Right scan %c cannot enter room at %d\n", amphipod.Display(), rightPos)
				continue
			}

			next := burrowState.NextRoomState(rightPos)
			//fmt.Printf("Right scan %c from %d lands at in room %d: %s -> %s", amphipod.Display(), pos, rightPos, this.String(path.burrowState), this.String(next))
			burrowSolution.Add(next, path, amphipod.energy(dist + rightPos - pos + this.room[toRoomIndex].EnterDistance(toRoomState)))
		}
	}
}

func (this Burrow) CanMove(from, to int, burrowState BurrowState) (int, bool) {
	delta := sign(to - from)
	for pos := from + delta; pos != to; pos += delta {
		if this.hallway.roomType[pos] == Empty {
			amphipod := burrowState.GetAmphipod(pos)
			if amphipod != Empty {
				return 0, false
			}
		}
	}

	return abs(to - from), true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

//------------------------------------------------------------------------------

type Path struct {
	burrowState BurrowState
	energy Energy
}

func NewPath(burrowState BurrowState, energy Energy) Path {
	return Path{burrowState, energy}
}

type BurrowSolution struct {
	bestEnergy map[BurrowState]Energy
	queue []*Path
}

func NewBurrowSolution() BurrowSolution {
	var burrowState BurrowState

	bestEnergy := make(map[BurrowState]Energy)

	path := NewPath(burrowState, 0)
	queue := []*Path{&path}

	return BurrowSolution{bestEnergy, queue}
}

func (this *BurrowSolution) Next() (*Path, bool) {
	// Loop until we find a path that we haven't already seen, or we run out.
	for {
		if len(this.queue) == 0 {
			return nil, false
		}

		path := this.queue[0]

		last := len(this.queue) - 1
		this.queue[0] = this.queue[last]
		this.queue = this.queue[0:last]

		downheap(this, 0)

		if bestEnergy, found := this.bestEnergy[path.burrowState]; found {
			// We've seen this state before. Are we in a better state now?
			if path.energy >= bestEnergy {
				continue
			}
		}

		this.bestEnergy[path.burrowState] = path.energy

		return path, true
	}
}

func (this *BurrowSolution) Add(burrowState BurrowState, prev *Path, energy Energy) {
	//fmt.Printf(" energy %d + %d -> %d\n", prev.energy, energy, prev.energy + energy)
	energy += prev.energy
	if bestEnergy, found := this.bestEnergy[burrowState]; found {
		if energy >= bestEnergy {
			// Already have a same or better path to this state. Ignore it.
			return
		}
	}

	// This is either a new state, or an existing state with a lower energy.
	path := NewPath(burrowState, energy)

	this.queue = append(this.queue, &path)
	upheap(this, len(this.queue)-1)
}

func (this *BurrowSolution) IsHigherPriority(parent, child int) bool {
	return this.queue[parent].energy < this.queue[child].energy
}

func (this *BurrowSolution) IsValid(i int) bool {
	return (i >= 0) && (i < len(this.queue))
}

func (this *BurrowSolution) Swap(i, j int) {
	this.queue[i], this.queue[j] = this.queue[j], this.queue[i]
}

//------------------------------------------------------------------------------

func main() {
	burrow := NewBurrow(aoc.GetFilename(), !false)
	fmt.Println(burrow)

	burrow.Solve()
}

//------------------------------------------------------------------------------

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

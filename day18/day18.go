package main

import (
	"advent-of-code/aoc"
	"fmt"
	"strconv"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

//------------------------------------------------------------------------------

func part1(lines []string) int {
	accum := MakeTree(lines[0])
	accum.Reduce()

	for _, line := range lines[1:] {
		next := MakeTree(line)
		next.Reduce()

		accum = accum.Add(next)
		accum.Reduce()
	}

	return accum.Magnitude()
}

func part2(lines []string) int {
	max := 0

	for i, lhs := range lines {
		for j, rhs := range lines {
			if i == j {
				continue
			}

			a := MakeTree(lhs)
			a.Reduce()

			b := MakeTree(rhs)
			b.Reduce()

			sum := a.Add(b)
			sum.Reduce()

			value := sum.Magnitude()
			if value > max {
				max = value
			}
		}
	}
	return max
}

//------------------------------------------------------------------------------

type NodeType int

const (
	Node NodeType = iota
	Leaf
)

type Tree struct {
	nodeType    NodeType
	value       int
	left, right *Tree
	parent      *Tree
}

func MakeTree(s string) *Tree {

	var parseTree func(int, *Tree) (*Tree, int)
	parseTree = func(offset int, parent *Tree) (*Tree, int) {
		//fmt.Printf("parseTree(%d) %s\n", offset, s[offset:])
		if s[offset] == '[' {
			node := &Tree{nodeType: Node, parent: parent}
			node.left, offset = parseTree(offset+1, node)
			if s[offset] != ',' {
				panic("Expected comma, got: " + s[offset:])
			}
			node.right, offset = parseTree(offset+1, node)
			if s[offset] != ']' {
				panic("Expected ], got: " + s[offset:])
			}
			return node, offset + 1
		}
		value := 0
		for ; isDigit(s[offset]); offset++ {
			value = value*10 + int(s[offset]-'0')
		}
		return &Tree{nodeType: Leaf, value: value, parent: parent}, offset
	}

	tree, rest := parseTree(0, nil)
	if rest < len(s) {
		panic("Garbage after tree parse: " + s[rest:])
	}

	return tree
}

func (t *Tree) Reduce() {
	//fmt.Printf("Reduce: %s\n", t.ToString())
	for i := 1; t.reduceStep(); i++ {
		//fmt.Printf("\nReduce %d: %s\n", i, t.ToString())
	}
	//fmt.Printf("--> Reduced to: %s\n\n", t.ToString())
}

func (t *Tree) reduceStep() bool {

	var explode func(*Tree, int) bool
	explode = func(node *Tree, depth int) bool {
		if node.IsLeaf() {
			return false
		}

		if explode(node.left, depth+1) {
			return true
		}

		if node.shouldExplode(depth) {
			node.explode()
			return true
		}

		return explode(node.right, depth+1)
	}

	var split func(*Tree, int) bool
	split = func(node *Tree, depth int) bool {
		if node.IsLeaf() {
			if node.shouldSplit() {
				node.split()
				return true
			}
			return false
		}

		return split(node.left, depth+1) || split(node.right, depth+1)
	}

	return explode(t, 0) || split(t, 0)
}

func (t *Tree) Magnitude() int {
	if t.IsLeaf() {
		return t.value
	}

	return 3*t.left.Magnitude() + 2*t.right.Magnitude()
}

func (t *Tree) ToString() string {
	if t.nodeType == Leaf {
		return strconv.Itoa(t.value)
	}
	return "[" + t.left.ToString() + "," + t.right.ToString() + "]"
}

func (t *Tree) Add(other *Tree) *Tree {
	root := &Tree{nodeType: Node, left: t, right: other}
	root.left.parent = root
	root.right.parent = root
	return root
}

func (t *Tree) shouldSplit() bool {
	// Assumes Leaf
	return (t.nodeType == Leaf) && (t.value >= 10)
}

func (t *Tree) split() {
	//fmt.Printf("Splitting: %s\n", t.ToString())
	t.nodeType = Node
	t.left = &Tree{nodeType: Leaf, value: t.value / 2, parent: t}
	t.right = &Tree{nodeType: Leaf, value: (t.value + 1) / 2, parent: t}
	t.value = 0
}

func (t *Tree) explode() {
	//fmt.Printf("Exploding: %s\n", t.ToString())
	if !(t.left.IsLeaf() && t.right.IsLeaf()) {
		panic("Trying to explode: " + t.ToString())
	}

	t.explodeLeft(t.left.value)
	t.explodeRight(t.right.value)
	t.nodeType = Leaf
	t.value = 0
	t.left = nil
	t.right = nil
}

func (t *Tree) explodeLeft(value int) {
	//fmt.Printf("ExplodeLeft(%d): %s\n", value, t.ToString())
	for {
		if t.IsRoot() {
			//fmt.Println("-- giving up")
			return
		}
		if t.isRightChild() {
			//fmt.Printf("-- found junction: %s\n", t.parent.ToString())
			t = t.parent.left
			break
		}
		t = t.parent
	}

	for {
		if t.IsNode() {
			//fmt.Printf("-- descending right leg: %s\n", t.right.ToString())
			t = t.right
		} else {
			//fmt.Printf("-- found rightmost: %d + %d => %d\n", t.value, value, t.value + value)
			t.value += value
			return
		}
	}
}

func (t *Tree) explodeRight(value int) {
	//fmt.Printf("ExplodeRight(%d): %s\n", value, t.ToString())
	for {
		if t.IsRoot() {
			//fmt.Println("-- giving up")
			return
		}
		if t.isLeftChild() {
			//fmt.Printf("-- found junction: %s\n", t.parent.ToString())
			t = t.parent.right
			break
		}
		t = t.parent
	}

	for {
		if t.IsNode() {
			//fmt.Printf("-- descending left leg: %s\n", t.left.ToString())
			t = t.left
		} else {
			//fmt.Printf("-- found leftmost: %d + %d => %d\n", t.value, value, t.value + value)
			t.value += value
			return
		}
	}
}

func (t *Tree) isLeftChild() bool {
	return !t.IsRoot() && t.parent.left == t
}

func (t *Tree) isRightChild() bool {
	return !t.IsRoot() && t.parent.right == t
}

func (t *Tree) shouldExplode(depth int) bool {
	// Assumes Node
	return (depth >= 4) && t.left.IsLeaf() && t.right.IsLeaf()
}

func (t *Tree) IsLeaf() bool {
	return t.nodeType == Leaf
}

func (t *Tree) IsNode() bool {
	return t.nodeType == Node
}

func (t *Tree) IsRoot() bool {
	return t.parent == nil
}

func isDigit(c byte) bool {
	return (c >= '0') && (c <= '9')
}

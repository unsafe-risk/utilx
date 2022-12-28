package avltreex

// Package avltree provides associative container that store elements formed by a combination of a key value and a mapped value, following a specific order.
// In a AVLTree, the key values are generally used to sort and uniquely identify the elements, while the mapped values store the content associated to this key.
// The types of key and mapped value may differ.
// Internally, the elements in a AVLTree are always sorted by its key following a specific strict weak ordering criterion
// indicated by its internal comparison object (of type Comparator).
// AVLTree containers are generally slower than go map container to access individual elements by their key,
// but they allow the direct iteration on subsets based on their order.

import (
	"errors"
	"fmt"
	"io"
	"math/bits"

	"golang.org/x/exp/constraints"
)

// Comparator is a function type that would be defined for a key type in the tree.
// This function compares two keys key1, key2 and returns -1 when key1 < key2, 1 when keu1 > key2, and 0 when key1 == key2
// You should pass such function into `NewAVLTree` function.
type Comparator[KeyT any] func(a KeyT, b KeyT) int

// Enumerator is a function type for AVLTree enumeration.
// See Enumerate and EnumerateDiapason for details.
type Enumerator[KeyT any, ValueT any] func(key KeyT, value ValueT) bool

// / Internall stuff
func max[KeyT constraints.Ordered](a KeyT, b KeyT) KeyT {
	if a > b {
		return a
	}
	return b
}

func orderedComparator[KeyT constraints.Ordered](a KeyT, b KeyT) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func orderedComparatorPtr[KeyT constraints.Ordered](a *KeyT, b *KeyT) int {
	if *a == *b {
		return 0
	}
	if *a < *b {
		return -1
	}
	return 1
}

type node[KeyT any, ValueT any] struct {
	key   KeyT
	value ValueT
	links [2]*node[KeyT, ValueT]

	// Balance:
	// -1 - balanced
	//  1 - right higher
	//  0 - left higher
	balance int
}

func getHeight[KeyT any, ValueT any](n *node[KeyT, ValueT]) int {
	if n == nil {
		return 0
	}
	return max(getHeight(n.links[0]), getHeight(n.links[1])) + 1
}

func (n *node[KeyT, ValueT]) getDirection(cmpResult int) int {
	// if cmpResult == -1 {
	// 	return 0
	// }
	// return 1
	return (cmpResult + 1) >> 1
}

func (n *node[KeyT, ValueT]) avlIsBalanced() bool {
	return n.balance < 0
}

type heightChecker func(lh int, rh int)

func recursiveCheckHeight[KeyT any, ValueT any](n *node[KeyT, ValueT], checker heightChecker) {
	if n == nil {
		return
	}
	recursiveCheckHeight(n.links[0], checker)
	recursiveCheckHeight(n.links[1], checker)
	checker(getHeight(n.links[0]), getHeight(n.links[1]))
}

func (t *AVLTree[KeyT, ValueT]) checkHeight(checker heightChecker) {
	recursiveCheckHeight(t.root, checker)
}

func rotate2[KeyT any, ValueT any](pathTop **node[KeyT, ValueT], dir int) *node[KeyT, ValueT] {
	nodeB := *pathTop
	nodeD := nodeB.links[dir]
	nodeC := nodeD.links[1-dir]
	nodeE := nodeD.links[dir]

	*pathTop = nodeD
	nodeD.links[1-dir] = nodeB
	nodeB.links[dir] = nodeC

	return nodeE
}

func rotate3[KeyT any, ValueT any](pathTop **node[KeyT, ValueT], dir int) {
	nodeB := *pathTop
	nodeF := nodeB.links[dir]
	nodeD := nodeF.links[1-dir]
	/* note: C and E can be nil */
	nodeC := nodeD.links[1-dir]
	nodeE := nodeD.links[dir]
	*pathTop = nodeD
	nodeD.links[1-dir] = nodeB
	nodeD.links[dir] = nodeF
	nodeB.links[dir] = nodeC
	nodeF.links[1-dir] = nodeE
}

func avlRotate2[KeyT any, ValueT any](pathTop **node[KeyT, ValueT], dir int) *node[KeyT, ValueT] {
	(*pathTop).balance = -1
	result := rotate2(pathTop, dir)
	(*pathTop).balance = -1
	return result
}

func avlRotate3[KeyT any, ValueT any](pathTop **node[KeyT, ValueT], dir int, third int) *node[KeyT, ValueT] {
	nodeB := *pathTop
	nodeF := nodeB.links[dir]
	nodeD := nodeF.links[1-dir]
	/* note: C and E can be nil */
	nodeC := nodeD.links[1-dir]
	nodeE := nodeD.links[dir]

	nodeB.balance = -1
	nodeF.balance = -1
	nodeD.balance = -1

	rotate3(pathTop, dir)

	if third == -1 {
		return nil
	} else if third == dir {
		/* E holds the insertion so B is unbalanced */
		nodeB.balance = 1 - dir
		return nodeE
	} else {
		/* C holds the insertion so F is unbalanced */
		nodeF.balance = dir
		return nodeC
	}
}

func avlInsert[KeyT any, ValueT any](root **node[KeyT, ValueT], key KeyT, value ValueT, cmp Comparator[KeyT]) bool {
	//Stage 1. Find a position in the tree and link a new node
	// by the way find and remember a node where the tree starts to be unbalanced.
	pathTop := root // Unbalanced node
	nodePtr := root // *nodePtr - a new node

	for nodePtr = root; *nodePtr != nil; {
		n := *nodePtr
		cmpRes := cmp(key, n.key)
		if cmpRes == 0 {
			break
		}
		if !n.avlIsBalanced() {
			pathTop = nodePtr
		}
		dir := n.getDirection(cmpRes)
		nodePtr = &(n.links[dir])
	}

	if *nodePtr != nil {
		return false //already has the key
	}

	*nodePtr = &node[KeyT, ValueT]{
		key:     key,
		value:   value,
		balance: -1,
	}

	//Stage 2. Rebalance
	path := *pathTop
	var first, second, third int
	if !path.avlIsBalanced() {
		first = path.getDirection(cmp(key, path.key))
		if path.balance != first {
			/* took the shorter path */
			path.balance = -1
			path = path.links[first]
		} else {
			second = path.links[first].getDirection(cmp(key, path.links[first].key))
			if first == second {
				/* just a two-point rotate */
				path = avlRotate2(pathTop, first)
			} else {
				/* fine details of the 3 point rotate depend on the third step.
				 * However there may not be a third step, if the third point of the
				 * rotation is the newly inserted point.  In that case we record
				 * the third step as NEITHER
				 */
				path = path.links[first].links[second]
				cmpRes := cmp(key, path.key)
				if cmpRes == 0 {
					third = -1
				} else {
					third = path.getDirection(cmpRes)
				}
				path = avlRotate3(pathTop, first, third)
			}
		}
	}

	//Stage 3. Update balance info in the each node
	for path != nil {
		cmpRes := cmp(key, path.key)
		if cmpRes == 0 {
			break
		}
		direction := path.getDirection(cmpRes)
		path.balance = direction
		path = path.links[direction]
	}
	return true
}

func avlErase[KeyT any, ValueT any](root **node[KeyT, ValueT], key KeyT, cmp Comparator[KeyT]) *node[KeyT, ValueT] {
	//Stage 1. lookup for the node that contain a key
	var targetPtr **node[KeyT, ValueT]
	var dir int
	pathTop := root // Adjust balance start node

	for nodePtr := root; *nodePtr != nil; {
		n := *nodePtr
		cmpRes := cmp(key, n.key)
		dir = n.getDirection(cmpRes)
		if cmpRes == 0 {
			targetPtr = nodePtr
		} else if n.links[dir] == nil {
			break
		} else if n.avlIsBalanced() || (n.balance == (1-dir) && n.links[1-dir].avlIsBalanced()) {
			pathTop = nodePtr
		}
		nodePtr = &n.links[dir]
	}
	if targetPtr == nil {
		return nil //key not found nothing to remove
	}

	/*
	 * Stage 2.
	 * adjust balance, but don't lose 'targetPtr'.
	 * each node from treep down towards target, but
	 * excluding the last, will have a subtree grow
	 * and need rebalancing
	 */
	treep := pathTop
	targetn := *targetPtr
	for {
		tree := *treep
		cmpRes := cmp(key, tree.key)
		bdir := tree.getDirection(cmpRes)
		if tree.links[bdir] == nil {
			break
		} else if tree.avlIsBalanced() {
			tree.balance = 1 - bdir
		} else if tree.balance == bdir {
			tree.balance = -1
		} else {
			second := tree.links[1-bdir].balance
			if second == bdir {
				avlRotate3(treep, 1-bdir, tree.links[1-bdir].links[bdir].balance)
			} else if second == -1 {
				avlRotate2(treep, 1-bdir)
				tree.balance = 1 - bdir
				(*treep).balance = bdir
			} else {
				avlRotate2(treep, 1-bdir)
			}
			if tree == targetn {
				targetPtr = &(*treep).links[bdir]
			}
		}
		treep = &(tree.links[bdir])
	}

	/*
	 * Stage 3.
	 * We have re-balanced everything, it remains only to
	 * swap the end of the path (*treep) with the deleted item
	 * (*targetp)
	 */
	tree := *treep
	targetn = *targetPtr
	*targetPtr = tree
	*treep = tree.links[1-dir]
	tree.links[0] = targetn.links[0]
	tree.links[1] = targetn.links[1]
	tree.balance = targetn.balance

	return targetn
}

func (t *AVLTree[KeyT, ValueT]) findEdgeNodeImpl(key KeyT, dir int) *node[KeyT, ValueT] {
	var n, candidate *node[KeyT, ValueT] = t.root, nil
	for n != nil {
		cmpRes := t.compare(key, n.key)
		if cmpRes == (2*dir - 1) {
			n = n.links[dir]
			continue
		}
		if cmpRes == (1 - 2*dir) {
			candidate = n
			n = n.links[1-dir]
			continue
		}
		if n.links[dir] == nil {
			return candidate
		}
		return edgeNodeImpl(n.links[dir], 1-dir)
	}

	return candidate
}

func edgeNodeImpl[KeyT any, ValueT any](n *node[KeyT, ValueT], dir int) *node[KeyT, ValueT] {
	if n == nil {
		return nil
	}
	for n.links[dir] != nil {
		n = n.links[dir]
	}
	return n
}

func (t *AVLTree[KeyT, ValueT]) lookupNode(key KeyT) *node[KeyT, ValueT] {
	n := t.root
	for n != nil {
		cmp := t.compare(key, n.key)
		if cmp == 0 {
			return n
		} else if cmp == -1 {
			n = n.links[0]
		} else {
			n = n.links[1]
		}
	}
	return n
}

type nodeEnumerator[NodeT any] func(node *NodeT) bool

func (t *AVLTree[KeyT, ValueT]) enumerateNodes(order EnumerationOrder, f nodeEnumerator[node[KeyT, ValueT]]) {
	n := t.root
	if n == nil {
		return
	}

	maxHeight := bits.Len(t.count)
	maxHeight += maxHeight / 2
	stack := make([]*node[KeyT, ValueT], maxHeight)
	stackPtr := 0
	goingDown := true
	for {
		if goingDown {
			//Going down as deep as possible
			for ; n.links[order] != nil; n = n.links[order] {
				stack[stackPtr] = n
				stackPtr++
			}
		}

		// Visit node
		if !f(n) {
			return
		}

		// Going down via second link or return up
		if next := n.links[1-order]; next != nil {
			n = next
			goingDown = true
		} else if stackPtr != 0 {
			stackPtr--
			n = stack[stackPtr]
			goingDown = false
		} else {
			break
		}
	}
}

/// Internall stuff END

// AVLTree is a sorted associative container that contains key-value pairs with unique keys.
// Keys are sorted by using the comparison function `Comparator`.
// Search, removal, and insertion operations have logarithmic complexity.
// KeyT - is a key type. It must be comparable and ordered.
// ValueT - is value type.
// As usual AVLTree instance creation is allowed via `NewAVLTreeOrderedKey` or `NewAVLTree`
type AVLTree[KeyT any, ValueT any] struct {
	root    *node[KeyT, ValueT]
	count   uint
	compare Comparator[KeyT]
}

// NewAVLTree creates a new AVLTree instance with the given Comparator
func NewAVLTree[KeyT any, ValueT any](c Comparator[KeyT]) *AVLTree[KeyT, ValueT] {
	return &AVLTree[KeyT, ValueT]{
		compare: c,
	}
}

// NewAVLTreeOrderedKey creates a new AVLTree instance where Key type is constraints.Ordered.
// This is actually the same as NewAVLTree but Comparator will be defined automaticaly inside the call.
func NewAVLTreeOrderedKey[KeyT constraints.Ordered, ValueT any]() *AVLTree[KeyT, ValueT] {
	return NewAVLTree[KeyT, ValueT](orderedComparator[KeyT])
}

// NewAVLTreeOrderedKeyPtr creates a new AVLTree instance where Key type is a pointer of constraints.Ordered.
// This is actually the same as NewAVLTreeOrderedKey but Comparator will be defined automaticaly inside the call.
func NewAVLTreeOrderedKeyPtr[KeyT constraints.Ordered, ValueT any]() *AVLTree[*KeyT, ValueT] {
	return NewAVLTree[*KeyT, ValueT](orderedComparatorPtr[KeyT])
}

// Size returns the number of elements
func (t *AVLTree[KeyT, ValueT]) Size() uint {
	return t.count
}

// Empty checks whether the container is empty
func (t *AVLTree[KeyT, ValueT]) Empty() bool {
	return t.count == 0
}

// Contains checks if the container contains element with the specific key
func (t *AVLTree[KeyT, ValueT]) Contains(key KeyT) bool {
	return t.lookupNode(key) != nil
}

// Find finds element with specific key
// Returns an pointer on associated with the key value.
// Value modification by the pointer is safe.
// When key isn't present returns nil pointer.
func (t *AVLTree[KeyT, ValueT]) Find(key KeyT) *ValueT {
	n := t.lookupNode(key)
	if n != nil {
		return &n.value
	}
	return nil
}

// FindPrevElement returns a key pointer and a value pointer that is nearest to the given key and lesser then given key.
// Can return (nil, nil) when no such node in the tree.
// Value modification by the pointer is safe.
// Key modification isn't safe!
func (t *AVLTree[KeyT, ValueT]) FindPrevElement(key KeyT) (*KeyT, *ValueT) {
	node := t.findEdgeNodeImpl(key, 0)
	if node != nil {
		return &node.key, &node.value
	}
	return nil, nil
}

// FindNextElement returns a key and a value with the key that is nearest to the given key and greater then given key.
// Can return (nil, nil) when no such node in the tree.
// Value modification by the pointer is safe.
// Key modification isn't safe!
func (t *AVLTree[KeyT, ValueT]) FindNextElement(key KeyT) (*KeyT, *ValueT) {
	node := t.findEdgeNodeImpl(key, 1)
	if node != nil {
		return &node.key, &node.value
	}
	return nil, nil
}

// Insert inserts an element with the given key and value.
// It the given key is already present returns an error.
func (t *AVLTree[KeyT, ValueT]) Insert(key KeyT, value ValueT) error {
	if avlInsert(&t.root, key, value, t.compare) {
		t.count++
		return nil
	}
	return errors.New("AVLTree: already contains key")
}

// EnumerationOrder  a type of enumeration for Enumerate, EnumerateDiapason methods
// There are two acceptable values -  ASCENDING and DESCENDING. All other values provides a runtime error.
// Unfortunately Go doesn't provide any possibiltiy to check wrong values for that in the compile time. So be careful here!
type EnumerationOrder int

const (
	//ASCENDING is an id for Enumerate and EnumerateDiapason methods
	ASCENDING = 0
	//DESCENDING is an id for Enumerate and EnumerateDiapason methods
	DESCENDING = 1
)

// First returns key, value pointers for the first tree node.
// Returns (nil, nil) when a tree is empty.
// Value modification by the pointer is safe.
// Key modification isn't safe!
func (t *AVLTree[KeyT, ValueT]) First() (*KeyT, *ValueT) {
	node := edgeNodeImpl(t.root, ASCENDING)
	if node == nil {
		return nil, nil
	}
	return &node.key, &node.value
}

// Last returns key, value pointers for the last tree node
// Returns (nil, nil) when a tree is empty
// Value modification by the pointer is safe.
// Key modification isn't safe!
func (t *AVLTree[KeyT, ValueT]) Last() (*KeyT, *ValueT) {
	node := edgeNodeImpl(t.root, DESCENDING)
	if node == nil {
		return nil, nil
	}
	return &node.key, &node.value
}

// Erase removes an element by the given key
// Can return an error when such Key wasn't present.
func (t *AVLTree[KeyT, ValueT]) Erase(key KeyT) error {
	if nil != avlErase(&t.root, key, t.compare) {
		t.count--
		return nil
	}
	return errors.New("AVLTree: key not found")
}

// Clear removes all tree content
func (t *AVLTree[KeyT, ValueT]) Clear() {
	t.root = nil
	t.count = 0
}

// Enumerate calls 'Enumerator' for every Tree's element.
// Enumeration order can be one from ASCENDING or DESCENDING
// Enumerator should return `false` for stop enumerating or `true` for continue
func (t *AVLTree[KeyT, ValueT]) Enumerate(order EnumerationOrder, f Enumerator[KeyT, ValueT]) {
	nodeFoo := func(n *node[KeyT, ValueT]) bool {
		return f(n.key, n.value)
	}
	t.enumerateNodes(order, nodeFoo)
}

// EnumerateDiapason works like Enumerate but has two additional args - left and right
// These are left and right borders for enumeration.
// Enumeration includes left and right borders.
// Note: left must be always lesser than right. Otherwise returns error
// Note: left and right should be nil. In means the lesser/greater key in the tree is a border.
//
//	So call EnumerateDiapason where both borders are nil is equivalent to call Enumerate.
//
// Note: If you want to enumerate whole tree call Enumerate since it`s faster!
func (t *AVLTree[KeyT, ValueT]) EnumerateDiapason(left, right *KeyT, order EnumerationOrder, f Enumerator[KeyT, ValueT]) error {
	if t.count == 0 {
		return nil
	}

	if left != nil && right != nil && t.compare(*left, *right) > 0 {
		return errors.New("AVLTree: left must be less rigth")
	}

	//find common sub-tree
	n := t.root
	for {
		if left != nil && t.compare(n.key, *left) < 0 {
			n = n.links[1]
			continue
		}
		if right != nil && t.compare(n.key, *right) > 0 {
			n = n.links[0]
			continue
		}
		break
	}

	fences := [2]*KeyT{left, right}
	maxHeight := bits.Len(t.count)
	maxHeight += maxHeight / 2
	stack := make([]*node[KeyT, ValueT], maxHeight)
	stackPtr := 0
	goingDown := true
loop:
	for {
		if goingDown {
			//Going down as deep as possible
			for {
				if fences[order] != nil && (1-2*int(order))*t.compare(n.key, *fences[order]) < 0 {
					// Try go down via second link
					if next := n.links[1-order]; next != nil && (fences[1-order] == nil || (fences[1-order] != nil && (1-2*int(order))*t.compare(next.key, *fences[1-order]) <= 0)) {
						n = next
						continue
					} else if stackPtr != 0 {
						//Or return up
						stackPtr--
						n = stack[stackPtr]
						goingDown = false
						continue loop
					} else {
						break loop
					}
				} else {
					if n.links[order] == nil {
						break
					}
					stack[stackPtr] = n
					stackPtr++
					n = n.links[order]
				}
			}
		}

		// Visit node
		if !f(n.key, n.value) {
			return nil
		}
		// Going down via second link
		if next := n.links[1-order]; next != nil {
			if fences[1-order] != nil && (1-2*int(order))*t.compare(next.key, *fences[1-order]) >= 0 {
				for ; next != nil; next = next.links[order] {
					if (1-2*int(order))*t.compare(next.key, *fences[1-order]) <= 0 {
						n = next
						goingDown = true
						continue loop
					}
				}
			} else {
				n = next
				goingDown = true
				continue
			}
		}
		if stackPtr != 0 {
			stackPtr--
			n = stack[stackPtr]
			goingDown = false
		} else {
			break
		}
	}
	return nil
}

// BSTDump writes a Tree in graphviz digraph textual format
// See here https://graphviz.org/ for the details
func (t *AVLTree[KeyT, ValueT]) BSTDump(w io.Writer) {
	io.WriteString(w, "digraph BST {\n")
	foo := func(n *node[KeyT, ValueT]) bool {
		if n != nil {
			io.WriteString(w, fmt.Sprintf("\"%v\"-> { ", n.key))
			if n.links[0] != nil {
				io.WriteString(w, fmt.Sprintf("\"%v\" ", n.links[0].key))
			}
			if n.links[1] != nil {
				io.WriteString(w, fmt.Sprintf("\"%v\" ", n.links[1].key))
			}
			io.WriteString(w, "}\n")
		}
		return true
	}
	t.enumerateNodes(ASCENDING, foo)
	io.WriteString(w, "}\n")
}

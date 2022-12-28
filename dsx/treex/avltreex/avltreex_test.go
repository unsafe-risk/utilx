package avltreex

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func createTestTree[KeyT constraints.Ordered](start KeyT, end KeyT, step KeyT) *AVLTree[KeyT, KeyT] {
	tree := NewAVLTreeOrderedKey[KeyT, KeyT]()
	for i := start; i <= end; i += step {
		tree.Insert(i, i)
	}
	return tree
}

func insertKeys[KeyT constraints.Ordered](tree *AVLTree[KeyT, KeyT], keys []KeyT) {
	for _, k := range keys {
		tree.Insert(k, k)
	}
}

func createAndFillTree[KeyT constraints.Ordered](keys []KeyT) *AVLTree[KeyT, interface{}] {
	tree := NewAVLTreeOrderedKey[KeyT, interface{}]()
	for _, k := range keys {
		tree.Insert(k, nil)
	}
	return tree
}

func implEraseTesting[KeyT constraints.Ordered](t *testing.T, keys []KeyT, key KeyT) {
	require := require.New(t)

	tree := createAndFillTree(keys)

	require.True(tree.Contains(key))
	require.Nil(tree.Erase(key))
	tree.checkHeight(func(hl int, hr int) {
		require.LessOrEqual(abs(hl-hr), 1)
	})
	require.False(tree.Contains(key))
	require.Error(tree.Erase(key))
	require.Equal(uint(len(keys)-1), tree.Size())
	for _, k := range keys {
		if k != key {
			require.True(tree.Contains(k))
		}
	}
}

func implInsertTesting(t *testing.T, keys []int) {
	require := require.New(t)

	tree := NewAVLTreeOrderedKey[int, interface{}]()

	size := uint(0)
	for _, key := range keys {
		require.False(tree.Contains(key))
		require.Nil(tree.Insert(key, 0))
		require.Error(tree.Insert(key, 0))
		require.True(tree.Contains(key))
		size++
		require.Equal(size, tree.Size())
		tree.checkHeight(func(hl int, hr int) {
			require.LessOrEqual(abs(hl-hr), 1)
		})
	}
}

func TestCreate(t *testing.T) {
	require := require.New(t)

	tree1 := NewAVLTreeOrderedKey[int, interface{}]()

	require.False(tree1 == nil)
	require.Equal(uint(0), tree1.Size())
	require.True(tree1.Empty())

	tree2 := NewAVLTreeOrderedKeyPtr[int, interface{}]()

	require.False(tree2 == nil)
	require.Equal(uint(0), tree2.Size())
	require.True(tree2.Empty())

}

func TestPointerTree(t *testing.T) {
	require := require.New(t)

	tree := NewAVLTreeOrderedKeyPtr[int, interface{}]()
	key1 := 10
	tree.Insert(&key1, nil)
	key2 := 20
	tree.Insert(&key2, nil)
	key3 := 30
	tree.Insert(&key3, nil)
	require.True(tree.Contains(&key1))
	require.True(tree.Contains(&key2))
	require.True(tree.Contains(&key3))
	key4 := 20
	require.True(tree.Contains(&key4))
}

func TestInsert(t *testing.T) {
	keys1L := []int{1, 2, 3}
	implInsertTesting(t, keys1L)

	keys1R := []int{3, 2, 1}
	implInsertTesting(t, keys1R)

	keys2L := []int{1, 3, 2}
	implInsertTesting(t, keys2L)

	keys2R := []int{3, 1, 2}
	implInsertTesting(t, keys2R)

	keys2a := []int{20, 4, 26, 3, 9, 15}
	implInsertTesting(t, keys2a)

	keys2b := []int{20, 4, 26, 3, 9, 8}
	implInsertTesting(t, keys2b)

	keys3a := []int{20, 4, 26, 3, 9, 21, 30, 2, 7, 11, 15}
	implInsertTesting(t, keys3a)

	keys3b := []int{20, 4, 26, 3, 9, 21, 30, 2, 7, 11, 8}
	implInsertTesting(t, keys3b)
}

func TestErase(t *testing.T) {
	keys1L := []int{3, 4, 5, 2, 6, 1, 7}
	implEraseTesting(t, keys1L, 1)

	keys1R := []int{5, 3, 6, 2, 4, 7, 1}
	implEraseTesting(t, keys1R, 7)

	keys2L := []int{5, 3, 10, 1, 4, 8, 11, 7, 9, 12, 2, 6}
	implEraseTesting(t, keys2L, 2)

	keys2R := []int{8, 3, 11, 2, 5, 9, 12, 1, 4, 6, 10, 7}
	implEraseTesting(t, keys2R, 10)

	keysCase2 := []int{6, 2, 9, 1, 4, 8, 11, 3, 5, 7, 10, 12, 13}
	implEraseTesting(t, keysCase2, 1)

	keysCase3 := []int{3, 1, 4, 2}
	implEraseTesting(t, keysCase3, 3)
}

func TestClear(t *testing.T) {
	require := require.New(t)

	tree := createTestTree(1, 10, 1)

	require.Equal(uint(10), tree.Size())
	require.False(tree.Empty())

	tree.Clear()

	require.Equal(uint(0), tree.Size())
	require.True(tree.Empty())
}

func TestFindExt(t *testing.T) {
	require := require.New(t)

	emptyTree := NewAVLTreeOrderedKey[int, interface{}]()
	require.Nil(emptyTree.First())
	require.Nil(emptyTree.Last())

	const (
		START  = 0
		FINISH = 100
		STEP   = 5
	)

	tree := createTestTree(START, FINISH, STEP)

	require.Nil(tree.Find(START - 1))
	require.Nil(tree.Find(FINISH + 1))
	require.Nil(tree.Find((START - 1) * 2))
	require.Nil(tree.Find((FINISH + 1) * 2))

	for i := START; i <= FINISH; i += STEP {
		require.True(tree.Contains(i))
		require.Equal(*tree.Find(i), i)
	}

	key, _ := tree.First()
	require.Equal(START, *key)

	key, _ = tree.Last()
	require.Equal(FINISH, *key)

	require.Nil(tree.FindPrevElement(START))
	require.Nil(tree.FindNextElement(FINISH))

	for i := START + 1; i <= FINISH; i += STEP {
		k, _ := tree.FindPrevElement(i)
		require.Equal(i-1, *k) //1
		k, _ = tree.FindPrevElement(i + 1)
		require.Equal(i-1, *k) //2
		k, _ = tree.FindPrevElement(i + 2)
		require.Equal(i-1, *k) //3
		k, _ = tree.FindPrevElement(i + 3)
		require.Equal(i-1, *k) //4
		k, _ = tree.FindPrevElement(i + 4)
		require.Equal(i-1, *k) //5
	}

	for i := FINISH - 1; i >= START; i -= STEP {
		k, _ := tree.FindNextElement(i)
		require.Equal(i+1, *k) //99
		k, _ = tree.FindNextElement(i - 1)
		require.Equal(i+1, *k) //98
		k, _ = tree.FindNextElement(i - 2)
		require.Equal(i+1, *k) //97
		k, _ = tree.FindNextElement(i - 3)
		require.Equal(i+1, *k) //96
		k, _ = tree.FindNextElement(i - 4)
		require.Equal(i+1, *k) //95
	}
}

func TestEnumerate(t *testing.T) {
	require := require.New(t)

	emptyTree := NewAVLTreeOrderedKey[int, interface{}]()
	called := false
	emptyTree.Enumerate(ASCENDING, func(k int, v interface{}) bool {
		called = true
		return true
	})
	require.False(called)

	const MIN = -100
	const MAX = 100
	const STEP = 1

	tree := createTestTree(MIN, MAX, STEP)

	i := MIN
	tree.Enumerate(ASCENDING, func(k int, v int) bool {
		require.Equal(i, k)
		i++
		return true
	})

	i = MAX
	tree.Enumerate(DESCENDING, func(k int, v int) bool {
		require.Equal(i, k)
		i--
		return true
	})

	i = MIN
	expectedInterrupt := MIN + 10
	tree.Enumerate(ASCENDING, func(k int, v int) bool {
		require.Equal(i, k)
		if k == expectedInterrupt {
			return false
		}
		i++
		return true
	})
	require.Equal(expectedInterrupt, i)

	i = MAX
	expectedInterrupt = MAX - 10
	tree.Enumerate(DESCENDING, func(k int, v int) bool {
		require.Equal(i, k)
		if k == expectedInterrupt {
			return false
		}
		i--
		return true
	})
	require.Equal(expectedInterrupt, i)
}

func TestEnumerateDiapason(t *testing.T) {
	require := require.New(t)

	emptyTree := NewAVLTreeOrderedKey[int, interface{}]()
	called := false
	emptyTree.EnumerateDiapason(nil, nil, ASCENDING, func(k int, v interface{}) bool {
		called = true
		return true
	})
	require.False(called)

	var (
		START  int = 0
		FINISH int = 100
		STEP   int = 5
	)

	tree := createTestTree(START, FINISH, STEP)

	// Cancel enumeration
	count := 0
	err := tree.EnumerateDiapason(nil, nil, ASCENDING, func(k int, v int) bool {
		count++
		require.Equal(k, START)
		return false
	})
	require.Equal(1, count)
	require.Nil(err)

	// Wrong diapason order
	l := 100
	r := 10
	err = tree.EnumerateDiapason(&l, &r, ASCENDING, func(k int, v int) bool {
		return true
	})
	require.Error(err)

	//0..100
	i := START
	tree.EnumerateDiapason(nil, nil, ASCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(nil, nil, DESCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	//0..100
	i = START
	tree.EnumerateDiapason(&START, &FINISH, ASCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(&START, &FINISH, DESCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	//0..100
	i = START
	left := START - 10
	right := FINISH + 10
	tree.EnumerateDiapason(&left, &right, ASCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	left = START - 10
	right = FINISH + 10
	tree.EnumerateDiapason(&left, &right, DESCENDING, func(k int, v int) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	for i := START; i <= FINISH; i = i + STEP {
		j := i
		tree.EnumerateDiapason(&i, nil, ASCENDING, func(k int, v int) bool {
			require.True(j >= i && j <= FINISH)
			require.Equal(j, k)
			j = j + STEP
			return true
		})

		l := FINISH
		tree.EnumerateDiapason(&l, nil, DESCENDING, func(k int, v int) bool {
			require.True(l >= i && l <= FINISH)
			require.Equal(l, k)
			l = l - STEP
			return true
		})
	}

	for i := START; i <= FINISH; i = i + STEP {
		j := START
		tree.EnumerateDiapason(nil, &i, ASCENDING, func(k int, v int) bool {
			require.True(j >= START && j <= i)
			require.Equal(j, k)
			j = j + STEP
			return true
		})

		l := i
		tree.EnumerateDiapason(nil, &i, DESCENDING, func(k int, v int) bool {
			require.True(l >= START && l <= i)
			require.Equal(l, k)
			l = l - STEP
			return true
		})
	}
}

func TestMaxHeight(t *testing.T) {
	require := require.New(t)

	keys := []int{33,
		20, 46,
		12, 28, 41, 51,
		07, 17, 25, 31, 38, 44, 49, 53,
		04, 10, 15, 19, 23, 27, 30, 32, 36, 40, 43, 45, 48, 50, 52,
		2, 6, 9, 11, 14, 16, 18, 22, 24, 26, 29, 35, 37, 39, 42, 47,
		1, 3, 5, 8, 13, 21, 34,
		0}
	tree := createAndFillTree(keys)

	i := 0
	tree.Enumerate(ASCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		i++
		return true
	})
	require.Equal(i, len(keys))
}

func TestBSTDump(t *testing.T) {
	require := require.New(t)
	const expected = "digraph BST {" +
		"\n\"1\"-> { }" +
		"\n\"2\"-> { \"1\" \"3\" }" +
		"\n\"3\"-> { }" +
		"\n}\n"

	keys := []int{2, 1, 3}
	tree := createAndFillTree(keys)

	builder := new(strings.Builder)
	tree.BSTDump(builder)

	require.Equal(expected, builder.String())
}

package gen

import (
	"runtime/debug"
	"testing"
)

func TestInsertCharacterError(t *testing.T) {
	alphabet := []rune{'a', 'b', 'c'}

	testOutput(
		t,
		"cab",
		InsertCharacterError([]rune("ab"), newGen(t, 2, 0), alphabet),
	)

	testOutput(
		t,
		"abc",
		InsertCharacterError([]rune("ab"), newGen(t, 2, 2), alphabet),
	)

	testOutput(
		t,
		"acb",
		InsertCharacterError([]rune("ab"), newGen(t, 2, 1), alphabet),
	)

	testOutput(
		t,
		"aacbb",
		InsertCharacterError([]rune("aabb"), newGen(t, 2, 2), alphabet),
	)

	testOutput(
		t,
		"a",
		InsertCharacterError([]rune(""), newGen(t, 0, 0), alphabet),
	)
}

func TestRemoveCharacterError(t *testing.T) {
	testOutput(
		t,
		"b",
		RemoveCharacterError([]rune("ab"), newGen(t, 0)),
	)

	testOutput(
		t,
		"ab",
		RemoveCharacterError([]rune("aab"), newGen(t, 0)),
	)

	testOutput(
		t,
		"a",
		RemoveCharacterError([]rune("ab"), newGen(t, 1)),
	)

	testOutput(
		t,
		"",
		RemoveCharacterError([]rune("a"), newGen(t, 0)),
	)

	testOutput(
		t,
		"",
		RemoveCharacterError([]rune(""), newGen(t)),
	)
}

func TestSwapAdjacentCharacterError(t *testing.T) {
	testOutput(
		t,
		"ba",
		SwapAdjacentCharacterError([]rune("ab"), newGen(t, 0)),
	)

	testOutput(
		t,
		"a",
		SwapAdjacentCharacterError([]rune("a"), newGen(t)),
	)

	testOutput(
		t,
		"bac",
		SwapAdjacentCharacterError([]rune("abc"), newGen(t, 0)),
	)
	testOutput(
		t,
		"acb",
		SwapAdjacentCharacterError([]rune("abc"), newGen(t, 1)),
	)
}

func TestReplaceCharacterError(t *testing.T) {
	alphabet := []rune{'a', 'b', 'c'}

	testOutput(
		t,
		"cb",
		ReplaceCharacterError([]rune("ab"), newGen(t, 2, 0), alphabet),
	)

	testOutput(
		t,
		"ac",
		ReplaceCharacterError([]rune("ab"), newGen(t, 2, 1), alphabet),
	)

	testOutput(
		t,
		"c",
		ReplaceCharacterError([]rune("a"), newGen(t, 2, 0), alphabet),
	)

	testOutput(
		t,
		"",
		ReplaceCharacterError([]rune(""), newGen(t), alphabet),
	)
}

type testRandIntGenerator struct {
	t *testing.T
	i int
	r []int
}

func (randGen *testRandIntGenerator) Intn(n int) int {
	out := randGen.r[randGen.i]

	if out >= n {
		randGen.t.Log(string(debug.Stack()))
		randGen.t.Fatal("return value is greater then the maximum allowed return value")
	}

	if out < 0 {
		randGen.t.Log(string(debug.Stack()))
		randGen.t.Fatal("return value should not be smaller then 0")
	}

	randGen.i++
	return out
}

func newGen(t *testing.T, r ...int) RandIntGenerator {
	return &testRandIntGenerator{
		t: t,
		i: 0,
		r: r,
	}
}

func testOutput(t *testing.T, expected string, actual []rune) {
	if expected != string(actual) {
		t.Log(string(debug.Stack()))
		t.Fatalf("expected '%s' but was '%s'", expected, string(actual))
	}
}

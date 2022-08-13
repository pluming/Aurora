package skiplist

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	l := 200000
	sl := MakeSkipList()
	for i := 0; i < 200000; i++ {
		m := fmt.Sprintf("m%d", i)
		sl.Insert(m, float64(i))

	}
	if sl.Length() != uint64(l) {
		t.Fatal("Length unMatch")
	}

	if sl.Level() > maxLevel {
		t.Fatal("Level more than 16")
	}

	deleted := 100
	for i := 0; i < deleted; i++ {
		m := fmt.Sprintf("m%d", i)
		sl.Remove(m, float64(i))
	}
	if sl.Length() != uint64(l)-uint64(deleted) {
		t.Fatal("Length unMatch")
	}

	if sl.Level() > maxLevel {
		t.Fatal("Level more than 16")
	}
	t.Log(sl)
}

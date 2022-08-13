package skiplist

import (
	"fmt"
	"testing"
)

func TestNormalSkip(t *testing.T) {
	l := 20
	sl := MakeSkipList()
	for i := 0; i < l; i++ {
		m := fmt.Sprintf("m%d", i+1)
		sl.Insert(m, float64(i+1))

	}
	if sl.Length() != int64(l) {
		t.Fatal("Length unMatch")
	}

	if sl.Level() > maxLevel {
		t.Fatal("Level more than 16")
	}

	t1 := fmt.Sprintf("m%d", 10)
	rank100 := sl.GetRank(t1, 10)
	if rank100 != 10 {
		t.Fatal("rank UnMatch")
	}

	tr := 6
	rank101Node := sl.GetByRank(int64(tr))
	if rank101Node == nil {
		t.Fatal("rank101Node fail")
	}

	if rank101Node.Score != float64(tr) {
		t.Fatal("rank101Node Score not eq 101")
	}

	deleted := 10
	for i := 0; i < deleted; i++ {
		m := fmt.Sprintf("m%d", i+1)
		sl.Remove(m, float64(i+1))
	}
	if sl.Length() != int64(l)-int64(deleted) {
		t.Fatal("Length unMatch")
	}

	if sl.Level() > maxLevel {
		t.Fatal("Level more than 16")
	}

	t.Log(sl)
}

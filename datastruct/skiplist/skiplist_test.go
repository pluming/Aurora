package skiplist

import (
	"fmt"
	"testing"
)

func TestNormalSkip(t *testing.T) {
	l := 2000
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

	rank100Tmp := 100
	t1 := fmt.Sprintf("m%d", rank100Tmp)
	rank100 := sl.GetRank(t1, float64(rank100Tmp))
	if rank100 != int64(rank100Tmp) {
		t.Fatal("rank UnMatch")
	}

	tr := 600
	rank101Node := sl.GetByRank(int64(tr))
	if rank101Node == nil {
		t.Fatal("rank101Node fail")
	}

	if rank101Node.Score != float64(tr) {
		t.Fatalf("rank101Node Score not eq %d", 600)
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

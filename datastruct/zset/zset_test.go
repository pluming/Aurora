package zset

import (
	"fmt"
	"testing"

	"github.com/pluming/aurora/datastruct/skiplist"
)

func TestNormalZSet(t *testing.T) {
	l := 2000
	zs := New()
	for i := 0; i < l; i++ {
		m := fmt.Sprintf("m%d", i+1)
		zs.Add(m, float64(i+1))
	}
	if zs.Len() != int64(l) {
		t.Fatal("Length unMatch")
	}
	var c float64 = 1
	zs.ForEach(0, zs.Len(), false, func(element *skiplist.Element) bool {
		if c != element.Score {
			t.Fatalf("Cur element score not eq %v", c)
		}
		c++
		return true
	})
	if int64(c-1) != zs.Len() {
		t.Fatalf("expect:%v but:%v", zs.Len(), int64(c))
	}

	elements := zs.Range(0, zs.Len(), false)

	if int64(len(elements)) != zs.Len() {
		t.Fatalf("expect:%v but:%v", zs.Len(), int64(len(elements)))
	}

	element := elements[zs.Len()-1]

	if element.Score != float64(l) {
		t.Fatalf("expect:%v but:%v", float64(l), element.Score)
	}

	t.Log(zs)
}

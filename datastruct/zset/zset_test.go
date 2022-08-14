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

	count := zs.Count(skiplist.NegativeInfBorder, skiplist.PositiveInfBorder)
	t.Log(count)
	if count != zs.Len() {
		t.Fatalf("zs.Count expect:%v but:%v", zs.Len(), count)
	}

	var firstN int64 = 10
	count2 := zs.Count(&skiplist.ScoreBorder{Value: float64(firstN)}, skiplist.PositiveInfBorder)
	t.Log(count2)
	if count2 != int64(zs.Len()-firstN+1) {
		t.Fatalf("zs.Count2 expect:%v but:%v", zs.Len(), count2)
	}

	var cc int64 = firstN
	zs.ForEachByScore(&skiplist.ScoreBorder{Value: float64(firstN)}, skiplist.PositiveInfBorder, 0, 10, false, func(element *skiplist.Element) bool {
		t.Log(element.Member)
		if element.Score != float64(cc) {
			t.Fatalf("zs.ForEachByScore expect:%v but:%v", float64(cc), element.Score)
		}
		cc++
		return true
	})

	var limit int64 = 10
	rangeByScore := zs.RangeByScore(&skiplist.ScoreBorder{Value: float64(firstN)}, skiplist.PositiveInfBorder, 0, limit, false)

	if int64(len(rangeByScore)) != limit {
		t.Fatalf("zs.RangeByScore expect:%v but:%v", int64(len(rangeByScore)), limit)
	}

	//delCount := zs.RemoveByScore(&skiplist.ScoreBorder{Value: float64(1)}, &skiplist.ScoreBorder{Value: float64(firstN)})
	delCount := zs.RemoveByRank(0, 10)
	t.Log(delCount)
	t.Log(zs)
}

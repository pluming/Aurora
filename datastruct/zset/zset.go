package zset

import (
	"strconv"

	"github.com/pluming/aurora/datastruct/skiplist"
)

const (
	initZSetDictSize = 16
)

type ZSet struct {
	//dict dict.Dict
	dict map[string]*skiplist.Element

	skipList *skiplist.SkipList
}

func New() *ZSet {
	return &ZSet{
		dict:     make(map[string]*skiplist.Element, initZSetDictSize),
		skipList: skiplist.MakeSkipList(),
	}
}

// Add puts member into set,  and returns whether has inserted new node
func (zs *ZSet) Add(member string, score float64) bool {
	element, ok := zs.dict[member]

	//exist and score eq ==> not modify
	if ok {
		if element.Score != score {
			zs.dict[member] = &skiplist.Element{
				Member: member,
				Score:  score,
			}
			zs.skipList.Remove(member, score)
			zs.skipList.Insert(member, score)
		}
		return false
	}
	zs.dict[member] = &skiplist.Element{
		Member: member,
		Score:  score,
	}
	zs.skipList.Insert(member, score)
	return true
}

func (zs *ZSet) Len() int64 {
	return int64(len(zs.dict))
}

func (zs *ZSet) Get(member string) (*skiplist.Element, bool) {
	element, ok := zs.dict[member]
	if !ok {
		return nil, false
	}
	return element, ok
}

func (zs *ZSet) Remove(member string) bool {
	element, ok := zs.dict[member]
	if ok {
		delete(zs.dict, member)
		zs.skipList.Remove(element.Member, element.Score)
		return true
	}
	return false
}

// ForEach visits each member which rank within [start, stop), sort by ascending order, rank starts from 0
func (zs *ZSet) ForEach(start, stop int64, desc bool, cb func(element *skiplist.Element) bool) {
	size := zs.Len()
	if start < 0 || start >= size {
		panic("illegal start " + strconv.FormatInt(start, 10))
	}

	if stop < start || stop > size {
		panic("illegal stop " + strconv.FormatInt(start, 10))
	}

	var node *skiplist.Node
	if desc {
		node = zs.skipList.Tail
		if start > 0 {
			node = zs.skipList.GetByRank(size - start)
		}
	} else {
		node = zs.skipList.Header.Level[0].Forward
		if start > 0 {
			node = zs.skipList.GetByRank(start + 1)
		}
	}

	fi := stop - start
	for i := int64(0); i < fi; i++ {
		if !cb(&node.Element) {
			break
		}
		if desc {
			node = node.Backward
		} else {
			node = node.Level[0].Forward
		}
	}
}

// Range returns members which rank within [start, stop), sort by ascending order, rank starts from 0
func (zs *ZSet) Range(start int64, stop int64, desc bool) []*skiplist.Element {
	sSize := stop - start
	slice := make([]*skiplist.Element, 0, sSize)
	zs.ForEach(start, stop, desc, func(element *skiplist.Element) bool {
		slice = append(slice, element)
		return true
	})
	return slice
}

func (zs *ZSet) Count(min *skiplist.ScoreBorder, max *skiplist.ScoreBorder) int64 {
	var c int64
	_, firstRank := zs.skipList.GetFirstInScoreRange(min, max)
	if firstRank < 0 {
		return 0
	}
	_, lastRank := zs.skipList.GetLastInScoreRange(min, max)
	if lastRank < 0 {
		return 0
	}
	zs.ForEach(firstRank-1, lastRank, false, func(element *skiplist.Element) bool {
		c++
		return true
	})
	return c
}

// ForEachByScore offset : 0-base
// ForEachByScore visits members which score within the given border
func (zs *ZSet) ForEachByScore(min, max *skiplist.ScoreBorder, offset, limit int64, desc bool, cb func(element *skiplist.Element) bool) {
	if offset < 0 {
		return
	}
	_, firstRank := zs.skipList.GetFirstInScoreRange(min, max)
	if firstRank < 0 {
		return
	}
	_, lastRank := zs.skipList.GetLastInScoreRange(min, max)
	if lastRank < 0 {
		return
	}

	start := firstRank + offset - 1
	stop := lastRank
	if start < 0 || start >= zs.Len() {
		return
	}
	if stop < start || stop > zs.Len() {
		return
	}
	var c int64
	zs.ForEach(start, stop, desc, func(element *skiplist.Element) bool {
		if limit < 0 || c >= limit {
			return false
		}
		c++
		if !cb(element) {
			return false
		}
		if c == limit {
			return false
		}
		return true
	})
}

// RangeByScore returns members which score within the given border
// param limit: <0 means no limit
func (zs *ZSet) RangeByScore(min *skiplist.ScoreBorder, max *skiplist.ScoreBorder, offset int64, limit int64, desc bool) []*skiplist.Element {
	if limit == 0 || offset < 0 {
		return make([]*skiplist.Element, 0)
	}
	slice := make([]*skiplist.Element, 0)
	zs.ForEachByScore(min, max, offset, limit, desc, func(element *skiplist.Element) bool {
		slice = append(slice, element)
		return true
	})
	return slice
}

// RemoveByScore removes members which score within the given border
func (zs *ZSet) RemoveByScore(min, max *skiplist.ScoreBorder) int64 {
	removed := zs.skipList.RemoveRangeByScore(min, max)
	for _, element := range removed {
		delete(zs.dict, element.Member)
	}
	return int64(len(removed))
}

// RemoveByRank removes member ranking within [start, stop)
// sort by ascending order and rank starts from 0
func (zs *ZSet) RemoveByRank(start, stop int64) int64 {
	removed := zs.skipList.RemoveRangeByRank(start+1, stop+1)
	for _, element := range removed {
		delete(zs.dict, element.Member)
	}
	return int64(len(removed))
}

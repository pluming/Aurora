package skiplist

import (
	"math/rand"
	"time"
)

const (
	maxLevel = 16
)

// Element key-score pair
type Element struct {
	Member string
	Score  float64
}

type Level struct {
	Forward *Node
	span    int64
}

type Node struct {
	Element
	Backward *Node
	Level    []*Level
}

type SkipList struct {
	Header *Node

	Tail *Node

	level int16

	length int64
}

func (sl *SkipList) Level() int16 {
	return sl.level
}

func (sl *SkipList) Length() int64 {
	return sl.length
}

func MakeNode(member string, score float64, level int16) *Node {
	n := &Node{
		Element: Element{
			Member: member,
			Score:  score,
		},
		Level: make([]*Level, level),
	}
	for i := 0; i < int(level); i++ {
		n.Level[i] = &Level{}
	}
	return n
}

func MakeSkipList() *SkipList {
	sl := &SkipList{
		Header: MakeNode("", 0, maxLevel),
		level:  1,
	}
	return sl
}

func init() {
	rand.Seed(time.Now().Unix())
}

func randomLevel() int16 {
	var level int16 = 1

	for float64(rand.Int31()&0xffff) < (0.25 * 0xffff) {
		level++
	}
	if level > maxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(member string, score float64) *Node {
	update := make([]*Node, maxLevel)
	rank := make([]int64, maxLevel)

	node := sl.Header
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] //starting store last rank
		}

		if node.Level[i] != nil {
			//Level traverse the SkipList
			for node.Level[i].Forward != nil &&
				(node.Level[i].Forward.Score < score ||
					(node.Level[i].Forward.Score == score && node.Level[i].Forward.Member < member)) {
				rank[i] += node.Level[i].span
				node = node.Level[i].Forward
			}
		}
		update[i] = node
	}
	level := randomLevel()
	if level > sl.level {
		//extend sk Level
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.Header
			update[i].Level[i].span = sl.length
		}
		sl.level = level
	}
	node = MakeNode(member, score, level)
	for i := int16(0); i < level; i++ {
		node.Level[i].Forward = update[i].Level[i].Forward
		update[i].Level[i].Forward = node
		node.Level[i].span = update[i].Level[i].span - (rank[0] - rank[i])
		update[i].Level[i].span = (rank[0] - rank[i]) + 1
		if rank[0]-rank[i] < 0 {
			panic("fsadfasdfas")
		}
	}
	//inc +1  untouched Level
	for i := level; i < sl.level; i++ {
		update[i].Level[i].span++
	}
	//update node bw
	if update[0] == sl.Header {
		node.Backward = nil
	} else {
		node.Backward = update[0]
	}
	//update node.Forward
	if node.Level[0].Forward != nil {
		node.Level[0].Forward.Backward = node
	} else {
		sl.Tail = node
	}
	sl.length++
	return nil
}

// Remove deleted and return true; otherwise false
func (sl *SkipList) Remove(member string, score float64) bool {
	update := make([]*Node, maxLevel)
	node := sl.Header
	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].Forward != nil &&
			(node.Level[i].Forward.Score < score ||
				(node.Level[i].Forward.Score == score && node.Level[i].Forward.Member < member)) {
			node = node.Level[i].Forward
		}
		update[i] = node
	}

	node = node.Level[0].Forward
	if node != nil && node.Score == score && node.Member == member {
		sl.removeNode(node, update)
		return true
	}

	return false
}

/*
n: target node to delete
update: Backward node of target
*/
func (sl *SkipList) removeNode(n *Node, update []*Node) {
	//update.Level.Forward=n.Level.Forward
	for i := int16(0); i < sl.level; i++ {
		if update[i].Level[i].Forward == n {
			update[i].Level[i].Forward = n.Level[i].Forward
			update[i].Level[i].span += n.Level[i].span - 1
		} else {
			update[i].Level[i].span--
		}
	}

	// not last
	if n.Level[0].Forward != nil {
		n.Level[0].Forward.Backward = n.Backward
	} else { //last
		sl.Tail = n.Backward
	}
	//desc -1 Level
	for sl.level > 1 && sl.Header.Level[sl.level-1].Forward == nil {
		sl.level--
	}
	sl.length--
}

// GetRank 0 not found ;rank return
func (sl *SkipList) GetRank(member string, score float64) int64 {
	var rank int64
	node := sl.Header
	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].Forward != nil &&
			(node.Level[i].Forward.Score < score ||
				(node.Level[i].Forward.Score == score && node.Level[i].Forward.Member <= member)) {
			rank += node.Level[i].span
			node = node.Level[i].Forward
		}
		if node.Member == member {
			return rank
		}
	}
	return 0
}

// GetByRank 1-based rank
// GetByRank nil not found
func (sl *SkipList) GetByRank(rank int64) *Node {
	// rank less 0 return
	if rank <= 0 {
		return nil
	}
	var r int64
	node := sl.Header
	// traverse sl; scan from top Level
	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].Forward != nil && (r+node.Level[i].span) <= rank {
			r += node.Level[i].span
			node = node.Level[i].Forward
		}
		if r == rank {
			return node
		}
	}
	return nil
}

// GetFirstInScoreRange 获取指定范围的Node ,返回rank; rank 1 base;-1 UnMatch
func (sl *SkipList) GetFirstInScoreRange(min *ScoreBorder, max *ScoreBorder) (*Node, int64) {
	var rank int64
	if !sl.hasInRange(min, max) {
		return nil, -1
	}
	node := sl.Header
	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].Forward != nil && !min.less(node.Level[i].Forward.Score) {
			rank += node.Level[i].span
			node = node.Level[i].Forward
		}
	}

	node = node.Level[0].Forward
	if node == nil {
		return nil, -1
	}
	if !max.greater(node.Score) {
		return nil, -1
	}
	rank++
	return node, rank
}

func (sl *SkipList) GetLastInScoreRange(min *ScoreBorder, max *ScoreBorder) (*Node, int64) {
	var rank int64
	if !sl.hasInRange(min, max) {
		return nil, -1
	}
	node := sl.Header
	// scan from top Level
	for level := sl.level - 1; level >= 0; level-- {
		for node.Level[level].Forward != nil && max.greater(node.Level[level].Forward.Score) {
			rank += node.Level[level].span
			node = node.Level[level].Forward
		}
	}

	if node == nil {
		return nil, -1
	}
	if !min.less(node.Score) {
		return nil, -1
	}
	return node, rank
}

func (sl *SkipList) hasInRange(min *ScoreBorder, max *ScoreBorder) bool {
	if min.Value > max.Value || (min.Value == max.Value && (min.Exclude || max.Exclude)) {
		return false
	}
	node := sl.Tail
	//空表；
	//min > Tail
	if node == nil || !min.less(node.Score) {
		return false
	}
	//max < min
	node = sl.Header.Level[0].Forward
	if node == nil || !max.greater(node.Score) {
		return false
	}
	return true
}

// RemoveRangeByScore return removed elements
func (sl *SkipList) RemoveRangeByScore(min *ScoreBorder, max *ScoreBorder) (removed []*Element) {
	update := make([]*Node, maxLevel)
	removed = make([]*Element, 0)
	// find Backward nodes (of target range) or last node of each Level
	node := sl.Header
	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].Forward != nil {
			if min.less(node.Level[i].Forward.Score) { // already in range
				break
			}
			node = node.Level[i].Forward
		}
		update[i] = node
	}

	// node is the first one within range
	node = node.Level[0].Forward

	// remove nodes in range
	for node != nil {
		if !max.greater(node.Score) { // already out of range
			break
		}
		next := node.Level[0].Forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		sl.removeNode(node, update)
		node = next
	}
	return removed
}

// RemoveRangeByRank 1-based rank, including start, exclude stop
func (sl *SkipList) RemoveRangeByRank(start int64, stop int64) (removed []*Element) {
	var i int64 = 0 // rank of iterator
	update := make([]*Node, maxLevel)
	removed = make([]*Element, 0)

	// scan from top Level
	node := sl.Header
	for level := sl.level - 1; level >= 0; level-- {
		for node.Level[level].Forward != nil && (i+node.Level[level].span) < start {
			i += node.Level[level].span
			node = node.Level[level].Forward
		}
		update[level] = node
	}

	i++
	node = node.Level[0].Forward // first node in range

	// remove nodes in range
	for node != nil && i < stop {
		next := node.Level[0].Forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		sl.removeNode(node, update)
		node = next
		i++
	}
	return removed
}

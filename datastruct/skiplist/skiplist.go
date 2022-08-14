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
	forward *Node
	span    int64
}

type Node struct {
	Element
	backward *Node
	level    []*Level
}

type skipList struct {
	header *Node

	tail *Node

	level int16

	length int64
}

func (sl *skipList) Level() int16 {
	return sl.level
}

func (sl *skipList) Length() int64 {
	return sl.length
}

func MakeNode(member string, score float64, level int16) *Node {
	n := &Node{
		Element: Element{
			Member: member,
			Score:  score,
		},
		level: make([]*Level, level),
	}
	for i := 0; i < int(level); i++ {
		n.level[i] = &Level{}
	}
	return n
}

func MakeSkipList() *skipList {
	sl := &skipList{
		header: MakeNode("", 0, maxLevel),
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

func (sl *skipList) Insert(member string, score float64) *Node {
	update := make([]*Node, maxLevel)
	rank := make([]int64, maxLevel)

	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] //starting store last rank
		}

		if node.level[i] != nil {
			//level traverse the skipList
			for node.level[i].forward != nil &&
				(node.level[i].forward.Score < score ||
					(node.level[i].forward.Score == score && node.level[i].forward.Member < member)) {
				rank[i] += node.level[i].span
				node = node.level[i].forward
			}
		}
		update[i] = node
	}
	level := randomLevel()
	if level > sl.level {
		//extend sk level
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = sl.length
		}
		sl.level = level
	}
	node = MakeNode(member, score, level)
	for i := int16(0); i < level; i++ {
		node.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = node
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
		if rank[0]-rank[i] < 0 {
			panic("fsadfasdfas")
		}
	}
	//inc +1  untouched level
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}
	//update node bw
	if update[0] == sl.header {
		node.backward = nil
	} else {
		node.backward = update[0]
	}
	//update node.forward
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node
	} else {
		sl.tail = node
	}
	sl.length++
	return nil
}

// Remove deleted and return true; otherwise false
func (sl *skipList) Remove(member string, score float64) bool {
	update := make([]*Node, maxLevel)
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			(node.level[i].forward.Score < score ||
				(node.level[i].forward.Score == score && node.level[i].forward.Member < member)) {
			node = node.level[i].forward
		}
		update[i] = node
	}

	node = node.level[0].forward
	if node != nil && node.Score == score && node.Member == member {
		sl.removeNode(node, update)
		return true
	}

	return false
}

/*
n: target node to delete
update: backward node of target
*/
func (sl *skipList) removeNode(n *Node, update []*Node) {
	//update.level.forward=n.level.forward
	for i := int16(0); i < sl.level; i++ {
		if update[i].level[i].forward == n {
			update[i].level[i].forward = n.level[i].forward
			update[i].level[i].span += n.level[i].span - 1
		} else {
			update[i].level[i].span--
		}
	}

	// not last
	if n.level[0].forward != nil {
		n.level[0].forward.backward = n.backward
	} else { //last
		sl.tail = n.backward
	}
	//desc -1 level
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

// GetRank 0 not found ;rank return
func (sl *skipList) GetRank(member string, score float64) int64 {
	var rank int64
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			(node.level[i].forward.Score < score ||
				(node.level[i].forward.Score == score && node.level[i].forward.Member <= member)) {
			rank += node.level[i].span
			node = node.level[i].forward
		}
		if node.Member == member {
			return rank
		}
	}
	return 0
}

// GetByRank 1-based rank
// GetByRank nil not found
func (sl *skipList) GetByRank(rank int64) *Node {
	// rank less 0 return
	if rank <= 0 {
		return nil
	}
	var r int64
	node := sl.header
	// traverse sl; scan from top level
	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil && (r+node.level[i].span) <= rank {
			r += node.level[i].span
			node = node.level[i].forward
		}
		if r == rank {
			return node
		}
	}
	return nil
}

func (sl *skipList) GetFirstInScoreRange(min *ScoreBorder, max *ScoreBorder) *Node {
	if !sl.hasInRange(min, max) {
		return nil
	}
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil && !min.less(node.level[i].forward.Score) {
			node = node.level[i].forward
		}
	}

	node = node.level[0].forward
	if node == nil {
		return nil
	}
	if !max.greater(node.Score) {
		return nil
	}
	return node
}

func (sl *skipList) getLastInScoreRange(min *ScoreBorder, max *ScoreBorder) *Node {
	if !sl.hasInRange(min, max) {
		return nil
	}
	node := sl.header
	// scan from top level
	for level := sl.level - 1; level >= 0; level-- {
		for node.level[level].forward != nil && max.greater(node.level[level].forward.Score) {
			node = node.level[level].forward
		}
	}

	if node == nil {
		return nil
	}
	if !min.less(node.Score) {
		return nil
	}
	return node
}

func (sl *skipList) hasInRange(min *ScoreBorder, max *ScoreBorder) bool {
	if min.Value > max.Value || (min.Value == max.Value && (min.Exclude || max.Exclude)) {
		return false
	}
	node := sl.tail
	//空表；
	//min > tail
	if node == nil || !min.less(node.Score) {
		return false
	}
	//max < min
	node = sl.header.level[0].forward
	if node == nil || !max.greater(node.Score) {
		return false
	}
	return true
}

// RemoveRangeByScore return removed elements
func (sl *skipList) RemoveRangeByScore(min *ScoreBorder, max *ScoreBorder) (removed []*Element) {
	update := make([]*Node, maxLevel)
	removed = make([]*Element, 0)
	// find backward nodes (of target range) or last node of each level
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil {
			if min.less(node.level[i].forward.Score) { // already in range
				break
			}
			node = node.level[i].forward
		}
		update[i] = node
	}

	// node is the first one within range
	node = node.level[0].forward

	// remove nodes in range
	for node != nil {
		if !max.greater(node.Score) { // already out of range
			break
		}
		next := node.level[0].forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		sl.removeNode(node, update)
		node = next
	}
	return removed
}

// RemoveRangeByRank 1-based rank, including start, exclude stop
func (sl *skipList) RemoveRangeByRank(start int64, stop int64) (removed []*Element) {
	var i int64 = 0 // rank of iterator
	update := make([]*Node, maxLevel)
	removed = make([]*Element, 0)

	// scan from top level
	node := sl.header
	for level := sl.level - 1; level >= 0; level-- {
		for node.level[level].forward != nil && (i+node.level[level].span) < start {
			i += node.level[level].span
			node = node.level[level].forward
		}
		update[level] = node
	}

	i++
	node = node.level[0].forward // first node in range

	// remove nodes in range
	for node != nil && i < stop {
		next := node.level[0].forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		sl.removeNode(node, update)
		node = next
		i++
	}
	return removed
}

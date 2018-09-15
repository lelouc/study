package skipList

import "math/rand"

var SKIP_LIST_P float64 = 0.25
var SKIP_LIST_MAXLEVEL uint32 = 32

type Comparator interface {
	CmpKey(interface{}, interface{}) int
	CmpScore(interface{}, interface{}) int
}

type SkipList struct {
	Comparator
	level  uint32
	length uint32
	head   *SkipListNode
	tail   *SkipListNode
}

func NewSkipList(cmp Comparator) *SkipList {
	head := NewSkipListNode(SKIP_LIST_MAXLEVEL, nil)
	for i := 0; i < int(SKIP_LIST_MAXLEVEL); i++ {
		head.SetForward(i, nil)
		head.SetSpan(i, 0)
	}
	head.SetBackward(nil)

	return &SkipList{
		Comparator: cmp,
		level:      1,
		length:     0,
		head:       head,
		tail:       nil,
	}
}

func (this *SkipList) Level() uint32        { return this.level }
func (this *SkipList) Length() uint32       { return this.length }
func (this *SkipList) Head() *SkipListNode  { return this.head }
func (this *SkipList) Tail() *SkipListNode  { return this.tail }
func (this *SkipList) First() *SkipListNode { return this.head.Forward(0) }

func (this *SkipList) Less(value1, value2 interface{}) bool {
	return this.CmpScore(value1, value2) < 0 ||
		(this.CmpScore(value1, value2) == 0 && this.CmpKey(value1, value2) < 0)
}

func (this *SkipList) Equal(value1, value2 interface{}) bool {
	return this.CmpScore(value1, value2) == 0 && this.CmpKey(value1, value2) == 0
}

func (this *SkipList) LessEqual(value1, value2 interface{}) bool {
	return this.CmpScore(value1, value2) < 0 || (this.CmpScore(value1, value2) == 0 && this.CmpKey(value1, value2) <= 0)
}

func (this *SkipList) RandomLevel() uint32 {
	level := uint32(1)
	for (rand.Uint32()&0xFFFF) < uint32(SKIP_LIST_P*0xFFFF) && level < SKIP_LIST_MAXLEVEL {
		level++
	}
	return level
}

func (this *SkipList) Insert(value interface{}) *SkipListNode {
	rank := make([]uint32, SKIP_LIST_MAXLEVEL)
	update := make([]*SkipListNode, SKIP_LIST_MAXLEVEL)

	curNode := this.head
	for i := int(this.level - 1); i >= 0; i-- {
		if i != int(this.length-1) {
			rank[i] = rank[i+1]
		}
		for nextNode := curNode.Forward(i); nextNode != nil && this.Less(nextNode.value, value); nextNode = nextNode.Forward(i) {
			rank[i] += nextNode.Span(i)
			curNode = nextNode
		}
		update[i] = curNode
	}
	return this.InsertNode(value, rank, update)
}

func (this *SkipList) InsertNode(value interface{}, rank []uint32, update []*SkipListNode) *SkipListNode {
	level := this.RandomLevel()
	if level > this.level {
		for i := int(this.level); i < int(level); i++ {
			rank[i] = 0
			update[i] = this.head
			update[i].SetSpan(i, this.length)
		}
		this.level = level
	}

	newNode := NewSkipListNode(level, value)
	for i := 0; i < int(level); i++ {
		newNode.SetForward(i, update[i].Forward(i))
		update[i].SetForward(i, newNode)

		newNode.SetSpan(i, update[i].Span(i)-(rank[0]-rank[i]))
		update[i].SetSpan(i, rank[0]-rank[i]+1)
	}

	for i := int(level); i < int(this.level); i++ {
		update[i].SetSpan(i, update[i].Span(i)+1)
	}

	if update[0] != this.head {
		newNode.SetBackward(update[0])
	} else {
		newNode.SetBackward(nil)
	}

	if nextNode := newNode.Forward(0); nextNode != nil {
		nextNode.SetBackward(newNode)
	} else {
		this.tail = newNode
	}

	this.length++
	return newNode
}

func (this *SkipList) Delete(value interface{}) bool {
	update := make([]*SkipListNode, SKIP_LIST_MAXLEVEL)

	curNode := this.head
	for i := int(this.level - 1); i >= 0; i-- {
		for nextNode := curNode.Forward(i); nextNode != nil && this.Less(nextNode.value, value); nextNode = nextNode.Forward(i) {
			curNode = nextNode
		}
		update[i] = curNode
	}

	if delNode := curNode.Forward(0); delNode != nil &&
		this.Equal(delNode.value, value) {
		this.DeleteNode(delNode, update)
		return true
	}
	return false
}

func (this *SkipList) DeleteNode(delNode *SkipListNode, update []*SkipListNode) {
	for i := int(this.level - 1); i >= 0; i-- {
		if update[i].Forward(i) == delNode {
			update[i].SetForward(i, delNode.Forward(i))
			update[i].SetSpan(i, update[i].Span(i)+delNode.Span(i)-1)
		} else {
			update[i].SetSpan(i, update[i].Span(i)-1)
		}
	}

	if delNode.Forward(0) != nil {
		delNode.Forward(0).SetBackward(delNode.Backward())
	} else {
		this.tail = delNode.Backward()
	}

	for this.level > 1 && this.head.Forward(int(this.level-1)) == nil {
		this.level--
	}
	this.length--
}

func (this *SkipList) GetRank(value interface{}) uint32 {
	var rank uint32
	curNode := this.head
	for i := int(this.level - 1); i >= 0; i-- {
		for nextNode := curNode.Forward(i); nextNode != nil && this.LessEqual(nextNode.value, value); nextNode = nextNode.Forward(i) {
			rank += curNode.Span(i)
			curNode = nextNode
		}
		if curNode != this.head && this.Equal(curNode.value, value) {
			return rank
		}
	}
	return 0
}

func (this *SkipList) GetNodeByRank(rank uint32) *SkipListNode {
	if rank > 0 {
		var traversed uint32
		var curNode *SkipListNode
		for i := int(this.level - 1); i >= 0; i-- {
			traversed = 0
			curNode = this.head
			for nextNode := curNode.Forward(i); nextNode != nil && traversed+curNode.Span(i) <= rank; nextNode = nextNode.Forward(i) {
				traversed += curNode.Span(i)
				curNode = nextNode
			}

			if traversed == rank {
				return curNode
			}
		}
	}
	return nil
}

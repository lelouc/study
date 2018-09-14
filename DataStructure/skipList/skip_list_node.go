package skipList

type SkipListLevel struct {
	span    uint32
	forward *SkipListNode
}

type SkipListNode struct {
	value    interface{}
	backward *SkipListNode
	level    []SkipListLevel
}

func NewSkipListNode(level uint32, value interface{}) *SkipListNode {
	return &SkipListNode{
		value: value,
		level: make([]SkipListLevel, level),
	}
}

func (this *SkipListNode) Value() interface{} {
	return this.Value
}

func (this *SkipListNode) Next() *SkipListNode {
	return this.level[0].forward
}

func (this *SkipListNode) Prev() *SkipListNode {
	return this.backward
}

func (this *SkipListNode) Span(i int) uint32 {
	return this.level[i].span
}

func (this *SkipListNode) SetSpan(i int, span uint32) {
	this.level[i].span = span
}

func (this *SkipListNode) Forward(i int) *SkipListNode {
	return this.level[i].forward
}

func (this *SkipListNode) SetForward(i int, node *SkipListNode) {
	this.level[i].forward = node
}

func (this *SkipListNode) Backward() *SkipListNode {
	return this.backward
}

func (this *SkipListNode) SetBackward(node *SkipListNode) {
	this.backward = node
}

package compute

type Stacker interface {
	Push(interface{})
	Pop() interface{}
	Size() int
	IsEmpty() bool
	Peek() interface{}
}

type Stack struct {
	stack []interface{}
}

func NewStack() Stacker {
	return &Stack{}
}

func (cs *Stack) Push(element interface{}) {
	cs.stack = append(cs.stack, element)
}

func (cs *Stack) Pop() interface{} {
	element := cs.stack[len(cs.stack)-1]
	cs.stack = cs.stack[:len(cs.stack)-1]
	return element
}

func (cs *Stack) Size() int {
	return len(cs.stack)
}

func (cs *Stack) IsEmpty() bool {
	return len(cs.stack) == 0
}

func (cs *Stack) Peek() interface{} {
	return cs.stack[len(cs.stack)-1]
}

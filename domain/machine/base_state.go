package machine

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine/ability"
)

type BaseState struct {
	Nexts map[byte]sets.Set[ability.State]
}

func NewBaseState() *BaseState {
	return &BaseState{
		Nexts: make(map[byte]sets.Set[ability.State]),
	}
}

// AddState 新增节点的转移状态
func (node *BaseState) AddState(c byte, state ability.State) {
	if _, ok := node.Nexts[c]; !ok {
		node.Nexts[c] = sets.NewSet[ability.State]()
	}
	node.Nexts[c].Add(state)
}

func (node *BaseState) RemoveState(c byte, state ability.State) {
	node.Nexts[c].Remove(state)
}

func (node *BaseState) GetNexts() map[byte]sets.Set[ability.State] {
	return node.Nexts
}

// Process 根据字符获取该节点的所有转移状态
func (node *BaseState) Process(c byte) []ability.State {
	return node.Nexts[c].ToSlices()
}

package common

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type CommonState struct {
	*machine.BaseState
}

// 创建新的状态节点
func newCommonState() *CommonState {
	return &CommonState{
		BaseState: machine.NewBaseState(),
	}
}

func (node *CommonState) GetNexts() map[byte]sets.Set[ability.State] {
	return node.Nexts
}

func (node *CommonState) Process(c byte) []ability.State {
	states := node.BaseState.Process(c)
	if extraStates, ok := node.Nexts['.']; ok {
		states = append(states, extraStates.ToSlices()...)
	}
	return states
}

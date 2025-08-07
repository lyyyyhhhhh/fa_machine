package common

import (
	"fa_machine/common/sets"
	"fa_machine/machine2/domain"
	"fa_machine/machine2/domain/ability"
)

type CommonState struct {
	*domain.BaseState
}

// 创建新的状态节点
func newCommonState() *CommonState {
	return &CommonState{
		BaseState: domain.NewBaseState(),
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

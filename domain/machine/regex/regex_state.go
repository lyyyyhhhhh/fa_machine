package regex

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type RegexState struct {
	*machine.BaseState
	pre ability.State
}

// 创建新的状态节点
func newRegexState() *RegexState {
	return &RegexState{
		BaseState: machine.NewBaseState(),
	}
}

func (node *RegexState) GetNexts() map[byte]sets.Set[ability.State] {
	return node.Nexts
}

func (node *RegexState) Process(c byte) []ability.State {
	states := node.BaseState.Process(c)
	if extraStates, ok := node.Nexts['.']; ok {
		states = append(states, extraStates.ToSlices()...)
	}
	return states
}

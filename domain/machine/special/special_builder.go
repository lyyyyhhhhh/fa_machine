package special

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type SpecMachineBuilder struct {
	machine *machine.BaseMachine

	cur            *SpecState
	visited        sets.Set[ability.State]
	byteHandlerMap map[byte]ability.ByteHandler
}

func NewSpecMachineBuilder() *SpecMachineBuilder {
	state := newSpecState()
	builder := &SpecMachineBuilder{
		machine: &machine.BaseMachine{
			StartState: state,
			EndStates:  make(sets.Set[ability.State]),
		},
		cur:     state,
		visited: sets.Set[ability.State]{},
	}
	builder.byteHandlerMap = map[byte]ability.ByteHandler{
		'*': builder.starHandler,
		'?': builder.queHandler,
	}
	return builder
}

func (s *SpecMachineBuilder) GetInitMachine() ability.MachineAbility {
	return s.machine
}

func (s *SpecMachineBuilder) GetByteHandler(c byte) ability.ByteHandler {
	if function, ok := s.byteHandlerMap[c]; ok {
		return function
	}
	return s.defaultHandler
}

func (s *SpecMachineBuilder) starHandler(pattern string, idx int) {
	s.cur.isStringMatch = true
}

func (s *SpecMachineBuilder) queHandler(pattern string, idx int) {
	s.cur.isSingleMatch = true
}

func (s *SpecMachineBuilder) defaultHandler(pattern string, idx int) {
	cur := s.cur
	ch := pattern[idx]
	cur.AddState(ch, newSpecState())
	// 跳转到匹配该字符的下个状态
	for _, nextState := range cur.Nexts[ch].ToSlices() {
		if !s.visited.Contains(nextState) {
			s.visited.Add(nextState)
			s.cur = nextState.(*SpecState)
			break
		}
	}
}

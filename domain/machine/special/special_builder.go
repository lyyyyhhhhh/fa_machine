package special

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type SpecMachineBuilder struct {
	*machine.BaseMachineBuilder

	cur               *SpecState
	visited           sets.Set[ability.State]
	handleCharFuncMap map[byte]ability.HandleCharFunc
}

func NewSpecMachineBuilder() *SpecMachineBuilder {
	startState := newSpecState()
	builder := &SpecMachineBuilder{
		BaseMachineBuilder: machine.NewBaseMachineBuilder(startState),
		visited:            sets.NewSet[ability.State](),
	}
	builder.cur = startState
	builder.handleCharFuncMap = map[byte]ability.HandleCharFunc{
		'*': builder.handlerStarFunc,
		'?': builder.handlerQueFunc,
	}
	return builder
}

func (s *SpecMachineBuilder) GetHandlerByteFunc(c byte) ability.HandleCharFunc {
	if function, ok := s.handleCharFuncMap[c]; ok {
		return function
	}
	return s.handlerDefaultFunc
}

func (s *SpecMachineBuilder) handlerStarFunc(pattern string, idx int) {
	s.cur.isStringMatch = true
}

func (s *SpecMachineBuilder) handlerQueFunc(pattern string, idx int) {
	s.cur.isSingleMatch = true
}

func (s *SpecMachineBuilder) handlerDefaultFunc(pattern string, idx int) {
	cur := s.cur
	ch := pattern[idx]
	cur.Nexts[ch].Add(newSpecState())
	// 跳转到匹配该字符的下个状态
	for _, nextState := range cur.Nexts[ch].ToSlices() {
		if !s.visited.Contains(nextState) {
			s.visited.Add(nextState)
			s.cur = nextState.(*SpecState)
			break
		}
	}
}

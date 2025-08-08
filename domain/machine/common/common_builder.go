package common

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type CommonMachineBuilder struct {
	*machine.BaseMachineBuilder

	cur               *CommonState
	visited           sets.Set[ability.State]
	uniqueState       map[int]ability.State
	handleCharFuncMap map[byte]ability.HandleCharFunc
}

func NewCommonMachineBuilder() *CommonMachineBuilder {
	commonState := newCommonState()
	builder := &CommonMachineBuilder{
		BaseMachineBuilder: machine.NewBaseMachineBuilder(commonState),
		visited:            sets.NewSet[ability.State](),
		uniqueState:        make(map[int]ability.State),
	}
	builder.cur = commonState

	builder.handleCharFuncMap = map[byte]ability.HandleCharFunc{
		'*': builder.handlerStarFunc,
		'+': builder.handlerAddFunc,
		'?': builder.handlerQueFunc,
	}
	return builder
}

func (c *CommonMachineBuilder) getMachine() *machine.BaseMachine {
	return c.BaseMachineBuilder.Machine
}

func (c *CommonMachineBuilder) GetHandlerByteFunc(ch byte) ability.HandleCharFunc {
	if function, ok := c.handleCharFuncMap[ch]; ok {
		return function
	}
	return c.handlerDefaultFunc
}

func (c *CommonMachineBuilder) handlerStarFunc(pattern string, idx int) {
	// 表示此状态可以通过前面的字符到达原状态
	ch := pattern[idx-1]
	c.cur.AddState(ch, c.cur)
}

func (c *CommonMachineBuilder) handlerAddFunc(pattern string, idx int) {
	// 表示此状态可以通过前面的字符到达原状态
	ch := pattern[idx-1]
	c.cur.AddState(ch, c.cur)
}

func (c *CommonMachineBuilder) handlerQueFunc(pattern string, idx int) {
}

func (c *CommonMachineBuilder) handlerDefaultFunc(pattern string, idx int) {
	cur := c.cur

	ch := pattern[idx]
	cur.AddState(ch, c.getOrSetState(idx))
	count := 0
	judgeIdx := idx + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
	for judgeIdx < len(pattern) && (pattern[judgeIdx] == '*' || pattern[judgeIdx] == '?') {
		if judgeIdx+1 < len(pattern) {
			cur.AddState(pattern[judgeIdx+1], c.getOrSetState(judgeIdx+1))
		} else {
			// 最后一个字符为*, 表示可以该状态为最终状态之一
			c.getMachine().EndStates.Add(cur)
		}
		count++
		judgeIdx = idx + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
	}
	// 跳转到匹配该字符的下个状态
	for _, nextState := range cur.Nexts[ch].ToSlices() {
		if !c.visited.Contains(nextState) {
			c.visited.Add(nextState)
			cur = nextState.(*CommonState)
			break
		}
	}

	c.cur = cur
}

// 获取或创建状态
func (c *CommonMachineBuilder) getOrSetState(idx int) ability.State {
	state := c.uniqueState[idx]
	if state == nil {
		state = newCommonState()
		c.uniqueState[idx] = state
	}
	return state
}

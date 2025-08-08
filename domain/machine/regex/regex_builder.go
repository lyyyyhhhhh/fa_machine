package regex

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type RegexMachineBuilder struct {
	*machine.BaseMachineBuilder

	cur               *RegexState
	visited           sets.Set[ability.State] // 构建状态时, 防止遇到'*'绕回去
	uniqueState       map[int]ability.State   // 状态数最多等于正则表达式长度, 构建指针时防止重复
	handleCharFuncMap map[byte]ability.HandleCharFunc
}

func NewRegexMachineBuilder() *RegexMachineBuilder {
	regexState := newRegexState()
	builder := &RegexMachineBuilder{
		BaseMachineBuilder: machine.NewBaseMachineBuilder(regexState),
		visited:            sets.NewSet[ability.State](),
		uniqueState:        make(map[int]ability.State),
	}
	builder.cur = regexState

	builder.handleCharFuncMap = map[byte]ability.HandleCharFunc{
		'*': builder.handleStarFunc,
	}
	return builder
}

func (r *RegexMachineBuilder) GetHandlerByteFunc(c byte) ability.HandleCharFunc {
	if function, ok := r.handleCharFuncMap[c]; ok {
		return function
	}
	return r.handlerDefaultFunc
}

func (r *RegexMachineBuilder) getMachine() *machine.BaseMachine {
	return r.BaseMachineBuilder.Machine
}

// 遇到*, 表示此状态可以通过前面的字符到达原状态, 不推进状态
func (r *RegexMachineBuilder) handleStarFunc(pattern string, idx int) {
	cur := r.cur

	ch := pattern[idx-1]
	tmp := cur
	cur = cur.pre.(*RegexState)
	cur.RemoveState(ch, tmp)
	cur.AddState(ch, cur)
	r.visited.Add(cur) // 标记原节点访问过了

	r.cur = cur
}

// 默认字符处理方式, 推进状态
func (r *RegexMachineBuilder) handlerDefaultFunc(pattern string, idx int) {
	cur := r.cur

	ch := pattern[idx]
	cur.AddState(ch, r.getOrSetState(cur, idx))
	count := 0
	judgeIdx := idx + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符跳到后一个状态
	for judgeIdx < len(pattern) && pattern[judgeIdx] == '*' {
		if judgeIdx+1 < len(pattern) {
			cur.AddState(pattern[judgeIdx+1], r.getOrSetState(cur, judgeIdx+1))
		} else {
			// 最后一个字符为*, 表示可以该状态为最终状态之一
			r.getMachine().EndStates.Add(cur)
		}
		count++
		judgeIdx = idx + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
	}
	// 跳转到匹配该字符的下个状态
	for _, nextState := range cur.Nexts[ch].ToSlices() {
		if !r.visited.Contains(nextState) {
			r.visited.Add(nextState)
			r.cur = nextState.(*RegexState)
			break
		}
	}
}

func (r *RegexMachineBuilder) getOrSetState(cur ability.State, idx int) ability.State {
	state := r.uniqueState[idx]
	if state == nil {
		state = newRegexState()
		r.uniqueState[idx] = state
	}
	state.(*RegexState).pre = cur
	return state
}

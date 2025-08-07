package regex

import (
	"fa_machine/common/sets"
	"fa_machine/machine2/domain"
	"fa_machine/machine2/domain/ability"
	"fmt"
)

type RegexMachineBuilder struct {
	machine *domain.BaseMachine

	visited     sets.Set[ability.State] // 构建状态时, 防止遇到'*'绕回去
	uniqueState map[int]ability.State   // 状态数最多等于正则表达式长度, 构建指针时防止重复
}

func NewRegexMachineBuilder() *RegexMachineBuilder {
	return &RegexMachineBuilder{
		machine: &domain.BaseMachine{
			StartState: newRegexState(),
			EndStates:  sets.NewSet[ability.State](),
		},
		visited:     sets.NewSet[ability.State](),
		uniqueState: make(map[int]ability.State),
	}
}

func (r *RegexMachineBuilder) BuildMachine(pattern string) (ability.MachineAbility, error) {
	machine := r.machine
	cur := machine.StartState.(*RegexState)
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			// 表示此状态可以通过前面的字符到达原状态
			ch := pattern[i-1]
			tmp := cur
			cur = cur.pre.(*RegexState)
			cur.Nexts[ch].Remove(tmp)
			cur.AddState(ch, cur)
			r.visited.Add(cur) // 原节点访问过了
		default:
			cur.AddState(c, r.getOrSetFA(cur, i))
			count := 0
			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符跳到后一个状态
			for judgeIdx < len(pattern) && pattern[judgeIdx] == '*' {
				if judgeIdx+1 < len(pattern) {
					cur.AddState(pattern[judgeIdx+1], r.getOrSetFA(cur, judgeIdx+1))
				} else {
					// 最后一个字符为*, 表示可以该状态为最终状态之一
					machine.EndStates.Add(cur)
				}
				count++
				judgeIdx = i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			}
			// 跳转到匹配该字符的下个状态
			for _, nextState := range cur.Nexts[c].ToSlices() {
				if !r.visited.Contains(nextState) {
					r.visited.Add(nextState)
					cur = nextState.(*RegexState)
					break
				}
			}
		}
	}
	machine.EndStates.Add(cur)
	fmt.Println(machine.ToDot())
	return machine, nil
}

func (r *RegexMachineBuilder) getOrSetFA(cur ability.State, idx int) ability.State {
	state := r.uniqueState[idx]
	if state == nil {
		state = newRegexState()
		r.uniqueState[idx] = state
	}
	state.(*RegexState).pre = cur
	return state
}

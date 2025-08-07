package common

import (
	"fa_machine/common/sets"
	"fa_machine/machine2/domain"
	"fa_machine/machine2/domain/ability"
	"fmt"
)

type CommonMachineBuilder struct {
	machine *domain.BaseMachine

	visited     sets.Set[ability.State]
	uniqueState map[int]ability.State
}

func NewCommonMachineBuilder() *CommonMachineBuilder {
	return &CommonMachineBuilder{
		machine: &domain.BaseMachine{
			StartState: newCommonState(),
			EndStates:  sets.NewSet[ability.State](),
		},
		visited:     sets.NewSet[ability.State](),
		uniqueState: make(map[int]ability.State),
	}
}

func (b *CommonMachineBuilder) BuildMachine(pattern string) (ability.MachineAbility, error) {
	var err error // 正则表达式错误处理
	machine := b.machine
	cur := machine.StartState.(*CommonState)
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*', '+':
			// 表示此状态可以通过前面的字符到达原状态
			ch := pattern[i-1]
			cur.AddState(ch, cur)
		case '?':
			// '?'不处理
		default:
			cur.AddState(c, b.getOrSetState(i))
			count := 0
			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			for judgeIdx < len(pattern) && (pattern[judgeIdx] == '*' || pattern[judgeIdx] == '?') {
				if judgeIdx+1 < len(pattern) {
					cur.AddState(pattern[judgeIdx+1], b.getOrSetState(judgeIdx+1))
				} else {
					// 最后一个字符为*, 表示可以该状态为最终状态之一
					machine.EndStates.Add(cur)
				}
				count++
				judgeIdx = i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			}
			// 跳转到匹配该字符的下个状态
			for _, nextState := range cur.Nexts[c].ToSlices() {
				if !b.visited.Contains(nextState) {
					b.visited.Add(nextState)
					cur = nextState.(*CommonState)
					break
				}
			}
		}
	}
	machine.EndStates.Add(cur)
	fmt.Println(machine.ToDot())
	return machine, err
}

// 获取或创建状态
func (b *CommonMachineBuilder) getOrSetState(idx int) ability.State {
	state := b.uniqueState[idx]
	if state == nil {
		state = newCommonState()
		b.uniqueState[idx] = state
	}
	return state
}

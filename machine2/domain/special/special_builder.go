package special

import (
	"fa_machine/common/sets"
	"fa_machine/machine2/domain"
	"fa_machine/machine2/domain/ability"
	"fmt"
)

type SpecMachineBuilder struct {
	machine *domain.BaseMachine

	visited sets.Set[ability.State]
}

func NewSpecMachineBuilder() *SpecMachineBuilder {
	return &SpecMachineBuilder{
		machine: &domain.BaseMachine{
			StartState: newSpecState(),
			EndStates:  sets.NewSet[ability.State](),
		},
		visited: sets.NewSet[ability.State](),
	}
}

func (b *SpecMachineBuilder) BuildMachine(pattern string) (ability.MachineAbility, error) {
	machine := b.machine
	var err error //表达式错误处理
	cur := machine.StartState.(*SpecState)
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			cur.isStringMatch = true
		case '?':
			cur.isSingleMatch = true
		default:
			cur.Nexts[c].Add(newSpecState())
			// 跳转到匹配该字符的下个状态
			for _, nextState := range cur.Nexts[c].ToSlices() {
				if !b.visited.Contains(nextState) {
					b.visited.Add(nextState)
					cur = nextState.(*SpecState)
					break
				}
			}
		}
	}
	machine.EndStates.Add(cur)
	fmt.Println(machine.ToDot())
	return machine, err
}

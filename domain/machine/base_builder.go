package machine

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine/ability"
	"fmt"
)

type BaseMachineBuilder struct {
	Machine *BaseMachine
}

// NewBaseMachineBuilder 注入子类状态节点
func NewBaseMachineBuilder(startState ability.State) *BaseMachineBuilder {
	return &BaseMachineBuilder{
		Machine: &BaseMachine{
			StartState: startState,
			EndStates:  sets.NewSet[ability.State](),
		},
	}
}

func (b *BaseMachineBuilder) BuildMachine(pattern string) (ability.MachineAbility, error) {
	machine := b.Machine
	cur := machine.StartState
	for i := 0; i < len(pattern); i++ {
		b.GetHandlerByteFunc(pattern[i])(pattern, i)
	}
	machine.EndStates.Add(cur)
	fmt.Println(machine.ToDot())
	return machine, nil
}

func (b *BaseMachineBuilder) GetHandlerByteFunc(c byte) ability.HandleCharFunc {
	return nil
}

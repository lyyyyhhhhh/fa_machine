package machine

import (
	"fa_machine/domain/machine/ability"
)

type MachineBuilder struct {
	Machine *BaseMachine
	this    ability.BuildStateAbility
}

func NewMachineBuilder() *MachineBuilder {
	return &MachineBuilder{}
}

// Init 注入初始状态机以及状态转移or其他能力
func (b *MachineBuilder) Init(ext ability.BuildStateAbility) {
	b.Machine = ext.GetInitMachine().(*BaseMachine)
	b.this = ext
}

// BuildMachine 能力编排 control
func (b *MachineBuilder) BuildMachine(pattern string) (ability.MachineAbility, error) {
	// 1. 构建状态机
	machine := b.Machine
	cur := machine.StartState
	for i := 0; i < len(pattern); i++ {
		b.this.GetByteHandler(pattern[i])(pattern, i)
	}
	machine.EndStates.Add(cur)
	// 2. 打印状态机
	machine.ToDot()
	// ...
	return machine, nil
}

package service

import (
	"fa_machine/machine2/constant"
	"fa_machine/machine2/domain/ability"
	"fa_machine/machine2/domain/common"
	"fa_machine/machine2/domain/regex"
	"fa_machine/machine2/domain/special"
	"fmt"
)

var ModelToBuilder = map[constant.MachineType]func() ability.MachineBuilder{
	constant.SpecialFA:    func() ability.MachineBuilder { return special.NewSpecMachineBuilder() },
	constant.CommonSelfFA: func() ability.MachineBuilder { return regex.NewRegexMachineBuilder() },
	constant.CommonFA:     func() ability.MachineBuilder { return common.NewCommonMachineBuilder() },
}

// IsMatch 控制
func IsMatch(model string, s string, p string) bool {
	// 1. 构建状态机
	faMachine, err := build(model, p)
	if err != nil {
		return false
	}
	// 2. 匹配
	return faMachine.IsMatch(s)
}

func build(model string, pattern string) (ability.MachineAbility, error) {
	builder, ok := ModelToBuilder[constant.MachineType(model)]
	if !ok {
		return nil, fmt.Errorf("unknown constant type: %s", model)
	}

	return builder().BuildMachine(pattern)
}

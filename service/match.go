package service

import (
	"fa_machine/constant"
	"fa_machine/domain/machine/ability"
	"fa_machine/domain/machine/common"
	"fa_machine/domain/machine/regex"
	"fa_machine/domain/machine/special"
	"fmt"
)

var (
	commonBuilder            *common.CommonMachineBuilder
	regexBuilder             *regex.RegexMachineBuilder
	specBuilder              *special.SpecMachineBuilder
	ModelToBuilder           map[constant.MachineType]ability.MachineBuilder
	ModelToHandleCharAbility map[constant.MachineType]ability.HandleCharAbility
)

func init() {
	regexBuilder = regex.NewRegexMachineBuilder()
	commonBuilder = common.NewCommonMachineBuilder()
	specBuilder = special.NewSpecMachineBuilder()

	// 工厂+策略模式, 只暴露方法
	ModelToBuilder = map[constant.MachineType]ability.MachineBuilder{
		constant.SpecialFA:    specBuilder,
		constant.CommonSelfFA: regexBuilder,
		constant.CommonFA:     commonBuilder,
	}

	ModelToHandleCharAbility = map[constant.MachineType]ability.HandleCharAbility{
		constant.SpecialFA:    specBuilder,
		constant.CommonSelfFA: regexBuilder,
		constant.CommonFA:     commonBuilder,
	}
}

// IsMatch 控制
func IsMatch(model string, s string, p string) bool {
	// 1. 构建状态机
	faMachine, err := build(model, p)
	if err != nil {
		return false
	}
	// 2. 匹配
	return faMachine.IsMatch2(s)
}

func build(model string, pattern string) (ability.MachineAbility, error) {
	builderAbility, ok := ModelToBuilder[constant.MachineType(model)]
	if !ok {
		return nil, fmt.Errorf("unknown constant type: %s", model)
	}
	handleAbility, ok := ModelToHandleCharAbility[constant.MachineType(model)]
	if !ok {
		return nil, fmt.Errorf("unknown constant type: %s", model)
	}
	return builderAbility.BuildMachine(pattern, handleAbility)
}

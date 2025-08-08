package domain

import (
	"fa_machine/archive/domain/fa"
	"fa_machine/constant"
	"fmt"
)

var modelToFA = map[constant.MachineType]MachineAbility{
	constant.SpecialFA:    fa.NewSpecFA(),
	constant.CommonSelfFA: fa.NewCommonSelfFA(),
	constant.CommonFA:     fa.NewCommonFA(),
}

func NewFaByModel(model string, regex string) (MachineAbility, error) {
	if faMachine, ok := modelToFA[constant.MachineType(model)]; ok {
		return faMachine, faMachine.Build(regex)
	}
	return nil, fmt.Errorf("unknown constant type: %s", model)
}

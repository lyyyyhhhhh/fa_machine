package domain

import (
	"fmt"
	"simpleCode/machine/domain/fa"
	"simpleCode/machine/domain/faAbility"
)

type FaModelType string

const (
	SpecialFA FaModelType = "special" // leetcode 第44题
	CommonFA  FaModelType = "common"  // leetcode 第10题
)

var modelToFA = map[FaModelType]func() faAbility.AutomatonAbility{
	SpecialFA: func() faAbility.AutomatonAbility { return fa.NewSpecFA() },
	CommonFA:  func() faAbility.AutomatonAbility { return fa.NewCommonFA() },
}

func NewFaByModel(model string, regex string) (faAbility.AutomatonAbility, error) {
	faModel := FaModelType(model)
	if factory, ok := modelToFA[faModel]; ok {
		faMachine := factory()
		return faMachine, faMachine.Build(regex)
	}
	return nil, fmt.Errorf("unknown model type: %s", model)
}

package domain

import (
	"simpleCode/machine/domain/fa"
	"simpleCode/machine/domain/faAbility"
)

type FaModelType string

const (
	SpecialFA FaModelType = "special"
	CommonFA  FaModelType = "common"
)

var modelToFA = map[FaModelType]func() faAbility.AutomatonAbility{
	SpecialFA: func() faAbility.AutomatonAbility { return fa.NewSpecFA() },
	CommonFA:  func() faAbility.AutomatonAbility { return fa.NewCommonFA() },
}

func NewFaByModel(model string) faAbility.AutomatonAbility {
	faModel := FaModelType(model)
	if factory, ok := modelToFA[faModel]; ok {
		return factory()
	}
	return nil
}

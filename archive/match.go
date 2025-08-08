package archive

import (
	"fa_machine/archive/domain"
)

// IsMatch 控制
func IsMatch(model string, s string, p string) bool {
	faMachine, err := domain.NewFaByModel(model, p)
	if err != nil {
		return false
	}
	return faMachine.IsMatch(s)
}

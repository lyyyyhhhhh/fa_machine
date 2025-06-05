package machine

import (
	"simpleCode/common/sets"
	"simpleCode/machine/domain"
	"simpleCode/machine/domain/faAbility"
)

// IsMatch 控制
func IsMatch(model string, s string, p string) bool {
	startState := domain.NewFaByModel(model)
	endState := startState.Build(p) //结束状态

	curStates := sets.NewSet[faAbility.AutomatonAbility]()
	curStates.Add(startState) //初始状态
	for i := 0; i < len(s); i++ {
		c := s[i]
		nextStates := sets.NewSet[faAbility.AutomatonAbility]()
		for state := range curStates {
			// 遍历状态集合, 得到此状态通过此字符可以流转到的下个状态集合
			nextStates.AddAll(state.Process(c))
		}
		if len(nextStates) == 0 {
			return false
		}
		curStates = nextStates
	}
	return curStates.ContainsOne(endState)
}

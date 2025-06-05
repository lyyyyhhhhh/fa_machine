package fa

import (
	"simpleCode/common/utils"
	"simpleCode/machine/domain/faAbility"
)

type SpecFA struct {
	nextChar      byte // 匹配指定字符
	isSingleMatch bool // 匹配任意多个字符
	isStringMatch bool // 匹配任意单个字符
	isEnd         bool
	next          *SpecFA
}

func NewSpecFA() *SpecFA {
	return &SpecFA{}
}

func (fa *SpecFA) Build(pattern string) []faAbility.AutomatonAbility {
	cur := fa
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			cur.isStringMatch = true
		case '?':
			cur.isSingleMatch = true
			cur.next = &SpecFA{}
			cur = cur.next
		default:
			cur.nextChar = c
			cur.next = &SpecFA{}
			cur = cur.next
		}
	}
	cur.isEnd = true
	endStates := []*SpecFA{cur}
	return utils.ToInterfaceSlice[*SpecFA, faAbility.AutomatonAbility](endStates)
}

// Process 根据输入的字符进行状态流转, 保存下一步的状态
func (fa *SpecFA) Process(c byte) []faAbility.AutomatonAbility {
	var states []*SpecFA
	if fa.isStringMatch {
		states = append(states, fa)
	}
	if fa.isEnd {
		return utils.ToInterfaceSlice[*SpecFA, faAbility.AutomatonAbility](states)
	}
	if fa.isSingleMatch || fa.nextChar == c {
		states = append(states, fa.next)
	}
	// go语言不允许切片自动转换
	return utils.ToInterfaceSlice[*SpecFA, faAbility.AutomatonAbility](states)
}

package special

import (
	"fa_machine/domain/machine"
	"fa_machine/domain/machine/ability"
)

type SpecState struct {
	*machine.BaseState
	isSingleMatch bool // 匹配任意多个字符
	isStringMatch bool // 匹配任意单个字符
}

// 创建新的状态节点
func newSpecState() *SpecState {
	return &SpecState{
		BaseState: machine.NewBaseState(),
	}
}

func (node *SpecState) Process(c byte) []ability.State {
	// 1. 基本转移能力
	states := node.BaseState.Process(c)

	// 2. 特殊转移逻辑
	if node.isStringMatch {
		// 如果匹配任意长的字符, 始终返回原状态
		states = append(states, node)
	} else if node.isSingleMatch {
		// 如果匹配任意一个字符, 返回原状态后清除标记
		states = append(states, node)
		node.isSingleMatch = false
	}

	return states
}

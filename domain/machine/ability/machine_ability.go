package ability

import "fa_machine/common/sets"

// MachineAbility 自动机匹配字符串的基本能力 control 不变
type MachineAbility interface {
	IsMatch2(s string) bool // 根据输入字符串进行匹配
}

// State 不同自动机匹配到不同字符的处理方式 logic 可变
type State interface {
	Process(c byte) []State

	GetNexts() map[byte]sets.Set[State] // 打印状态节点, 与逻辑无关
}

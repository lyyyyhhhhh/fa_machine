package domain

// MachineAbility 自动机匹配字符串的基本能力 control 不变
type MachineAbility interface {
	Build(pattern string) error
	IsMatch(s string) bool // 根据输入字符串进行匹配
}

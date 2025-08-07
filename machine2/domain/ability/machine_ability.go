package ability

// MachineAbility 自动机基本能力
type MachineAbility interface {
	// IsMatch 根据输入字符串进行匹配
	IsMatch(s string) bool
}

// State 节点的状态转移能力
type State interface {
	Process(c byte) []State
}

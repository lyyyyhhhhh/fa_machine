package faAbility

// AutomatonAbility 自动机基本能力 计算
type AutomatonAbility interface {
	// Build 根据表达式构建自动机
	Build(pattern string) []AutomatonAbility
	// Process 根据输入字符进行状态转移
	Process(c byte) []AutomatonAbility
}

package faAbility

// AutomatonAbility 自动机基本能力 计算
type AutomatonAbility interface {
	// Build 根据表达式构建自动机
	Build(pattern string) error
	// IsMatch 根据输入字符串进行匹配
	IsMatch(s string) bool
}

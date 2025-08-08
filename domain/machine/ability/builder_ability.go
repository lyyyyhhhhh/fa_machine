package ability

// MachineBuilder 根据模式串构建自动机的基础逻辑 control 不变
type MachineBuilder interface {
	BuildMachine(pattern string, buildLogic HandleCharAbility) (MachineAbility, error) // 根据表达式构建自动机
}

type HandleCharFunc func(pattern string, idx int)

// HandleCharAbility 构建自动机时不同字符的处理方式 logic 可变
type HandleCharAbility interface {
	GetHandlerByteFunc(c byte) HandleCharFunc // 不同字符处理方式
}

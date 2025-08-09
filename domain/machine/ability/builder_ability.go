package ability

// MachineBuilder 根据模式串构建自动机的基础逻辑
type MachineBuilder interface {
	Init(ext BuildStateAbility)
	BuildMachine(pattern string) (MachineAbility, error) // 根据表达式构建自动机 control 不变
}

type ByteHandler func(pattern string, idx int)

// BuildStateAbility 构建自动机时不同字符的处理方式 logic 可变
type BuildStateAbility interface {
	GetInitMachine() MachineAbility
	GetByteHandler(c byte) ByteHandler // 不同字符处理方式
}

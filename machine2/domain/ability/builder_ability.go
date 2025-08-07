package ability

type MachineBuilder interface {
	// BuildMachine 根据表达式构建自动机
	BuildMachine(pattern string) (MachineAbility, error)
}

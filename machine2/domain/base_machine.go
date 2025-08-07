package domain

import (
	"fa_machine/common/sets"
	"fa_machine/machine2/domain/ability"
	"fmt"
	"strings"
)

type BaseMachine struct {
	StartState ability.State
	EndStates  sets.Set[ability.State]
}

type BaseState struct {
	Nexts map[byte]sets.Set[ability.State]
}

func NewBaseState() *BaseState {
	return &BaseState{
		Nexts: make(map[byte]sets.Set[ability.State]),
	}
}

// AddState 新增节点的转移状态
func (node *BaseState) AddState(c byte, state ability.State) {
	if _, ok := node.Nexts[c]; !ok {
		node.Nexts[c] = sets.NewSet[ability.State]()
	}
	node.Nexts[c].Add(state)
}

func (node *BaseState) GetNexts() map[byte]sets.Set[ability.State] {
	return node.Nexts
}

// Process 根据字符获取该节点的所有转移状态
func (node *BaseState) Process(c byte) []ability.State {
	return node.Nexts[c].ToSlices()
}

func (m *BaseMachine) IsMatch2(s string) bool {
	curStates := sets.NewSet[ability.State]()
	curStates.Add(m.StartState)

	for i := 0; i < len(s); i++ {
		c := s[i]
		nextStates := sets.NewSet[ability.State]()
		for state := range curStates {
			nextStates.AddAll(state.Process(c))
		}
		if len(nextStates) == 0 {
			return false
		}
		curStates = nextStates
	}

	return curStates.ContainsOne(m.EndStates.ToSlices())
}

// ToDot 画出节点状态图
func (m *BaseMachine) ToDot() string {
	var builder strings.Builder
	builder.WriteString("digraph FA {\n")
	builder.WriteString("  rankdir=LR;\n")
	builder.WriteString("  node [shape=circle, fontsize=2]\n")
	builder.WriteString("  edge [fontsize=20, fontcolor=red];\n")

	visited := make(map[ability.State]bool)
	var dfs func(node ability.State)

	nodeID := func(ptr ability.State) string {
		return fmt.Sprintf("node_%p", ptr)
	}

	dfs = func(node ability.State) {
		if visited[node] {
			return
		}
		visited[node] = true

		// 终止状态标记为双圈
		if m.EndStates.Contains(node) {
			builder.WriteString(fmt.Sprintf("  %s [shape=doublecircle];\n", nodeID(node)))
		} else {
			builder.WriteString(fmt.Sprintf("  %s;\n", nodeID(node)))
		}

		extendNode := node.(ability.ExtendableState)

		for ch, nextList := range extendNode.GetNexts() {
			for next := range nextList {
				label := string(ch)
				if ch == '.' {
					label = "."
				}
				builder.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n",
					nodeID(node), nodeID(next), label))
				dfs(next)
			}
		}
	}

	dfs(m.StartState)
	builder.WriteString("}")
	return builder.String()
}

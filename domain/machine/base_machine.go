package machine

import (
	"fa_machine/common/sets"
	"fa_machine/domain/machine/ability"
	"fmt"
	"strings"
)

type BaseMachine struct {
	StartState ability.State
	EndStates  sets.Set[ability.State]
}

func (m *BaseMachine) IsMatch2(s string) bool {
	curStates := sets.NewSet[ability.State]()
	curStates.Add(m.StartState)

	for _, c := range s {
		nextStates := sets.NewSet[ability.State]()
		for state := range curStates {
			nextStates.AddAll(state.Process(byte(c)))
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

		for ch, nextList := range node.GetNexts() {
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
	fmt.Println(builder.String())
	return builder.String()
}

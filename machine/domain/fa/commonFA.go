package fa

import (
	"fmt"
	"simpleCode/common/sets"
	"simpleCode/common/utils"
	"simpleCode/machine/domain/faAbility"
	"strings"
)

type CommonFA struct {
	uniqueState map[int]*CommonFA   // 状态数最多等于正则表达式长度
	visited     sets.Set[*CommonFA] // 构建状态时, 遇到'*'防止绕回去
	pre         *CommonFA
	nexts       map[byte][]*CommonFA //根据指定字符可以到达的状态集合, 一个字符到达的状态可能为多个, 例如a*a
	endStates   sets.Set[*CommonFA]
}

func NewCommonFA() *CommonFA {
	return &CommonFA{
		uniqueState: make(map[int]*CommonFA),
		nexts:       map[byte][]*CommonFA{},
		endStates:   make(sets.Set[*CommonFA]),
		visited:     make(sets.Set[*CommonFA]),
	}
}

func (fa *CommonFA) Build(pattern string) []faAbility.AutomatonAbility {
	cur := fa
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*', '+':
			// 表示此状态可以通过前面的字符到达原状态
			ch := pattern[i-1]
			cur.nexts[ch] = append(cur.nexts[ch], cur)
		case '?':
			// '?'不处理
		default:
			cur.nexts[c] = append(cur.nexts[c], fa.getOrSetFA(cur, i))
			count := 0
			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			for judgeIdx < len(pattern) && (pattern[judgeIdx] == '*' || pattern[judgeIdx] == '?') {
				if judgeIdx+1 < len(pattern) {
					cur.nexts[pattern[judgeIdx+1]] = append(cur.nexts[pattern[judgeIdx+1]], fa.getOrSetFA(cur, judgeIdx+1))
				} else {
					// 最后一个字符为*, 表示可以该状态为最终状态之一
					fa.endStates.Add(cur)
				}
				count++
				judgeIdx = i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			}
			// 跳转到匹配该字符的下个状态
			for _, nextState := range cur.nexts[c] {
				if !fa.visited.Contains(nextState) {
					fa.visited.Add(nextState)
					cur = nextState
					break
				}
			}
		}
	}
	fa.endStates.Add(cur)
	fmt.Println(fa.toDot())
	return utils.ToInterfaceSlice[*CommonFA, faAbility.AutomatonAbility](fa.endStates.ToSlices())
}

// Process 根据输入的字符进行状态流转, 保存下一步的状态
func (fa *CommonFA) Process(c byte) []faAbility.AutomatonAbility {
	states := fa.nexts[c]
	if extras, ok := fa.nexts['.']; ok {
		states = append(states, extras...)
	}
	return utils.ToInterfaceSlice[*CommonFA, faAbility.AutomatonAbility](states)
}

func (fa *CommonFA) getOrSetFA(cur *CommonFA, idx int) *CommonFA {
	state := fa.uniqueState[idx]
	if state == nil {
		state = NewCommonFA()
		fa.uniqueState[idx] = state
	}
	state.pre = cur
	return state
}

func (fa *CommonFA) toDot() string {
	var builder strings.Builder
	builder.WriteString("digraph FA {\n")
	builder.WriteString("  rankdir=LR;\n")
	builder.WriteString("  node [shape=circle, fontsize=2]\n")
	builder.WriteString("  edge [fontsize=20, fontcolor=red];\n")

	visited := make(map[*CommonFA]bool)
	var dfs func(node *CommonFA)

	nodeID := func(ptr *CommonFA) string {
		return fmt.Sprintf("node_%p", ptr)
	}

	dfs = func(node *CommonFA) {
		if visited[node] {
			return
		}
		visited[node] = true

		// 终止状态标记为双圈
		if fa.endStates.Contains(node) {
			builder.WriteString(fmt.Sprintf("  %s [shape=doublecircle];\n", nodeID(node)))
		} else {
			builder.WriteString(fmt.Sprintf("  %s;\n", nodeID(node)))
		}

		for ch, nextList := range node.nexts {
			for _, next := range nextList {
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

	dfs(fa)
	builder.WriteString("}")
	return builder.String()
}

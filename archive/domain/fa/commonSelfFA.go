package fa

import (
	"fa_machine/common/sets"
	"fmt"
	"strings"
)

type CommonSelfFA struct {
	uniqueState map[int]*CommonSelfFA   // 状态数最多等于正则表达式长度, 构建指针时防止重复
	visited     sets.Set[*CommonSelfFA] // 构建状态时, 防止遇到'*'绕回去
	pre         *CommonSelfFA
	nexts       map[byte]sets.Set[*CommonSelfFA] // 根据指定字符可以到达的状态集合, 一个字符到达的状态可能为多个, 例如a*a
	endStates   sets.Set[*CommonSelfFA]          // 结束状态集合, 存在'*', 0个字符也能匹配, 因此结束状态可能有多个
}

func NewCommonSelfFA() *CommonSelfFA {
	return &CommonSelfFA{
		uniqueState: make(map[int]*CommonSelfFA),
		visited:     make(sets.Set[*CommonSelfFA]),
		nexts:       map[byte]sets.Set[*CommonSelfFA]{},
		endStates:   make(sets.Set[*CommonSelfFA]),
	}
}

func (fa *CommonSelfFA) getOrSetFA(cur *CommonSelfFA, idx int) *CommonSelfFA {
	state := fa.uniqueState[idx]
	if state == nil {
		state = NewCommonSelfFA()
		fa.uniqueState[idx] = state
	}
	state.pre = cur
	return state
}

func (fa *CommonSelfFA) Build(pattern string) error {
	cur := fa
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			// 表示此状态可以通过前面的字符到达原状态
			ch := pattern[i-1]
			tmp := cur
			cur = cur.pre
			cur.nexts[ch].Remove(tmp)
			if cur.nexts[ch] == nil {
				cur.nexts[ch] = make(sets.Set[*CommonSelfFA])
			}
			cur.nexts[ch].Add(cur)
			fa.visited.Add(cur) // 原节点访问过了
		default:
			if cur.nexts[c] == nil {
				cur.nexts[c] = make(sets.Set[*CommonSelfFA])
			}
			cur.nexts[c].Add(fa.getOrSetFA(cur, i))
			count := 0
			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符跳到后一个状态
			for judgeIdx < len(pattern) && pattern[judgeIdx] == '*' {
				if judgeIdx+1 < len(pattern) {
					ch := pattern[judgeIdx+1]
					if cur.nexts[ch] == nil {
						cur.nexts[ch] = make(sets.Set[*CommonSelfFA])
					}
					cur.nexts[ch].Add(fa.getOrSetFA(cur, judgeIdx+1))
				} else {
					// 最后一个字符为*, 表示可以该状态为最终状态之一
					fa.endStates.Add(cur)
				}
				count++
				judgeIdx = i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			}
			// 跳转到匹配该字符的下个状态
			for _, nextState := range cur.nexts[c].ToSlices() {
				if !fa.visited.Contains(nextState) {
					fa.visited.Add(nextState)
					cur = nextState
					break
				}
			}
		}
	}
	fa.endStates.Add(cur)
	fmt.Println(fa.ToDot())
	return nil
}

// IsMatch 控制
func (fa *CommonSelfFA) IsMatch(s string) bool {
	cur := fa
	curStates := sets.NewSet[*CommonSelfFA]()
	curStates.Add(cur) //初始状态
	for i := 0; i < len(s); i++ {
		c := s[i]
		nextStates := sets.NewSet[*CommonSelfFA]()
		for state := range curStates {
			// 遍历状态集合, 得到此状态通过此字符可以流转到的下个状态集合
			nextStates.AddAll(state.Process(c))
		}
		if len(nextStates) == 0 {
			return false
		}
		curStates = nextStates
	}
	return curStates.ContainsOne(fa.endStates.ToSlices())
}

// Process 根据输入的字符进行状态流转, 保存下一步的状态
func (fa *CommonSelfFA) Process(c byte) []*CommonSelfFA {
	states := fa.nexts[c].ToSlices()
	if extras, ok := fa.nexts['.']; ok {
		states = append(states, extras.ToSlices()...)
	}
	return states
}

func (fa *CommonSelfFA) ToDot() string {
	var builder strings.Builder
	builder.WriteString("digraph FA {\n")
	builder.WriteString("  rankdir=LR;\n")
	builder.WriteString("  node [shape=circle, fontsize=2]\n")
	builder.WriteString("  edge [fontsize=40, fontcolor=red];\n")

	visited := make(map[*CommonSelfFA]bool)
	var dfs func(node *CommonSelfFA)

	nodeID := func(ptr *CommonSelfFA) string {
		return fmt.Sprintf("node_%p", ptr)
	}

	dfs = func(node *CommonSelfFA) {
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
			for _, next := range nextList.ToSlices() {
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

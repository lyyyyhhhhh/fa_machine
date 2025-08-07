package fa

import (
	"fmt"
	"strings"

	"fa_machine/common/sets"
)

type CommonFA struct {
	uniqueState map[int]*CommonFA    // 状态数最多等于正则表达式长度
	visited     sets.Set[*CommonFA]  // 构建状态时, 遇到'*'防止绕回去
	nexts       map[byte][]*CommonFA //根据指定字符可以到达的状态集合, 一个字符到达的状态可能为多个, 例如a*a
	endStates   sets.Set[*CommonFA]
}

type RegexMachine struct {
	StartState  *RegexState
	Visited     sets.Set[*RegexState]
	UniqueState map[int]*CommonFA
}

type RegexState struct {
	nexts map[byte][]*RegexState
}

func NewCommonFA() *CommonFA {
	return &CommonFA{
		uniqueState: make(map[int]*CommonFA),
		nexts:       map[byte][]*CommonFA{},
		endStates:   make(sets.Set[*CommonFA]),
		visited:     make(sets.Set[*CommonFA]),
	}
}

func (fa *CommonFA) Build(pattern string) error {
	var err error // 正则表达式错误处理
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
			cur.nexts[c] = append(cur.nexts[c], fa.getOrSetFA(i))
			count := 0
			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
			for judgeIdx < len(pattern) && (pattern[judgeIdx] == '*' || pattern[judgeIdx] == '?') {
				if judgeIdx+1 < len(pattern) {
					cur.nexts[pattern[judgeIdx+1]] = append(cur.nexts[pattern[judgeIdx+1]], fa.getOrSetFA(judgeIdx+1))
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
	return err
}

// IsMatch 控制
func (fa *CommonFA) IsMatch(s string) bool {
	cur := fa
	curStates := sets.NewSet[*CommonFA]()
	curStates.Add(cur) //初始状态
	for i := 0; i < len(s); i++ {
		c := s[i]
		nextStates := sets.NewSet[*CommonFA]()
		for state := range curStates {
			// 遍历状态集合, 得到此状态通过此字符可以流转到的下个状态集合
			nextStates.AddAll(state.process(c))
		}
		if len(nextStates) == 0 {
			return false
		}
		curStates = nextStates
	}
	return curStates.ContainsOne(fa.endStates.ToSlices())
}

// Process 根据输入的字符进行状态流转, 保存下一步的状态
func (fa *CommonFA) process(c byte) []*CommonFA {
	states := fa.nexts[c]
	if extras, ok := fa.nexts['.']; ok {
		states = append(states, extras...)
	}
	return states
}

func (fa *CommonFA) getOrSetFA(idx int) *CommonFA {
	state := fa.uniqueState[idx]
	if state == nil {
		state = NewCommonFA()
		fa.uniqueState[idx] = state
	}
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

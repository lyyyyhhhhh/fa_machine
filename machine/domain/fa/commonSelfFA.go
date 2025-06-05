package fa

//
//import (
//	"fmt"
//	"simpleCode/common/sets"
//	"simpleCode/common/utils"
//	"simpleCode/machine/domain/faAbility"
//	"strings"
//)
//
//type CommonFA struct {
//	uniqueState map[int]*CommonFA   // 状态数最多等于正则表达式长度, 构建指针时防止重复
//	visited     sets.Set[*CommonFA] // 构建状态时, 防止遇到'*'绕回去
//	pre         *CommonFA
//	nexts       map[byte]sets.Set[*CommonFA] // 根据指定字符可以到达的状态集合, 一个字符到达的状态可能为多个, 例如a*a
//	endStates   sets.Set[*CommonFA]          // 结束状态集合, 存在'*', 0个字符也能匹配, 因此结束状态可能有多个
//}
//
//func NewCommonFA() *CommonFA {
//	return &CommonFA{
//		uniqueState: make(map[int]*CommonFA),
//		visited:     make(sets.Set[*CommonFA]),
//		nexts:       map[byte]sets.Set[*CommonFA]{},
//		endStates:   make(sets.Set[*CommonFA]),
//	}
//}
//
//func (fa *CommonFA) getOrSetFA(cur *CommonFA, idx int) *CommonFA {
//	state := fa.uniqueState[idx]
//	if state == nil {
//		state = NewCommonFA()
//		fa.uniqueState[idx] = state
//	}
//	state.pre = cur
//	return state
//}
//
//func (fa *CommonFA) Build(pattern string) []faAbility.AutomatonAbility {
//	cur := fa
//	for i := 0; i < len(pattern); i++ {
//		c := pattern[i]
//		switch c {
//		case '*':
//			// 表示此状态可以通过前面的字符到达原状态
//			ch := pattern[i-1]
//			tmp := cur
//			cur = cur.pre
//			cur.nexts[ch].Remove(tmp)
//			if cur.nexts[ch] == nil {
//				cur.nexts[ch] = make(sets.Set[*CommonFA])
//			}
//			cur.nexts[ch].Add(cur)
//			fa.visited.Add(cur) // 原节点访问过了
//		default:
//			if cur.nexts[c] == nil {
//				cur.nexts[c] = make(sets.Set[*CommonFA])
//			}
//			cur.nexts[c].Add(fa.getOrSetFA(cur, i))
//			count := 0
//			judgeIdx := i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符跳到后一个状态
//			for judgeIdx < len(pattern) && pattern[judgeIdx] == '*' {
//				if judgeIdx+1 < len(pattern) {
//					ch := pattern[judgeIdx+1]
//					if cur.nexts[ch] == nil {
//						cur.nexts[ch] = make(sets.Set[*CommonFA])
//					}
//					cur.nexts[ch].Add(fa.getOrSetFA(cur, judgeIdx+1))
//				} else {
//					// 最后一个字符为*, 表示可以该状态为最终状态之一
//					fa.endStates.Add(cur)
//				}
//				count++
//				judgeIdx = i + 1 + count*2 // 判断指定位置是否为*, 如果为*, 表示可以从下一个字符调到后一个状态
//			}
//			// 跳转到匹配该字符的下个状态
//			for _, nextState := range cur.nexts[c].ToSlices() {
//				if !fa.visited.Contains(nextState) {
//					fa.visited.Add(nextState)
//					cur = nextState
//					break
//				}
//			}
//		}
//	}
//	fa.endStates.Add(cur)
//	fmt.Println(fa.ToDot())
//	return utils.ToInterfaceSlice[*CommonFA, faAbility.AutomatonAbility](fa.endStates.ToSlices())
//}
//
//// Process 根据输入的字符进行状态流转, 保存下一步的状态
//func (fa *CommonFA) Process(c byte) []faAbility.AutomatonAbility {
//	states := fa.nexts[c].ToSlices()
//	if extras, ok := fa.nexts['.']; ok {
//		states = append(states, extras.ToSlices()...)
//	}
//	return utils.ToInterfaceSlice[*CommonFA, faAbility.AutomatonAbility](states)
//}
//
//func (fa *CommonFA) ToDot() string {
//	var builder strings.Builder
//	builder.WriteString("digraph FA {\n")
//	builder.WriteString("  rankdir=LR;\n")
//	builder.WriteString("  node [shape=circle, fontsize=2]\n")
//	builder.WriteString("  edge [fontsize=40, fontcolor=red];\n")
//
//	visited := make(map[*CommonFA]bool)
//	var dfs func(node *CommonFA)
//
//	nodeID := func(ptr *CommonFA) string {
//		return fmt.Sprintf("node_%p", ptr)
//	}
//
//	dfs = func(node *CommonFA) {
//		if visited[node] {
//			return
//		}
//		visited[node] = true
//
//		// 终止状态标记为双圈
//		if fa.endStates.Contains(node) {
//			builder.WriteString(fmt.Sprintf("  %s [shape=doublecircle];\n", nodeID(node)))
//		} else {
//			builder.WriteString(fmt.Sprintf("  %s;\n", nodeID(node)))
//		}
//
//		for ch, nextList := range node.nexts {
//			for _, next := range nextList.ToSlices() {
//				label := string(ch)
//				if ch == '.' {
//					label = "."
//				}
//				builder.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n",
//					nodeID(node), nodeID(next), label))
//				dfs(next)
//			}
//		}
//	}
//
//	dfs(fa)
//	builder.WriteString("}")
//	return builder.String()
//}

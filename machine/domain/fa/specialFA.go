package fa

import (
	"simpleCode/common/sets"
)

type SpecFA struct {
	nextChar      byte // 匹配指定字符
	isSingleMatch bool // 匹配任意多个字符
	isStringMatch bool // 匹配任意单个字符
	isEnd         bool
	next          *SpecFA
	endStates     *sets.Set[*SpecFA]
}

func NewSpecFA() *SpecFA {
	return &SpecFA{}
}

func (fa *SpecFA) Build(pattern string) error {
	var err error //表达式错误处理
	cur := fa
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			cur.isStringMatch = true
		case '?':
			cur.isSingleMatch = true
			cur.next = &SpecFA{}
			cur = cur.next
		default:
			cur.nextChar = c
			cur.next = &SpecFA{}
			cur = cur.next
		}
	}
	cur.isEnd = true
	fa.endStates.Add(cur)
	return err
}

// IsMatch 控制
func (fa *SpecFA) IsMatch(s string) bool {
	cur := fa
	curStates := sets.NewSet[*SpecFA]()
	curStates.Add(cur) //初始状态
	for i := 0; i < len(s); i++ {
		c := s[i]
		nextStates := sets.NewSet[*SpecFA]()
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
func (fa *SpecFA) process(c byte) []*SpecFA {
	var states []*SpecFA
	if fa.isStringMatch {
		states = append(states, fa)
	}
	if fa.isEnd {
		return states
	}
	if fa.isSingleMatch || fa.nextChar == c {
		states = append(states, fa.next)
	}
	return states
}

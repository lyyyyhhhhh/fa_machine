package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	Match = 256
	Split = 257
)

type State struct {
	c        int
	out      *State
	out1     *State
	lastList int
}

var matchState = &State{c: Match}
var nstate int

func state(c int, out, out1 *State) *State {
	nstate++
	return &State{c: c, out: out, out1: out1}
}

type Ptrlist struct {
	s    **State
	next *Ptrlist
}

type Frag struct {
	start *State
	out   *Ptrlist
}

func frag(start *State, out *Ptrlist) Frag {
	return Frag{start, out}
}

func list1(outp **State) *Ptrlist {
	return &Ptrlist{s: outp}
}

func appendPtrlist(l1, l2 *Ptrlist) *Ptrlist {
	if l1 == nil {
		return l2
	}
	head := l1
	for l1.next != nil {
		l1 = l1.next
	}
	l1.next = l2
	return head
}

func patch(l *Ptrlist, s *State) {
	for l != nil {
		*l.s = s
		l = l.next
	}
}

func re2post(re string) string {
	var buf strings.Builder
	type parenFrame struct{ nalt, natom int }
	paren := make([]parenFrame, 100)
	p := 0
	nalt, natom := 0, 0

	for i := 0; i < len(re); i++ {
		c := re[i]
		switch c {
		case '(':
			if natom > 1 {
				natom--
				buf.WriteByte('.')
			}
			if p >= len(paren) {
				return ""
			}
			paren[p] = parenFrame{nalt, natom}
			p++
			nalt = 0
			natom = 0
		case '|':
			if natom == 0 {
				return ""
			}
			for natom > 1 {
				buf.WriteByte('.')
				natom--
			}
			nalt++
		case ')':
			if p == 0 || natom == 0 {
				return ""
			}
			for natom > 1 {
				buf.WriteByte('.')
				natom--
			}
			for ; nalt > 0; nalt-- {
				buf.WriteByte('|')
			}
			p--
			nalt = paren[p].nalt
			natom = paren[p].natom
			natom++
		case '*', '+', '?':
			if natom == 0 {
				return ""
			}
			buf.WriteByte(c)
		default:
			if natom > 1 {
				natom--
				buf.WriteByte('.')
			}
			buf.WriteByte(c)
			natom++
		}
	}
	if p != 0 {
		return ""
	}
	for natom > 1 {
		buf.WriteByte('.')
		natom--
	}
	for ; nalt > 0; nalt-- {
		buf.WriteByte('|')
	}
	return buf.String()
}

func post2nfa(postfix string) *State {
	stack := make([]Frag, 0, 1000)

	push := func(f Frag) {
		stack = append(stack, f)
	}
	pop := func() Frag {
		last := len(stack) - 1
		f := stack[last]
		stack = stack[:last]
		return f
	}

	for i := 0; i < len(postfix); i++ {
		c := postfix[i]
		switch c {
		default:
			s := state(int(c), nil, nil)
			push(frag(s, list1(&s.out)))
		case '.':
			e2 := pop()
			e1 := pop()
			patch(e1.out, e2.start)
			push(frag(e1.start, e2.out))
		case '|':
			e2 := pop()
			e1 := pop()
			s := state(Split, e1.start, e2.start)
			push(frag(s, appendPtrlist(e1.out, e2.out)))
		case '?':
			e := pop()
			s := state(Split, e.start, nil)
			push(frag(s, appendPtrlist(e.out, list1(&s.out1))))
		case '*':
			e := pop()
			s := state(Split, e.start, nil)
			patch(e.out, s)
			push(frag(s, list1(&s.out1)))
		case '+':
			e := pop()
			s := state(Split, e.start, nil)
			patch(e.out, s)
			push(frag(e.start, list1(&s.out1)))
		}
	}
	if len(stack) != 1 {
		return nil
	}
	e := pop()
	patch(e.out, matchState)
	return e.start
}

type List struct {
	s []*State
	n int
}

var l1, l2 List
var listid int

func addstate(l *List, s *State) {
	if s == nil || s.lastList == listid {
		return
	}
	s.lastList = listid
	if s.c == Split {
		addstate(l, s.out)
		addstate(l, s.out1)
		return
	}
	l.s[l.n] = s
	l.n++
}

func startlist(start *State, l *List) *List {
	l.n = 0
	listid++
	addstate(l, start)
	return l
}

func ismatch(l *List) bool {
	for i := 0; i < l.n; i++ {
		if l.s[i] == matchState {
			return true
		}
	}
	return false
}

func step(clist *List, c int, nlist *List) {
	listid++
	nlist.n = 0
	for i := 0; i < clist.n; i++ {
		s := clist.s[i]
		if s.c == c {
			addstate(nlist, s.out)
		}
	}
}

type DState struct {
	l     List
	next  [256]*DState
	left  *DState
	right *DState
}

var freelist *DState
var alldstates *DState
var nstates int
var maxstates = 32

func ptrcmp(a, b *State) int {
	switch {
	case a == b:
		return 0
	case a.c < b.c:
		return -1
	default:
		return 1
	}
}

func listcmp(l1, l2 *List) int {
	if l1.n < l2.n {
		return -1
	}
	if l1.n > l2.n {
		return 1
	}
	for i := 0; i < l1.n; i++ {
		if l1.s[i] != l2.s[i] {
			return ptrcmp(l1.s[i], l2.s[i])
		}
	}
	return 0
}

func allocdstate() *DState {
	if freelist != nil {
		d := freelist
		freelist = d.left
		return d
	}
	d := &DState{}
	d.l.s = make([]*State, nstate)
	return d
}

func freestates(d *DState) {
	if d == nil {
		return
	}
	freestates(d.left)
	freestates(d.right)
	d.left = freelist
	freelist = d
}

func freecache() {
	freestates(alldstates)
	alldstates = nil
	nstates = 0
}

func dstate(l *List, nextp **DState) *DState {
	sort.Slice(l.s[:l.n], func(i, j int) bool {
		return l.s[i].c < l.s[j].c
	})
	dp := &alldstates
	for *dp != nil {
		d := *dp
		switch cmp := listcmp(l, &d.l); {
		case cmp < 0:
			dp = &d.left
		case cmp > 0:
			dp = &d.right
		default:
			return d
		}
	}

	if nstates >= maxstates {
		freecache()
		dp = &alldstates
		nextp = nil
	}

	d := allocdstate()
	copy(d.l.s, l.s[:l.n])
	d.l.n = l.n
	*dp = d
	nstates++
	if nextp != nil {
		*nextp = d
	}
	return d
}

func startdstate(start *State) *DState {
	return dstate(startlist(start, &l1), nil)
}

func nextstate(d *DState, c int) *DState {
	step(&d.l, c, &l1)
	return dstate(&l1, &d.next[c])
}

func matchDFA(start *DState, s string) bool {
	d := start
	for i := 0; i < len(s); i++ {
		c := s[i]
		if d.next[c] == nil {
			d.next[c] = nextstate(d, int(c))
		}
		d = d.next[c]
	}
	return ismatch(&d.l)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: nfa regexp string...\n")
		os.Exit(1)
	}
	post := re2post(os.Args[1])
	if post == "" {
		fmt.Fprintf(os.Stderr, "bad regexp %s\n", os.Args[1])
		os.Exit(1)
	}
	start := post2nfa(post)
	if start == nil {
		fmt.Fprintf(os.Stderr, "error in post2nfa %s\n", post)
		os.Exit(1)
	}
	l1.s = make([]*State, nstate)
	l2.s = make([]*State, nstate)
	for _, arg := range os.Args[2:] {
		if matchDFA(startdstate(start), arg) {
			fmt.Println(arg)
		}
	}
}

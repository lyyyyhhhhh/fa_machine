package main

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

type PtrList struct {
	s    **State
	next *PtrList
}

type Frag struct {
	start *State
	out   *PtrList
}

func frag(start *State, out *PtrList) Frag {
	return Frag{start: start, out: out}
}

func list1(outp **State) *PtrList {
	return &PtrList{s: outp}
}

func patch(l *PtrList, s *State) {
	for ; l != nil; l = l.next {
		*l.s = s
	}
}

func appendPtrList(l1, l2 *PtrList) *PtrList {
	if l1 == nil {
		return l2
	}
	p := l1
	for p.next != nil {
		p = p.next
	}
	p.next = l2
	return l1
}

// re2post converts infix regexp to postfix with explicit concat ('.')
func re2post(re string) string {
	type paren struct {
		nalt  int
		natom int
	}
	var buf []rune
	stack := make([]paren, 0, 100)
	nalt, natom := 0, 0

	for _, r := range re {
		switch r {
		case '(':
			if natom > 1 {
				natom--
				buf = append(buf, '.')
			}
			stack = append(stack, paren{nalt, natom})
			nalt, natom = 0, 0
		case '|':
			if natom == 0 {
				return ""
			}
			for {
				natom--
				if natom <= 0 {
					break
				}
				buf = append(buf, '.')
			}
			nalt++
		case ')':
			if len(stack) == 0 || natom == 0 {
				return ""
			}
			for natom > 1 {
				buf = append(buf, '.')
				natom--
			}
			for ; nalt > 0; nalt-- {
				buf = append(buf, '|')
			}
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nalt, natom = p.nalt, p.natom+1
		case '*', '+', '?':
			if natom == 0 {
				return ""
			}
			buf = append(buf, r)
		default:
			if natom > 1 {
				natom--
				buf = append(buf, '.')
			}
			buf = append(buf, r)
			natom++
		}
	}

	if len(stack) != 0 {
		return ""
	}
	for natom > 1 {
		buf = append(buf, '.')
		natom--
	}
	for ; nalt > 0; nalt-- {
		buf = append(buf, '|')
	}
	return string(buf)
}

func post2nfa(postfix string) *State {
	stack := make([]Frag, 0, 1000)

	for _, r := range postfix {
		switch r {
		default:
			s := state(int(r), nil, nil)
			stack = append(stack, frag(s, list1(&s.out)))
		case '.':
			e2, e1 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			patch(e1.out, e2.start)
			stack = append(stack, frag(e1.start, e2.out))
		case '|':
			e2, e1 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			s := state(Split, e1.start, e2.start)
			stack = append(stack, frag(s, appendPtrList(e1.out, e2.out)))
		case '?':
			e := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			s := state(Split, e.start, nil)
			stack = append(stack, frag(s, appendPtrList(e.out, list1(&s.out1))))
		case '*':
			e := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			s := state(Split, e.start, nil)
			patch(e.out, s)
			stack = append(stack, frag(s, list1(&s.out1)))
		case '+':
			e := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			s := state(Split, e.start, nil)
			patch(e.out, s)
			stack = append(stack, frag(e.start, list1(&s.out1)))
		}
	}

	if len(stack) != 1 {
		return nil
	}
	e := stack[0]
	patch(e.out, matchState)
	return e.start
}

type List struct {
	states []*State
	n      int
}

var l1, l2 List
var listID int

func addState(l *List, s *State) {
	if s == nil || s.lastList == listID {
		return
	}
	s.lastList = listID
	if s.c == Split {
		addState(l, s.out)
		addState(l, s.out1)
	} else {
		l.states[l.n] = s
		l.n++
	}
}

func startList(start *State, l *List) *List {
	l.n = 0
	listID++
	addState(l, start)
	return l
}

func step(clist *List, c int, nlist *List) {
	listID++
	nlist.n = 0
	for i := 0; i < clist.n; i++ {
		s := clist.states[i]
		if s.c == c {
			addState(nlist, s.out)
		}
	}
}

func isMatch(l *List) bool {
	for i := 0; i < l.n; i++ {
		if l.states[i] == matchState {
			return true
		}
	}
	return false
}

func match(start *State, s string) bool {
	clist := startList(start, &l1)
	nlist := &l2
	for _, r := range s {
		step(clist, int(r), nlist)
		clist, nlist = nlist, clist
	}
	return isMatch(clist)
}

//func main() {
//	if len(os.Args) < 3 {
//		fmt.Println("usage: go run nfa.go <regexp> <string1> [string2 ...]")
//		return
//	}
//	post := re2post(os.Args[1])
//	if post == "" {
//		fmt.Println("bad regexp", os.Args[1])
//		return
//	}
//	start := post2nfa(post)
//	if start == nil {
//		fmt.Println("error in post2nfa", post)
//		return
//	}
//	l1.states = make([]*State, nstate)
//	l2.states = make([]*State, nstate)
//	for _, s := range os.Args[2:] {
//		if match(start, s) {
//			fmt.Println(s)
//		}
//	}
//}

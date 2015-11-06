package tree
import (
	"fmt"
)


type Tree struct {
	subtrees []*Tree
	node interface{}
}

type Visit func(t interface{})
type Breadth func(q []interface{})

func Postorder(t *Tree, v Visit) {
	if t == nil {
		return
	}
	for _,t := range t.subtrees {
		Preorder(t, v)
	}
	v(t.node)
}

func Preorder(t *Tree, v Visit) {
	if t == nil {
		return
	}
	v(t.node)
	for _,t := range t.subtrees {
		Preorder(t, v)
	}
}

func BreadthFirst(t *Tree, v Visit) {
	if t == nil {
		return
	}
	queue := make([]*Tree, 0)
	queue = append(queue, t)

	for len(queue) > 0 {
		q := queue[0]
		queue = queue[1:]

		v(q.node)

		if q.subtrees != nil && len(q.subtrees) > 0 {
			queue = append(queue, q.subtrees...)
		}
	}
}

func T(d interface{}, sub ...*Tree) *Tree {
	return &Tree{
		node: d,
		subtrees: sub,
	}
}

func Ex() {
	a := T("00",
			T("11"),
			T("12",
				T("23"), T("24")),
			T("13"),
			T("14"))

	visit := func(a interface{}) {
		fmt.Println(a)
	}

	Preorder(a, visit)
	fmt.Println()
	Postorder(a, visit)
	fmt.Println()

	BreadthFirst(a, visit)
}


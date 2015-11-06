package build

import (
	"os"
	"fmt"
	"sync"
)

type Group struct {
	name string
	wait *sync.WaitGroup
}
type Groups []*Group
func NewGroups(g *Group) Groups {
	return []*Group{ g }
}
func (g Groups) Add(gr ...*Group) Groups {
	return append(g, gr...)
}
func (g Groups) Len() int {
	return len(g)
}

func YamlBuild() {
	plays := ReadYaml(os.Args[1]).Load().YamlPlays()
	first := plays[0]

	in := make(CmdChannel,0)
	done := make(chan bool, 0)
	go Overhead(in, done)

	fmt.Println("Building Job:", first.Name)

	groups := Traverse(in, first)
	done <- true

	fmt.Println("Job Built:", first.Name)
	fmt.Println("# Of Groups", groups.Len())

	for _,g := range groups {
		fmt.Printf("Preparing to execute task: %s\n", g.name)
	}

	for _,g := range groups {
		if g.wait != nil {
			g.wait.Wait()
		}
	}
}

func Traverse(ch CmdChannel, t *Task) Groups {
	if t == nil {
		return nil
	}

	if t.IsParllel() {
		currGroup := &Group{
			name: t.Name,
		}
		groups := NewGroups(currGroup)

		steps := make([]*CmdStep, t.Parallel.Len())
		for i,p := range t.Parallel {
			if p.HasSteps() {
				groups = groups.Add(Traverse(ch, p)...)
			} else {
				cmd := NewScript(p.Name, p.Sh).ScriptStep()
				steps[i] = cmd
			}
		}

		currGroup.wait = Process(ch, steps...)
		return groups

	} else if t.IsSeries() {

		groups := NewGroups(&Group{
			name: t.Name,
		})

		for _,r := range t.Series {
			if r.HasSteps() {
				groups = groups.Add(Traverse(ch, r)...)
			} else {
				cmd := NewScript(r.Name, r.Sh).ScriptStep()
				wg := Process(ch, cmd)
				groups = groups.Add(&Group{
					name: r.Name,
					wait: wg,
				})
			}
		}

		return groups

	} else {
		panic("Task is incorrect. A task can be in series or parallel, but not both.")
	}
}


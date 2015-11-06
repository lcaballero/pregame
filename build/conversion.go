package build

import (
	"github.com/aymerick/raymond"
	"log"
)

type Task struct {
	Name string
	Sh string
	Series   Tasks
	Parallel Tasks
}

type Tasks []*Task

func (t *Task) IsParllel() bool {
	return t.Parallel.Len() > 0 && t.Series.Len() == 0
}

func (t *Task) IsSeries() bool {
	return t.Series.Len() > 0 && t.Parallel.Len() == 0
}

func (p *Task) HasSteps() bool {
	return p.IsParllel() || p.IsSeries()
}

func (t Tasks) Len() int {
	if t == nil {
		return 0
	} else {
		return len(t)
	}
}

func (t Tasks) First() *Task {
	if t == nil || len(t) <= 0 {
		return nil
	} else {
		return t[0]
	}
}

func (v *Value) YamlPlays() Tasks {
	tasks := make(Tasks, 1)
	tasks[0] = &Task{
		Name: v.String("name"),
		Series: v.ToSeries(),
		Parallel: v.ToParallel(),
	}
	return tasks
}

func (v *Value) ToSeries() Tasks {
	return v.Get("series").ToPlays()
}

func (v *Value) ToParallel() Tasks {
	return v.Get("parallel").ToPlays()
}

func (vip *Value) ToPlays() Tasks {
	ar := vip.ToSlice()
	rs := make([]*Task, 0)

	for i,_ := range ar {
		a := vip.In(i)
		runs := a.RunTemplating()
		for _,r := range runs {
			r.Series = a.ToSeries()
			r.Parallel = a.ToParallel()
		}
		rs = append(rs, runs...)
	}

	return rs
}

func (v *Value) RunTemplating() []*Task {
	base := &Task{
		Name: v.Get("name").ToString(),
		Sh: v.Get("sh").ToString(),
	}
	vals := v.Get("with_items").ToSlice()
	rs := make([]*Task, 0)

	for _,item := range vals {
		m := map[string]interface{}{
			"item":item,
		}
		st := &Task{
			Name: Render(base.Name, m),
			Sh: Render(base.Sh, m),
		}

		rs = append(rs, st)
	}
	if len(vals) == 0 {
		rs = append(rs, base)
	}

	return rs
}

func Render(raw string, item map[string]interface{}) string {
	res, err := raymond.Render(raw, item)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

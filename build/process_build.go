package build

import (
	"fmt"
	"os"
	"sync"
)

func wd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return dir
}

type CmdChannel chan *StepProcess
type StepProcess struct {
	step int
	max int
	cmd *CmdStep
	waitGroup *sync.WaitGroup
}
type SyncGroup struct {
	wg sync.WaitGroup
}
func (sg *SyncGroup) Add(n int) {
	sg.wg.Add(n)
}
func (sg *SyncGroup) Done() {
	sg.wg.Done()
}
func (sg *SyncGroup) Wait() {
	sg.wg.Wait()
}

func Overhead(in CmdChannel, done chan bool) {
	var c *StepProcess
	smoke := true
	handle := func(sp *StepProcess) {
		fn := MakeStepRoutine(sp)
		if smoke {
			fmt.Printf("%s\n", c.cmd.name)
			fmt.Printf("%s\n", IndentLines("   ", c.cmd.rawCode))
		} else {
			c.waitGroup.Add(1)
			go fn()
		}
	}
	for {
		select {
		case c = <-in:
			handle(c)
		case <-done:
			break
		}
	}
	for len(in) > 0 {
		handle(<-in)
	}
}

func MakeStepRoutine(step *StepProcess) func() {
	return func() {
		fmt.Println("Processing:\n", step.cmd.code)
		step.cmd.Write()
		step.cmd.Build()
		defer func() {
			step.cmd.ShellScript.Dispose()
			step.waitGroup.Done()
		}()

		err := step.cmd.Exec(step)
		if err != nil {
			fmt.Printf("Proccssing Error: %s\n", step.cmd)
		}
	}
}

func Process(ch CmdChannel, cmds ...*CmdStep) *sync.WaitGroup {
	var wg sync.WaitGroup
	var max = len(cmds)
	for i, c := range cmds {
		ch <- &StepProcess{
			step: i,
			max: max,
			cmd: c,
			waitGroup: &wg,
		}
	}
	return &wg
}

func (cmd *CmdStep) Exec(step *StepProcess) error {
	cmd.Start()
	if cmd.HasError() {
		fmt.Println("HasError: ", cmd.err)
		return cmd.err
	}
	cmd.Wait()
	if cmd.HasError() {
		return cmd.err
	} else {
		fmt.Printf("Completed %d of %d, %s",
			step.step, step.max, cmd.ShellScript.name)
		return nil
	}
}



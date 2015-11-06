package build

import (
	"fmt"
	"os/exec"
	"path/filepath"
)


// A CmdStep is a hybrid type that is either a single command step or
// a set of commands (not both).  In this way a tree of commands can be
// build and executed where everything at each level can be ran
// concurrently.
type CmdStep struct {
	ShellScript
	description string
	err  error
	cmd *exec.Cmd
	children *CmdGroup
}

type CmdGroup struct {
	name string
	description string
	kids []*CmdStep
}

func (cmd *CmdStep) String() string {
	return fmt.Sprintf("err: %s\nname: %s\ncode: %s\ndir: %s\nfile: %s\n",
		cmd.err,
		cmd.ShellScript.name,
		cmd.ShellScript.code,
		cmd.cmd.Dir,
		cmd.filename,
	)
}

func (cmd *CmdStep) Build() *CmdStep {
	cmd.cmd = cmd.toCmd()

	if cmd.ShellScript.dir != "" {
		cmd.cmd.Dir = filepath.Join(wd(), cmd.ShellScript.dir)
	} else {
		cmd.cmd.Dir = wd()
	}

	return cmd
}

func (cmd *CmdStep) toCmd() *exec.Cmd {
	return exec.Command(cmd.ShellScript.cmd, cmd.ShellScript.filename)
}

func (cmd *CmdStep) Wait() error {
	cmd.err = cmd.cmd.Wait()
	return cmd.err
}

func (step *CmdStep) Start() {
	step.err = step.cmd.Start()
}

func (step *CmdStep) HasError() bool {
	return step.err != nil
}

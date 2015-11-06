package build

import (
	"fmt"
	"os"
	"io/ioutil"
	"crypto/md5"
	"io"
	"encoding/hex"
)

func Script(s string) string {
	sh := `
#!/bin/bash

%s

`
	return fmt.Sprintf(sh, s)
}

type ShellScript struct {
	name string
	code string
	cmd string
	filename string
	dir string
	rawCode string
}

func NewScript(name, code string) ShellScript {
	return ShellScript{
		cmd: "/bin/sh",
		name: name,
		code: Script(code),
		rawCode: code,
	}
}

func (ss *ShellScript) Dispose() error {
	return os.Remove(ss.filename)
}

func (ss *ShellScript) Write() error {
	h := md5.New()
	io.WriteString(h, ss.code)
	src := h.Sum(nil)

	md5hash := fmt.Sprintf("%s.sh", hex.EncodeToString(src))

	file, err := ioutil.TempFile("", md5hash)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	ss.filename = file.Name()

	_, err = file.Write([]byte(ss.code))
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

func (ss ShellScript) InDir(dir string) ShellScript {
	ss.dir = dir
	return ss
}

func (ss ShellScript) ScriptStep() *CmdStep {
	return &CmdStep{
		ShellScript: ss,
	}
}


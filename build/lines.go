package build

import (
	"io/ioutil"
	"strings"
	"bytes"
	"bufio"
)

func IndentLines(indent, code string) string {
	reader := bufio.NewReader(bytes.NewBufferString(code))
	buf := bytes.NewBufferString("")
	line, prefix, err := reader.ReadLine()

	for err == nil && line != nil {
		if !prefix {
			buf.WriteString(indent)
		}
		s := strings.Trim(string(line), " \t\r\n")
		buf.WriteString(s)
		if !prefix {
			buf.WriteString("\n")
		}
		line, prefix, err = reader.ReadLine()
	}

	return buf.String()
}

func ReadLines(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}
	n := strings.Split(string(content), "\n")
	lines := make([]string, 0, len(n))

	for _, m := range n {
		line := strings.Trim(m, " \t\r\n")
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func SimpleCmd(cmd string) []string {
	return strings.Split(cmd, " ")
}

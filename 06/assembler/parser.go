package main

import (
	"bufio"
	"io"
	"strings"
)

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) Parse(input io.Reader) (cmds []*Command) {
	lines := parseInput(input)

	for _, l := range lines {
		cmds = append(cmds, parseLine(l))
	}

	return cmds
}

func parseInput(input io.Reader) (lines []string) {
	s := bufio.NewScanner(input)
	for s.Scan() {
		line := s.Text()
		line = strings.ReplaceAll(line, " ", "")         // 空白除去
		if i := strings.LastIndex(line, "//"); i != -1 { // コメントアウトを除去
			line = line[:i]
		}
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func parseLine(line string) *Command {
	cmd := new(Command)

	if line[0] == '@' {
		cmd.Type = A_COMMAND
		cmd.Symbol = line[1:]
		return cmd
	}

	cmd.Type = C_COMMAND
	cmd.Comp = line
	if ss := strings.Split(line, "="); len(ss) == 2 {
		cmd.Dest = ss[0]
		cmd.Comp = strings.ReplaceAll(cmd.Comp, cmd.Dest+"=", "")
	}
	if ss := strings.Split(line, ";"); len(ss) == 2 {
		cmd.Jump = ss[1]
		cmd.Comp = strings.ReplaceAll(cmd.Comp, ";"+cmd.Jump, "")
	}
	return cmd
}

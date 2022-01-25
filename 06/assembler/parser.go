package main

import (
	"strings"
)

type cmdType uint8

const (
	A_COMMAND cmdType = iota
	C_COMMAND
)

type Parser struct{}

type Command struct {
	Type             cmdType
	Symbol           string
	Dest, Comp, Jump string
}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) Parse(s string) *Command {
	line := trimLine(s)
	if line == "" {
		return nil
	}

	cmd := new(Command)

	if s[0] == '@' {
		cmd.Type = A_COMMAND
		cmd.Symbol = s[1:]
		return cmd
	}

	cmd.Type = C_COMMAND
	cmd.Comp = s
	if ss := strings.Split(s, "="); len(ss) == 2 {
		cmd.Dest = ss[0]
		cmd.Comp = strings.ReplaceAll(cmd.Comp, cmd.Dest+"=", "")
	}
	if ss := strings.Split(s, ";"); len(ss) == 2 {
		cmd.Jump = ss[1]
		cmd.Comp = strings.ReplaceAll(cmd.Comp, ";"+cmd.Jump, "")
	}
	return cmd
}

func parseACommand(cmd *Command, str string) (*Command, error) {
	return cmd, nil
}

func trimLine(s string) string {
	s = strings.ReplaceAll(s, " ", "")            // 空白除去
	if i := strings.LastIndex(s, "//"); i != -1 { // コメントアウトを除去
		s = s[:i]
	}
	return s
}

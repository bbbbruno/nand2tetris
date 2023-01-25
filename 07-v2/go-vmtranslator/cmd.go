package vmtranslate

import (
	"strconv"

	"golang.org/x/exp/slices"
)

type CmdType int

const (
	Arithmetic CmdType = iota + 1
	Push
	Pop
)

type Cmd struct {
	Type    CmdType
	Command string
	Segment string
	Index   int
}

var arithmeticCommands = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}

func IsArithmeticCmd(s string) bool {
	return slices.Contains(arithmeticCommands, s)
}

func NewArithmeticCmd(command string) *Cmd {
	return &Cmd{Type: Arithmetic, Command: command}
}

var segments = []string{"argument", "local", "static", "constant", "this", "that", "pointer", "temp"}

func IsPushPopCmd(command, segment, index string) bool {
	if command != "push" && command != "pop" {
		return false
	}
	if !slices.Contains(segments, segment) {
		return false
	}
	if i, err := strconv.Atoi(index); i < 0 || err != nil {
		return false
	}

	return true
}

func NewPushCmd(command, segment, index string) *Cmd {
	i, _ := strconv.Atoi(index)
	return &Cmd{Type: Push, Command: command, Segment: segment, Index: i}
}

func NewPopCmd(command, segment, index string) *Cmd {
	i, _ := strconv.Atoi(index)
	return &Cmd{Type: Pop, Command: command, Segment: segment, Index: i}
}

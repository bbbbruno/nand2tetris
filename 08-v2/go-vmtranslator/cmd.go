package vmtranslate

import (
	"regexp"
	"strconv"

	"golang.org/x/exp/slices"
)

type CmdType int

const (
	Arithmetic CmdType = iota + 1
	Push
	Pop
	Flow
	Function
	Return
	Call
)

type Cmd struct {
	Type    CmdType
	Command string
	Arg1    string
	Arg2    int
}

var arithmeticCommands = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}

func IsArithmeticCmd(s string) bool {
	return slices.Contains(arithmeticCommands, s)
}

func NewArithmeticCmd(command string) *Cmd {
	return &Cmd{Type: Arithmetic, Command: command}
}

var segments = []string{"argument", "local", "static", "constant", "this", "that", "pointer", "temp"}

func IsPushPopCmd(command, arg1, arg2 string) bool {
	if command != "push" && command != "pop" {
		return false
	}
	if !slices.Contains(segments, arg1) {
		return false
	}
	if i, err := strconv.Atoi(arg2); i < 0 || err != nil {
		return false
	}

	return true
}

func NewPushCmd(command, arg1, arg2 string) *Cmd {
	i, _ := strconv.Atoi(arg2)
	return &Cmd{Type: Push, Command: command, Arg1: arg1, Arg2: i}
}

func NewPopCmd(command, arg1, arg2 string) *Cmd {
	i, _ := strconv.Atoi(arg2)
	return &Cmd{Type: Pop, Command: command, Arg1: arg1, Arg2: i}
}

var flowCommands = []string{"label", "goto", "if-goto"}

func IsFlowCmd(command, arg1 string) bool {
	if !slices.Contains(flowCommands, command) {
		return false
	}
	if matched, err := regexp.MatchString(`^[^0-9][a-zA-Z0-9_.:]*$`, arg1); !matched || err != nil {
		return false
	}

	return true
}

func NewFlowCmd(command, arg1 string) *Cmd {
	return &Cmd{Type: Flow, Command: command, Arg1: arg1}
}

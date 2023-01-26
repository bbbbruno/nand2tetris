package vmtranslate

import (
	"strconv"
)

type CmdType int

const (
	Arithmetic CmdType = iota + 1
	Push
	Pop
	Flow
	Function
)

type Cmd struct {
	Type    CmdType
	Command string
	Arg1    string
	Arg2    int
}

func NewArithmeticCmd(command string) *Cmd {
	return &Cmd{Type: Arithmetic, Command: command}
}

func NewPushCmd(command, arg1, arg2 string) *Cmd {
	i, _ := strconv.Atoi(arg2)
	return &Cmd{Type: Push, Command: command, Arg1: arg1, Arg2: i}
}

func NewPopCmd(command, arg1, arg2 string) *Cmd {
	i, _ := strconv.Atoi(arg2)
	return &Cmd{Type: Pop, Command: command, Arg1: arg1, Arg2: i}
}

func NewFlowCmd(command, arg1 string) *Cmd {
	return &Cmd{Type: Flow, Command: command, Arg1: arg1}
}

func NewFunctionCmd(command, arg1, arg2 string) *Cmd {
	if command == "return" {
		return &Cmd{Type: Function, Command: command}
	} else {
		i, _ := strconv.Atoi(arg2)
		return &Cmd{Type: Function, Command: command, Arg1: arg1, Arg2: i}
	}
}

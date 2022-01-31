package vmcommand

import (
	"strconv"
	"vmconverter/utils/slices"
)

type CommandType int

const (
	C_ARITHMETIC CommandType = iota + 1
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

var arithmeticCommands = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}
var memoryAccessCommands = []string{"", "", "push", "pop", "label", "goto", "if-goto", "function", "return", "call"}

type VMCommand interface {
	CommandType() CommandType
	Instruction() string
	Arg1() string
	Arg2() int
	IsArithmetic() bool
}
type vmCommand struct {
	instruction, arg1, arg2 string
}

func NewVMCommand(instruction, arg1, arg2 string) VMCommand {
	return &vmCommand{instruction, arg1, arg2}
}

func (cmd *vmCommand) CommandType() CommandType {
	if cmd.IsArithmetic() {
		return C_ARITHMETIC
	}
	if i := slices.FindIndex(memoryAccessCommands, cmd.instruction); i != -1 {
		return CommandType(i)
	}

	return 0
}

// VMコマンドの種類を返す。
func (cmd *vmCommand) Instruction() string {
	return cmd.instruction
}

// VMコマンドの第一引数を返す。
// VMコマンドの種類が算術コマンド（C_ARITHMETIC）である場合はコマンド自体を返す。
func (cmd *vmCommand) Arg1() string {
	if cmd.CommandType() == C_ARITHMETIC {
		return cmd.instruction
	} else {
		return cmd.arg1
	}
}

// VMコマンドの第二引数を数値で返す。
func (cmd *vmCommand) Arg2() int {
	i, _ := strconv.Atoi(cmd.arg2)
	return i
}

// 算術コマンドかどうか
func (cmd *vmCommand) IsArithmetic() bool {
	return slices.Contains(arithmeticCommands, cmd.instruction)
}

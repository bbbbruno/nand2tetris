package vmtranslate

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"
)

func Parse(line string) (cmd *Cmd, err error) {
	var command, arg1, arg2 string

	ss := strings.Split(line, " ")
	command = ss[0]
	if len(ss) > 1 {
		arg1 = ss[1]
	}
	if len(ss) > 2 {
		arg2 = ss[2]
	}

	switch command {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		if !IsValidArithmeticCmd(command, arg1, arg2) {
			return nil, xerrors.Errorf("Invalid arithmetic command: %s", line)
		}
		cmd = NewArithmeticCmd(command)
	case "push", "pop":
		if !IsValidPushPopCmd(command, arg1, arg2) {
			return nil, xerrors.Errorf("Invalid push/pop command: %s", line)
		}
		if command == "push" {
			cmd = NewPushCmd(command, arg1, arg2)
		} else {
			cmd = NewPopCmd(command, arg1, arg2)
		}
	case "label", "goto", "if-goto":
		if !IsValidFlowCmd(command, arg1, arg2) {
			return nil, xerrors.Errorf("Invalid program flow command: %s", line)
		}
		cmd = NewFlowCmd(command, arg1)
	case "function", "call", "return":
		if !IsValidFlowCmd(command, arg1, arg2) {
			return nil, xerrors.Errorf("Invalid program flow command: %s", line)
		}
		cmd = NewFunctionCmd(command, arg1, arg2)
	default:
		return nil, xerrors.Errorf("Invalid command: %s", line)
	}

	return cmd, nil
}

var arithmeticCommands = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}

func IsValidArithmeticCmd(command, arg1, arg2 string) bool {
	return slices.Contains(arithmeticCommands, command) && arg1 == "" && arg2 == ""
}

var segments = []string{"argument", "local", "static", "constant", "this", "that", "pointer", "temp"}

func IsValidPushPopCmd(command, arg1, arg2 string) bool {
	if command != "push" && command != "pop" {
		return false
	}
	if !slices.Contains(segments, arg1) {
		return false
	}
	if command == "pop" && arg1 == "constant" {
		return false
	}
	if i, err := strconv.Atoi(arg2); i < 0 || err != nil {
		return false
	}

	return true
}

var flowCommands = []string{"label", "goto", "if-goto"}

func IsValidFlowCmd(command, arg1, arg2 string) bool {
	if arg2 != "" {
		return false
	}
	if !slices.Contains(flowCommands, command) {
		return false
	}
	if matched, err := regexp.MatchString(`^[^0-9][a-zA-Z0-9_.:]*$`, arg1); !matched || err != nil {
		return false
	}

	return true
}

var functionCommands = []string{"function", "call", "return"}

func IsValidFunctionCmd(command, arg1, arg2 string) bool {
	if !slices.Contains(functionCommands, command) {
		return false
	}
	if (command == "function" || command == "call") && arg2 == "" {
		return false
	} else if command == "return" && (arg1 != "" || arg2 != "") {
		return false
	}

	return true
}

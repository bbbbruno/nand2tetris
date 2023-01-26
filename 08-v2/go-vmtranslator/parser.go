package vmtranslate

import (
	"strings"

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

	switch {
	case IsArithmeticCmd(command):
		return NewArithmeticCmd(command), nil
	case IsPushPopCmd(command, arg1, arg2):
		if command == "push" {
			return NewPushCmd(command, arg1, arg2), nil
		} else {
			return NewPopCmd(command, arg1, arg2), nil
		}
	case IsFlowCmd(command, arg1):
		return NewFlowCmd(command, arg1), nil
	default:
		return nil, xerrors.Errorf("Invalid command: %s", line)
	}
}

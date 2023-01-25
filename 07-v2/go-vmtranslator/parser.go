package vmtranslate

import (
	"strings"

	"golang.org/x/xerrors"
)

func Parse(line string) (cmd *Cmd, err error) {
	var command, segment, index string

	ss := strings.Split(line, " ")
	command = ss[0]
	if len(ss) > 1 {
		segment = ss[1]
	}
	if len(ss) > 2 {
		index = ss[2]
	}

	switch {
	case IsArithmeticCmd(command):
		return NewArithmeticCmd(command), nil
	case IsPushPopCmd(command, segment, index):
		if command == "push" {
			return NewPushCmd(command, segment, index), nil
		} else {
			return NewPopCmd(command, segment, index), nil
		}
	default:
		return nil, xerrors.Errorf("Invalid command: %s", line)
	}
}

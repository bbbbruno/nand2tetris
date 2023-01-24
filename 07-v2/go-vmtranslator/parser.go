package vmtranslate

import (
	"strings"

	"golang.org/x/xerrors"
)

func Parse(line string) (cmd Cmd, err error) {
	ss := strings.Split(line, " ")

	switch command := ss[0]; {
	case IsArithmeticCmd(command) && len(ss) == 1:
		cmd, err = NewArithmeticCmd(command)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
	case IsPushPopCmd(command) && len(ss) == 3:
		segment, index := ss[1], ss[2]
		if command == "push" {
			cmd, err = NewPushCmd(command, segment, index)
		} else {
			cmd, err = NewPopCmd(command, segment, index)
		}
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
	default:
		return nil, xerrors.Errorf("Invalid command: %s", line)
	}

	return cmd, nil
}

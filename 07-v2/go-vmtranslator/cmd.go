package vmtranslate

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"
)

var randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

type Cmd interface {
	Translate() string
}

type ArithmeticCmd struct {
	Command string
	r       *rand.Rand
}

var arithmeticCommands []string = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}

func IsArithmeticCmd(s string) bool {
	return slices.Contains(arithmeticCommands, s)
}

func NewArithmeticCmd(command string) (*ArithmeticCmd, error) {
	if !IsArithmeticCmd(command) {
		return nil, xerrors.Errorf("Failed to create ArithmeticCmd: invalid command %s", command)
	}

	return &ArithmeticCmd{Command: command, r: randomizer}, nil
}

func (c *ArithmeticCmd) Translate() string {
	switch c.Command {
	case "eq", "gt", "lt":
		return fmt.Sprintf(COMPARE_ASSEMBLY, c.r.Intn(1_000_000), c.r.Intn(1_000_000), jmpAssembly[c.Command])
	default:
		return arithmeticAssembly[c.Command]
	}
}

type PushPopCmd struct {
	Command string
	Segment string
	Index   int
}

var segments []string = []string{"argument", "local", "static", "constant", "this", "that", "pointer", "temp"}

func IsPushPopCmd(s string) bool {
	return s == "push" || s == "pop"
}

func NewPushPopCmd(command, segment, index string) (*PushPopCmd, error) {
	if !IsPushPopCmd(command) {
		return nil, xerrors.Errorf("Failed to create PushPopCmd: invalid command %s", command)
	}
	if !slices.Contains(segments, segment) {
		return nil, xerrors.Errorf("Failed to create PushPopCmd: invalid segment %s", segment)
	}
	i, err := strconv.Atoi(index)
	if err != nil {
		return nil, xerrors.Errorf("Failed to create PushPopCmd: invalid index %d", i)
	} else if i < 0 {
		return nil, xerrors.Errorf("Failed to create PushPopCmd: index %d must be greater than 0", i)
	}

	return &PushPopCmd{Command: command, Segment: segment, Index: i}, nil
}

func (c *PushPopCmd) Translate() string {
	if c.Command == "push" {
		return c.segmentAssembly() + push()
	} else {
		return ""
	}
}

func (c *PushPopCmd) segmentAssembly() string {
	switch c.Segment {
	case "constant":
		return fmt.Sprintf(`@%d
D=A
`, c.Index)
	default:
		return ""
	}
}

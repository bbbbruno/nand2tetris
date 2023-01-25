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

type PushCmd struct {
	*PushPopCmd
}

type PopCmd struct {
	*PushPopCmd
}

var segments []string = []string{"argument", "local", "static", "constant", "this", "that", "pointer", "temp"}

func IsPushPopCmd(s string) bool {
	return s == "push" || s == "pop"
}

func NewPushCmd(command, segment, index string) (*PushCmd, error) {
	i, err := validatePushPopCmd(command, segment, index)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &PushCmd{&PushPopCmd{Command: command, Segment: segment, Index: i}}, nil
}

func NewPopCmd(command, segment, index string) (*PopCmd, error) {
	i, err := validatePushPopCmd(command, segment, index)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &PopCmd{&PushPopCmd{Command: command, Segment: segment, Index: i}}, nil
}

func validatePushPopCmd(command, segment, index string) (int, error) {
	if !IsPushPopCmd(command) {
		return 0, xerrors.Errorf("Failed to create PushPopCmd: invalid command %s", command)
	}
	if !slices.Contains(segments, segment) {
		return 0, xerrors.Errorf("Failed to create PushPopCmd: invalid segment %s", segment)
	}
	i, err := strconv.Atoi(index)
	if err != nil {
		return 0, xerrors.Errorf("Failed to create PushPopCmd: invalid index %d", i)
	} else if i < 0 {
		return 0, xerrors.Errorf("Failed to create PushPopCmd: index %d must be greater than 0", i)
	}

	return i, nil
}

func (c *PushCmd) Translate() string {
	return c.indexAssembly() + c.segmentAssembly() + push()
}

func (c *PushCmd) segmentAssembly() string {
	sym := symbolAssembly[c.Segment]
	switch c.Segment {
	case "constant":
		return ""
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%s
A=M+D
D=M
`, sym)
	case "pointer", "temp":
		return fmt.Sprintf(`@%s
A=A+D
D=M
`, sym)
	default:
		return ""
	}
}

func (c *PopCmd) Translate() string {
	return c.indexAssembly() + c.segmentAssembly() + pop("M") + `@R13
A=M
M=D
`
}

func (c *PopCmd) segmentAssembly() string {
	sym := symbolAssembly[c.Segment]
	switch c.Segment {
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%s
D=M+D
@R13
M=D
`, sym)
	case "pointer", "temp":
		return fmt.Sprintf(`@%s
D=A+D
@R13
M=D
`, sym)
	default:
		return ""
	}
}

func (c *PushPopCmd) indexAssembly() string {
	return fmt.Sprintf(`@%d
D=A
`, c.Index)
}

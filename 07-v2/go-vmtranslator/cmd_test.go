package vmtranslate_test

import (
	"fmt"
	"math/rand"
	"testing"
	"vmtranslate"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewArithmeticCmd(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   vmtranslate.Cmd
		hasErr bool
	}{
		{
			name:   "create add command successfully",
			arg:    "add",
			want:   &vmtranslate.ArithmeticCmd{Command: "add"},
			hasErr: false,
		},
		{
			name:   "create not command successfully",
			arg:    "not",
			want:   &vmtranslate.ArithmeticCmd{Command: "not"},
			hasErr: false,
		},
		{
			name:   "can't create arithmetic command",
			arg:    "plus",
			want:   (*vmtranslate.ArithmeticCmd)(nil),
			hasErr: true,
		},
	}

	opt := cmpopts.IgnoreUnexported(vmtranslate.ArithmeticCmd{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vmtranslate.NewArithmeticCmd(tt.arg)
			if (err != nil) != tt.hasErr {
				t.Errorf("NewArithmeticCmd() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if !cmp.Equal(got, tt.want, opt) {
				t.Errorf("NewArithmeticCmd() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestArithmeticCmdTranslate(t *testing.T) {
	tests := []struct {
		name   string
		cmd    *vmtranslate.ArithmeticCmd
		want   string
		hasErr bool
	}{
		{
			name: "translate add command",
			cmd:  &vmtranslate.ArithmeticCmd{Command: "add"},
			want: `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D+M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			name: "translate eq command successfully",
			cmd: func() *vmtranslate.ArithmeticCmd {
				cmd := &vmtranslate.ArithmeticCmd{Command: "eq"}
				cmd.ExportSetRandomizer(rand.New(rand.NewSource(100)))
				return cmd
			}(),
			want: func() string {
				r := rand.New(rand.NewSource(100))
				return fmt.Sprintf(`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@TRUE%[1]d
D;JEQ
D=0
@FINALLY%[2]d
0;JMP
(TRUE%[1]d)
D=-1
(FINALLY%[2]d)
@SP
A=M
M=D
@SP
M=M+1
`, r.Intn(1_000_000), r.Intn(1_000_000))
			}(),
		},
		{
			name: "translate lt command successfully",
			cmd: func() *vmtranslate.ArithmeticCmd {
				cmd := &vmtranslate.ArithmeticCmd{Command: "lt"}
				cmd.ExportSetRandomizer(rand.New(rand.NewSource(100)))
				return cmd
			}(),
			want: func() string {
				r := rand.New(rand.NewSource(100))
				return fmt.Sprintf(`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@TRUE%[1]d
D;JLT
D=0
@FINALLY%[2]d
0;JMP
(TRUE%[1]d)
D=-1
(FINALLY%[2]d)
@SP
A=M
M=D
@SP
M=M+1
`, r.Intn(1_000_000), r.Intn(1_000_000))
			}(),
		},
		{
			name: "translate and command",
			cmd:  &vmtranslate.ArithmeticCmd{Command: "and"},
			want: `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D&M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			name: "translate not command",
			cmd:  &vmtranslate.ArithmeticCmd{Command: "not"},
			want: `@SP
M=M-1
A=M
D=!M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.Translate()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestNewPushPopCmd(t *testing.T) {
	type args struct {
		command, segment, index string
	}
	tests := []struct {
		name string
		args
		want   vmtranslate.Cmd
		hasErr bool
	}{
		{
			name:   "create push constant 3 command successfully",
			args:   args{command: "push", segment: "constant", index: "3"},
			want:   &vmtranslate.PushPopCmd{Command: "push", Segment: "constant", Index: 3},
			hasErr: false,
		},
		{
			name:   "create pop constant 0 command successfully",
			args:   args{command: "pop", segment: "constant", index: "0"},
			want:   &vmtranslate.PushPopCmd{Command: "pop", Segment: "constant", Index: 0},
			hasErr: false,
		},
		{
			name:   "can't create command with invalid segment",
			args:   args{command: "push", segment: "const", index: "0"},
			want:   (*vmtranslate.PushPopCmd)(nil),
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vmtranslate.NewPushPopCmd(tt.args.command, tt.args.segment, tt.args.index)
			if (err != nil) != tt.hasErr {
				t.Errorf("NewPushPopCmd() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NewPushPopCmd() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPushPopCmdTranslate(t *testing.T) {
	tests := []struct {
		name   string
		cmd    *vmtranslate.PushPopCmd
		want   string
		hasErr bool
	}{
		{
			name: "translate push constant command",
			cmd:  &vmtranslate.PushPopCmd{Command: "push", Segment: "constant", Index: 10},
			want: fmt.Sprintf(`@%d
D=A
@SP
A=M
M=D
@SP
M=M+1
`, 10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.Translate()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

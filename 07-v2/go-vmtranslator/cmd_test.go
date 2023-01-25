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
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
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

func TestNewPushCmd(t *testing.T) {
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
			name:   "create push constant 10 command successfully",
			args:   args{command: "push", segment: "constant", index: "10"},
			want:   &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "constant", Index: 10}},
			hasErr: false,
		},
		{
			name:   "create push local 1 command successfully",
			args:   args{command: "push", segment: "local", index: "1"},
			want:   &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "local", Index: 1}},
			hasErr: false,
		},
		{
			name:   "create push argument 0 command successfully",
			args:   args{command: "push", segment: "argument", index: "0"},
			want:   &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "argument", Index: 0}},
			hasErr: false,
		},
		{
			name:   "can't create command with invalid segment",
			args:   args{command: "push", segment: "const", index: "0"},
			want:   (*vmtranslate.PushCmd)(nil),
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vmtranslate.NewPushCmd(tt.args.command, tt.args.segment, tt.args.index)
			if (err != nil) != tt.hasErr {
				t.Errorf("NewPushCmd() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestPushCmdTranslate(t *testing.T) {
	tests := []struct {
		name   string
		cmd    *vmtranslate.PushCmd
		want   string
		hasErr bool
	}{
		{
			name: "translate push constant command",
			cmd:  &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "constant", Index: 10}},
			want: `@10
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			name: "translate push local command",
			cmd:  &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "local", Index: 1}},
			want: `@1
D=A
@LCL
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			name: "translate push pointer command",
			cmd:  &vmtranslate.PushCmd{&vmtranslate.PushPopCmd{Command: "push", Segment: "pointer", Index: 1}},
			want: `@1
D=A
@R3
A=A+D
D=M
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

func TestPopCmdTranslate(t *testing.T) {
	tests := []struct {
		name   string
		cmd    *vmtranslate.PopCmd
		want   string
		hasErr bool
	}{

		{
			name: "translate pop local command",
			cmd:  &vmtranslate.PopCmd{&vmtranslate.PushPopCmd{Command: "pop", Segment: "local", Index: 1}},
			want: `@1
D=A
@LCL
D=M+D
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			name: "translate pop pointer command",
			cmd:  &vmtranslate.PopCmd{&vmtranslate.PushPopCmd{Command: "pop", Segment: "pointer", Index: 1}},
			want: `@1
D=A
@R3
D=A+D
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
A=M
M=D
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

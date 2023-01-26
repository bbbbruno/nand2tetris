package vmtranslate_test

import (
	"fmt"
	"math/rand"
	"testing"
	"vmtranslate"

	"github.com/google/go-cmp/cmp"
)

var tr = vmtranslate.NewTranslator("Test")

func TestTranslateArithmeticCmd(t *testing.T) {
	r1, r2 := rand.New(rand.NewSource(100)), rand.New(rand.NewSource(100))
	tr.ExportSetRandomizer(r1)
	tests := []struct {
		name   string
		arg    *vmtranslate.Cmd
		want   string
		hasErr bool
	}{
		{
			name: "translate add command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "add"},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "eq"},
			want: fmt.Sprintf(`@SP
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
`, r2.Intn(1_000_000), r2.Intn(1_000_000)),
		},
		{
			name: "translate lt command successfully",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "lt"},
			want: fmt.Sprintf(`@SP
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
`, r2.Intn(1_000_000), r2.Intn(1_000_000)),
		},
		{
			name: "translate and command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "and"},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "not"},
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
			got := tr.Translate(tt.arg)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestTranslatePushCmd(t *testing.T) {
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		want string
	}{
		{
			name: "translate push constant command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "constant", Arg2: 10},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "local", Arg2: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "pointer", Arg2: 1},
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
		{
			name: "translate push static command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "static", Arg2: 3},
			want: `@Test.3
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
			got := tr.Translate(tt.arg)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestTranslatePopCmd(t *testing.T) {
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		want string
	}{

		{
			name: "translate pop local command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Arg1: "local", Arg2: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Arg1: "pointer", Arg2: 1},
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
		{
			name: "translate pop static command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Arg1: "static", Arg2: 3},
			want: `@SP
M=M-1
A=M
D=M
@Test.3
M=D
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.Translate(tt.arg)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestTranslateFlowCmd(t *testing.T) {
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		fn   string
		want string
	}{
		{
			name: "translate label command",
			fn:   "",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "label", Arg1: "LOOP_START"},
			want: `(LOOP_START)
`,
		},
		{
			name: "translate label command inside function",
			fn:   "Test.sum",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "label", Arg1: "LOOP_START"},
			want: `(Test.sum$LOOP_START)
`,
		},
		{
			name: "translate goto command",
			fn:   "",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "goto", Arg1: "LOOP_START"},
			want: `@LOOP_START
0;JMP
`,
		},
		{
			name: "translate goto command inside function",
			fn:   "Test.sum",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "goto", Arg1: "LOOP_START"},
			want: `@Test.sum$LOOP_START
0;JMP
`,
		},
		{
			name: "translate if-goto command",
			fn:   "",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "if-goto", Arg1: "LOOP_START"},
			want: `@SP
M=M-1
A=M
D=M
@LOOP_START
D;JNE
`,
		},
		{
			name: "translate if-goto command inside function",
			fn:   "Test.sum",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "if-goto", Arg1: "LOOP_START"},
			want: `@SP
M=M-1
A=M
D=M
@Test.sum$LOOP_START
D;JNE
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.ExportSetFunctionName(tt.fn)
			got := tr.Translate(tt.arg)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestTranslateFunctionCmd(t *testing.T) {
	r1, r2 := rand.New(rand.NewSource(100)), rand.New(rand.NewSource(100))
	tr.ExportSetRandomizer(r1)
	ret := r2.Intn(1_000_000)
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		want string
	}{
		{
			name: "translate function command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Function, Command: "function", Arg1: "Test.sum", Arg2: 2},
			want: `(Test.sum)
@SP
A=M
M=0
@SP
M=M+1
@SP
A=M
M=0
@SP
M=M+1
`,
		},
		{
			name: "translate call command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Function, Command: "call", Arg1: "Test.sum", Arg2: 2},
			want: fmt.Sprintf(`@RETURN%[1]d
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@2
D=A
@5
D=A+D
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Test.sum
0;JMP
(RETURN%[1]d)
`, ret),
		},
		{
			name: "translate return command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Function, Command: "return"},
			want: `@LCL
D=M
@R13
M=D
@5
D=A
@R13
A=M-D
D=M
@R14
M=D
@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@1
D=A
@R13
A=M-D
D=M
@THAT
M=D
@2
D=A
@R13
A=M-D
D=M
@THIS
M=D
@3
D=A
@R13
A=M-D
D=M
@ARG
M=D
@4
D=A
@R13
A=M-D
D=M
@LCL
M=D
@R14
A=M
0;JMP
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.Translate(tt.arg)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

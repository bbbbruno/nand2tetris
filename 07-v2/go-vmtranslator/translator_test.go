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
	tr := vmtranslate.NewTranslator("Test")
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		want string
	}{
		{
			name: "translate push constant command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Segment: "constant", Index: 10},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Segment: "local", Index: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Segment: "pointer", Index: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Segment: "static", Index: 3},
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
	tr := vmtranslate.NewTranslator("Test")
	tests := []struct {
		name string
		arg  *vmtranslate.Cmd
		want string
	}{

		{
			name: "translate pop local command",
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Segment: "local", Index: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Segment: "pointer", Index: 1},
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
			arg:  &vmtranslate.Cmd{Type: vmtranslate.Pop, Command: "pop", Segment: "static", Index: 3},
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

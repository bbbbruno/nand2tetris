package assemble_test

import (
	"assemble"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAssembleNoSymbol(t *testing.T) {
	in := `// This is a test input

  @2
  D=A
  @3
  D=D+A  // Add D and A
  D;JEQ  // if D = 0 jump
`
	want := `0000000000000010
1110110000010000
0000000000000011
1110000010010000
1110001100000010
`

	r := bytes.NewBufferString(in)
	w := bytes.NewBufferString("")
	if err := assemble.Assemble(r, w); err != nil {
		t.Errorf("unexpected error occured: %v", err)
	}
	got := w.String()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
	}
}

func TestAssembleWithSymbol(t *testing.T) {
	in := `// sum from 1 to 100

    @i
    M=1 // i=1
    @sum
    M=0
(LOOP)
    @i
    D=M // D=i
    @100
    D=D-A // D=i-100
    @END
    D;JGT // if (i-100)>0 goto END
    @i
    D=M // D=i
    @sum
    M=D+M // sum=sum+1
    @i
    M=M+1 // i=i+1
    @LOOP
    0;JMP // goto LOOP
(END)
    @END
    0;JMP // infinity loop
`
	want := `0000000000010000
1110111111001000
0000000000010001
1110101010001000
0000000000010000
1111110000010000
0000000001100100
1110010011010000
0000000000010010
1110001100000001
0000000000010000
1111110000010000
0000000000010001
1111000010001000
0000000000010000
1111110111001000
0000000000000100
1110101010000111
0000000000010010
1110101010000111
`

	r := bytes.NewBufferString(in)
	w := bytes.NewBufferString("")
	if err := assemble.Assemble(r, w); err != nil {
		t.Errorf("unexpected error occured: %v", err)
	}
	got := w.String()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
	}
}

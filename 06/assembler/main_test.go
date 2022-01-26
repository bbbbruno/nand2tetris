package main

import (
	"bytes"
	"testing"
)

func TestAssemble(t *testing.T) {
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
	w := bytes.NewBuffer([]byte(""))
	Assemble(r, w)
	got := w.String()
	if got != want {
		t.Errorf("assemble(%#v) got %#v, want %#v", in, got, want)
	}
}

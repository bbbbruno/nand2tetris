package main

import (
	"bytes"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			`//this is test file

add
`,
			`// add
@SP
M=M-1
A=M
D=M
M=0
@SP
M=M-1
A=M
D=D+M
M=0
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, test := range tests {
		r := bytes.NewBufferString(test.in)
		w := bytes.NewBufferString("")
		err := convert(r, w)
		if err != nil {
			t.Errorf("error: got %#v", err)
		} else if got := w.String(); got != test.want {
			t.Errorf("test failed: expected %#v, got %#v", test.want, got)
		}
	}
}

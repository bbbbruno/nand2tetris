package vmtranslate_test

import (
	"testing"
	"vmtranslate"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   *vmtranslate.Cmd
		hasErr bool
	}{
		{
			name:   "parse add command",
			arg:    "add",
			want:   &vmtranslate.Cmd{Type: vmtranslate.Arithmetic, Command: "add"},
			hasErr: false,
		},
		{
			name:   "parse invalid arithmetic command",
			arg:    "div",
			want:   nil,
			hasErr: true,
		},
		{
			name:   "parse push constant command",
			arg:    "push constant 10",
			want:   &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "constant", Arg2: 10},
			hasErr: false,
		},
		{
			name:   "parse error on invalid pop constant command",
			arg:    "pop constant 10",
			want:   nil,
			hasErr: true,
		},
		{
			name:   "parse push command",
			arg:    "push local 1",
			want:   &vmtranslate.Cmd{Type: vmtranslate.Push, Command: "push", Arg1: "local", Arg2: 1},
			hasErr: false,
		},
		{
			name:   "parse error on invalid push command",
			arg:    "push special 1",
			want:   nil,
			hasErr: true,
		},
		{
			name:   "parse flow command",
			arg:    "label LOOP_START",
			want:   &vmtranslate.Cmd{Type: vmtranslate.Flow, Command: "label", Arg1: "LOOP_START"},
			hasErr: false,
		},
		{
			name:   "parse error on flow command with invalid symbol",
			arg:    "label LOOP_START!",
			want:   nil,
			hasErr: true,
		},
		{
			name:   "parse error on flow command with invalid index",
			arg:    "label LOOP_START 3",
			want:   nil,
			hasErr: true,
		},
		{
			name:   "parse error on invalid command",
			arg:    "reverse argument 5",
			want:   nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vmtranslate.Parse(tt.arg)
			if (err != nil) != tt.hasErr {
				t.Errorf("Parse() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

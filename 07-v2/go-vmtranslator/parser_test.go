package vmtranslate_test

import (
	"testing"
	"vmtranslate"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   vmtranslate.Cmd
		hasErr bool
	}{
		{
			name:   "parse add command",
			arg:    "add",
			want:   &vmtranslate.ArithmeticCmd{Command: "add"},
			hasErr: false,
		},
		{
			name:   "parse push/pop constant command",
			arg:    "push constant 10",
			want:   &vmtranslate.PushPopCmd{Command: "push", Segment: "constant", Index: 10},
			hasErr: false,
		},
		{
			name:   "parse error on invalid command",
			arg:    "reverse argument 5",
			want:   nil,
			hasErr: true,
		},
	}

	opt := cmpopts.IgnoreUnexported(vmtranslate.ArithmeticCmd{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vmtranslate.Parse(tt.arg)
			if (err != nil) != tt.hasErr {
				t.Errorf("Parse() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if !cmp.Equal(got, tt.want, opt) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

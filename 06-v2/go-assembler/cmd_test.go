package assemble_test

import (
	"assemble"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewACmd(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   assemble.Cmd
		hasErr bool
	}{
		{
			name:   "create A command with constant value",
			arg:    "@100",
			want:   &assemble.ACmd{Symbol: "", Value: 100},
			hasErr: false,
		},
		{
			name:   "create A command with valid symbol",
			arg:    "@i",
			want:   &assemble.ACmd{Symbol: "i", Value: 16},
			hasErr: false,
		},
		{
			name:   "create another A command with valid symbol",
			arg:    "@sum",
			want:   &assemble.ACmd{Symbol: "sum", Value: 17},
			hasErr: false,
		},
		{
			name:   "create A command with invalid symbol",
			arg:    "@sum!",
			want:   (*assemble.ACmd)(nil),
			hasErr: true,
		},
	}
	st := assemble.NewSymTable()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assemble.NewACmd(tt.arg, st)
			if (err != nil) != tt.hasErr {
				t.Errorf("Parse() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			fmt.Println(cmp.Equal(nil, nil))
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewCCmd(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   assemble.Cmd
		hasErr bool
	}{
		{
			name:   "create C command with only comp",
			arg:    "D+1",
			want:   &assemble.CCmd{Dest: "", Comp: "D+1", Jump: ""},
			hasErr: false,
		},
		{
			name:   "create C command with dest and comp",
			arg:    "M=M+1",
			want:   &assemble.CCmd{Dest: "M", Comp: "M+1", Jump: ""},
			hasErr: false,
		},
		{
			name:   "create C command with comp and jump",
			arg:    "D;JGT",
			want:   &assemble.CCmd{Dest: "", Comp: "D", Jump: "JGT"},
			hasErr: false,
		},
		{
			name:   "create C command with dest and comp and jump",
			arg:    "D=M;JMP",
			want:   &assemble.CCmd{Dest: "D", Comp: "M", Jump: "JMP"},
			hasErr: false,
		},
		{
			name:   "can't create C command",
			arg:    "DMT",
			want:   (*assemble.CCmd)(nil),
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assemble.NewCCmd(tt.arg)
			if (err != nil) != tt.hasErr {
				t.Errorf("Parse() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

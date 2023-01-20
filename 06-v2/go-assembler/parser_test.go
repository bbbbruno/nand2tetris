package assemble_test

import (
	"assemble"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseLables(t *testing.T) {
	type wants struct {
		lines    []string
		symtable *assemble.SymTable
	}
	tests := []struct {
		name   string
		in     string
		wants  *wants
		hasErr bool
	}{
		{
			name: "parse labels successfully",
			in: `// sum from 1 to 100

    @i
    M=1 // i=1
    @sum
    M=0
    (LOOP)
    @100
    D=D-A // D=i-100
    (END)
    @END
    0;JMP // infinity loop
`,
			wants: &wants{
				lines:    []string{"", "", "@i", "M=1", "@sum", "M=0", "(LOOP)", "@100", "D=D-A", "(END)", "@END", "0;JMP"},
				symtable: assemble.NewSymTable().AddSymbolWithAddr("LOOP", 4).AddSymbolWithAddr("END", 6),
			},
			hasErr: false,
		},
		{
			name: "invalid symbols for labels",
			in: `// sum from 1 to 100

    @i
    M=1 // i=1
    @sum
    M=0
    (LOOP!)
    @100
    D=D-A // D=i-100
    (001)
    @END
    0;JMP // infinity loop
`,
			wants:  &wants{nil, nil},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.in)
			lines, st, err := assemble.ParseLabels(r)
			if (err != nil) != tt.hasErr {
				t.Errorf("ParseLabels() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if diff := cmp.Diff(lines, tt.wants.lines); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
			if diff := cmp.Diff(st, tt.wants.symtable); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

func TestParseLines(t *testing.T) {
	type args struct {
		lines []string
		st    *assemble.SymTable
	}
	tests := []struct {
		name   string
		args   *args
		want   []assemble.Cmd
		hasErr bool
	}{
		{
			name: "parse lines successfully",
			args: &args{
				lines: []string{"", "", "@i", "M=1", "@sum", "M=0", "(LOOP)", "@100", "D=D-A", "(END)", "@END", "0;JMP"},
				st:    assemble.NewSymTable().AddSymbolWithAddr("LOOP", 4).AddSymbolWithAddr("END", 6),
			},
			want: []assemble.Cmd{
				&assemble.ACmd{Symbol: "i", Value: 16},
				&assemble.CCmd{Dest: "M", Comp: "1", Jump: ""},
				&assemble.ACmd{Symbol: "sum", Value: 17},
				&assemble.CCmd{Dest: "M", Comp: "0", Jump: ""},
				&assemble.ACmd{Symbol: "", Value: 100},
				&assemble.CCmd{Dest: "D", Comp: "D-A", Jump: ""},
				&assemble.ACmd{Symbol: "END", Value: 6},
				&assemble.CCmd{Dest: "", Comp: "0", Jump: "JMP"},
			},
			hasErr: false,
		},
		{
			name: "error while parsing lines",
			args: &args{
				lines: []string{"", "", "@0i", "M=100", "@sum!", "M=0", "(LOOP", "@100", "D=D-A", "(END)", "@END", "0;JMP"},
				st:    assemble.NewSymTable().AddSymbolWithAddr("LOOP", 5).AddSymbolWithAddr("END", 7),
			},
			want:   nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmds, err := assemble.ParseLines(tt.args.lines, tt.args.st)
			if (err != nil) != tt.hasErr {
				t.Errorf("ParseLines() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if diff := cmp.Diff(cmds, tt.want); diff != "" {
				t.Errorf("expected value is mismatch (-got +want):%s\n", diff)
			}
		})
	}
}

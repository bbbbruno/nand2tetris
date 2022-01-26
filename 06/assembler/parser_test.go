package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	p := NewParser()
	text := `// This is a test input

  @i // i=1
  M=1
  @sum // sum=0
  M=0
(LOOP)
  @i
  D=M
  @100
  D=D-A
  @LOOP
  D;JGT
`
	input := bytes.NewReader([]byte(text))
	cmds := []*Command{
		{A_COMMAND, "i", "", "", ""},
		{C_COMMAND, "", "M", "1", ""},
		{A_COMMAND, "sum", "", "", ""},
		{C_COMMAND, "", "M", "0", ""},
		{A_COMMAND, "i", "", "", ""},
		{C_COMMAND, "", "D", "M", ""},
		{A_COMMAND, "100", "", "", ""},
		{C_COMMAND, "", "D", "D-A", ""},
		{A_COMMAND, "LOOP", "", "", ""},
		{C_COMMAND, "", "", "D", "JGT"},
	}
	st := NewSymbolTableWithOpts(18, map[symbol]addr{"i": 16, "sum": 17, "LOOP": 4})
	gotCmds, gotSt := p.Parse(input)
	if !reflect.DeepEqual(gotCmds, cmds) {
		t.Errorf("p.Parse([]byte(%#v)) got %#v, want %#v", text, gotCmds, cmds)
	}
	if !reflect.DeepEqual(gotSt, st) {
		t.Errorf("p.Parse([]byte(%#v)) got %#v, want %#v", text, gotSt, st)
	}
}

func TestParseInput(t *testing.T) {
	text := `// This is a test input

  @100
  D=A
  @i
(LOOP)
  D=D+A  // Add D and A
  D;JEQ  // if D = 0 jump
`
	input := bytes.NewReader([]byte(text))
	want := []string{"@100", "D=A", "@i", "(LOOP)", "D=D+A", "D;JEQ"}
	got := parseInput(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseInput([]byte(%#v)) got %#v, want %#v", text, got, want)
	}
}

func TestParseSymbols(t *testing.T) {
	in := []string{"@100", "D=A", "@i", "(LOOP)", "D=D+A", "D;JEQ"}
	want := NewSymbolTableWithOpts(16, map[symbol]addr{"LOOP": 3})
	got := parseSymbols(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseSymbols(%#v) got %#v, want %#v", in, got, want)
	}
}

func TestParseLines(t *testing.T) {
	lines := []string{"@100", "D=A", "@i", "(LOOP)", "D=D+A", "D;JEQ"}
	st := NewSymbolTableWithOpts(16, map[symbol]addr{"LOOP": 3})
	wantCmds := []*Command{
		{A_COMMAND, "100", "", "", ""},
		{C_COMMAND, "", "D", "A", ""},
		{A_COMMAND, "i", "", "", ""},
		{C_COMMAND, "", "D", "D+A", ""},
		{C_COMMAND, "", "", "D", "JEQ"},
	}
	wantSt := NewSymbolTableWithOpts(17, map[symbol]addr{"LOOP": 3, "i": 16})
	gotCmds, gotSt := parseLines(lines, st)
	if !reflect.DeepEqual(gotCmds, wantCmds) {
		t.Errorf("parseLines(%#v, %#v) got %#v, want %#v", lines, st, gotCmds, wantCmds)
	}
	if !reflect.DeepEqual(gotSt, wantSt) {
		t.Errorf("parseLines(%#v, %#v) got %#v, want %#v", lines, st, gotSt, wantSt)
	}
}

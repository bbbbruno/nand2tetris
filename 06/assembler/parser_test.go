package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	p := NewParser()
	text := `// This is a test input

  @2
  D=A
  @3
  D=D+A  // Add D and A
  D;JEQ  // if D = 0 jump

`
	input := bytes.NewReader([]byte(text))
	want := []*Command{
		{A_COMMAND, "2", "", "", ""},
		{C_COMMAND, "", "D", "A", ""},
		{A_COMMAND, "3", "", "", ""},
		{C_COMMAND, "", "D", "D+A", ""},
		{C_COMMAND, "", "", "D", "JEQ"},
	}
	got := p.Parse(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("p.Parse([]byte(%#v)) got %#v, want %#v", text, got, want)
	}
}

func TestParseInput(t *testing.T) {
	text := `// This is a test input

  @2
  D=A
  @3
  D=D+A  // Add D and A
  D;JEQ  // if D = 0 jump

`
	input := bytes.NewReader([]byte(text))
	want := []string{"@2", "D=A", "@3", "D=D+A", "D;JEQ"}
	got := parseInput(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseInput([]byte(%#v)) got %#v, want %#v", text, got, want)
	}
}

func TestParseLine(t *testing.T) {
	testCases := []struct {
		in   string
		want *Command
	}{
		{in: "@2", want: &Command{A_COMMAND, "2", "", "", ""}},
		{in: "D=A", want: &Command{C_COMMAND, "", "D", "A", ""}},
		{in: "@3", want: &Command{A_COMMAND, "3", "", "", ""}},
		{in: "D=D+A", want: &Command{C_COMMAND, "", "D", "D+A", ""}},
		{in: "D;JEQ", want: &Command{C_COMMAND, "", "", "D", "JEQ"}},
	}
	for _, test := range testCases {
		got := parseLine(test.in)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseLine(%#v) got %#v, want %#v", test.in, got, test.want)
		}
	}
}

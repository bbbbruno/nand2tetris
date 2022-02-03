package tokenizer

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	in := `// (identical to projects/09/Average/Main.jack)

/** Computes the average of a sequence of integers. */
class Main {
    function void main() {
        var int i;
	
	let length = Keyboard.readInt("HOW MANY NUMBERS? ");
	let i = 0;`
	want := []*token{
		{KEYWORD, "class"},
		{IDENTIFIER, "Main"},
		{SYMBOL, "{"},
		{KEYWORD, "function"},
		{KEYWORD, "void"},
		{IDENTIFIER, "main"},
		{SYMBOL, "("},
		{SYMBOL, ")"},
		{SYMBOL, "{"},
		{KEYWORD, "var"},
		{KEYWORD, "int"},
		{IDENTIFIER, "i"},
		{SYMBOL, ";"},
		{KEYWORD, "let"},
		{IDENTIFIER, "length"},
		{SYMBOL, "="},
		{IDENTIFIER, "Keyboard"},
		{SYMBOL, "."},
		{IDENTIFIER, "readInt"},
		{SYMBOL, "("},
		{STRING_CONST, "HOW MANY NUMBERS? "},
		{SYMBOL, ")"},
		{SYMBOL, ";"},
		{KEYWORD, "let"},
		{IDENTIFIER, "i"},
		{SYMBOL, "="},
		{INT_CONST, "0"},
		{SYMBOL, ";"},
	}
	r := bufio.NewReader(bytes.NewBufferString(in))
	tkz := &tokenizer{r, nil, nil, make([]byte, 0)}
	for i := 0; tkz.HasMoreTokens(); i++ {
		tkz.Advance()
		if !reflect.DeepEqual(tkz.currentToken, want[i]) {
			t.Errorf("FAILED: expected %#v, got %#v", want[i], tkz.currentToken)
		}
	}
}

func TestSkipComments(t *testing.T) {
	tests := []struct {
		in       string
		nextByte byte
	}{
		{in: "// initializes sum with 0\nlet sum = 0", nextByte: byte('l')},
		{in: "/* initializes sum with 0 */let sum = 0", nextByte: byte('l')},
		{in: `/** initializes
sum
with 0 */let sum = 0`, nextByte: byte('l')},
	}

	for _, test := range tests {
		r := bufio.NewReader(bytes.NewBufferString(test.in))
		tkr := &tokenizer{r, nil, nil, make([]byte, 0)}
		_, _ = tkr.Reader.ReadByte()
		b, _ := tkr.Reader.ReadByte()
		err := tkr.skipComments(b)
		nextByte, _ := tkr.Reader.ReadByte()
		if err != nil {
			t.Errorf("ERROR: got %#v", err)
		} else if nextByte != test.nextByte {
			t.Errorf("FAILED: expected %#v, got %#v", string(test.nextByte), string(nextByte))
		}
	}
}

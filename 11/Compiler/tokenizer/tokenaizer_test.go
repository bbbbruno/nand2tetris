package tokenizer

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func newTokenizer(in string) *tokenizer {
	r := bufio.NewReader(bytes.NewBufferString(in))
	return &tokenizer{r, nil, nil, make([]byte, 0)}
}

// Advance()とPeek()を実行したときのCurrentToken()の状態をテストする。
func TestCurrentToken(t *testing.T) {
	in := `class Test {
  static boolean test; // test file

  function void test() {
    return;
  }
}`
	tkz := newTokenizer(in)
	tkz.Advance()
	tkz.Advance()
	tkz.Advance()
	if token := tkz.CurrentToken(); !token.IsSymbol("{") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "{", token.Content())
	}
	tkz.Advance()
	// Peek()を実行すると次のトークンを先読みして返し、CurrentToken()は変わらない。
	if token := tkz.Peek(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	} else if token := tkz.CurrentToken(); !token.IsKeyword("static") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "static", token.Content())
	}
	// Peek()を何回実行してもそれ以上先へは進まない。
	if token := tkz.Peek(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	} else if token := tkz.CurrentToken(); !token.IsKeyword("static") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "static", token.Content())
	}
	// Advance()を実行するとcurrentToken()が先へ進む。
	tkz.Advance()
	if token := tkz.CurrentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	} else if token := tkz.CurrentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	}
}

func TestParse(t *testing.T) {
	in := `// (identical to projects/09/Average/Main.jack)

/** Computes the average of a sequence of integers. */
class Main {
    function void main() {
        var int i;
	
	let length = Keyboard.readInt("HOW MANY NUMBERS? ");
	let i = 0;`
	want := []*Token{
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
	tkz := newTokenizer(in)
	for i := 0; i < len(want)-1; i++ {
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
	* sum
	* with 0
	*/let sum = 0`, nextByte: byte('l')},
	}

	for _, test := range tests {
		tkz := newTokenizer(test.in)
		_, _ = tkz.r.ReadByte()
		b, _ := tkz.r.ReadByte()
		err := tkz.skipComments(b)
		nextByte, _ := tkz.r.ReadByte()
		if err != nil {
			t.Errorf("ERROR: got %#v", err)
		} else if nextByte != test.nextByte {
			t.Errorf("FAILED: expected %#v, got %#v", string(test.nextByte), string(nextByte))
		}
	}
}

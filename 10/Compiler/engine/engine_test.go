package engine

import (
	"bufio"
	"bytes"
	"compiler/tokenizer"
	"testing"
)

var defaultBuf = bytes.NewBuffer(nil)

func newEngine() *engine {
	defaultBuf.Reset()
	r := bytes.NewBufferString(`class Test {
  static boolean test; // test file

  function void test() {
    return;
  }
}`)
	tkz := tokenizer.New(r)
	b := bufio.NewWriter(defaultBuf)
	return &engine{tkz, nil, b, 0}
}

func TestCompile(t *testing.T) {
	b := defaultBuf
	e := newEngine()
	want := `<class>
  <keyword> class </keyword>
  <identifier> Test </identifier>
  <symbol> { </symbol>
  <classVarDec>
    <keyword> static </keyword>
    <keyword> boolean </keyword>
    <identifier> test </identifier>
    <symbol> ; </symbol>
  </classVarDec>
  <subroutineDec>
    <keyword> function </keyword>
    <keyword> void </keyword>
    <identifier> test </identifier>
    <symbol> ( </symbol>
    <parameterList>
    </parameterList>
    <symbol> ) </symbol>
    <subroutineBody>
      <symbol> { </symbol>
      <statements>
        <returnStatement>
          <keyword> return </keyword>
          <symbol> ; </symbol>
        </returnStatement>
      </statements>
      <symbol> } </symbol>
    </subroutineBody>
  </subroutineDec>
  <symbol> } </symbol>
</class>
`
	if err := e.Compile(); err != nil {
		t.Errorf("ERROR: expected no error, got %#v", err)
	} else if b.String() != want {
		t.Errorf("FAILED: expected %#v, got %#v", want, b.String())
	}
}

// advanceToken()とpeekToken()を実行したときのcurrentToken()の状態をテストする。
func TestCurrentToken(t *testing.T) {
	e := newEngine()
	e.advanceToken()
	e.advanceToken()
	e.advanceToken()
	if token := e.currentToken(); !token.IsSymbol("{") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "{", token.Content())
	}
	// peekToken()を実行すると、内部のtokenizerのトークンだけ先に進む。
	e.advanceToken()
	e.peekToken()
	if token := e.currentToken(); !token.IsKeyword("static") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "static", token.Content())
	} else if token := e.tkz.CurrentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	}
	// peekToken()を何回実行してもそれ以上先へは進まない。
	e.peekToken()
	if token := e.currentToken(); !token.IsKeyword("static") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "static", token.Content())
	} else if token := e.tkz.CurrentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	}
	// acvanceToken()を実行するとcurrentToken()が先へ進み内部のtokenizerのトークンに追いつく。
	e.advanceToken()
	if token := e.currentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	} else if token := e.tkz.CurrentToken(); !token.IsKeyword("boolean") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "boolean", token.Content())
	}
}

func TestAdvanceToken(t *testing.T) {
	e := newEngine()
	// e.tokenがnilの場合はトークンが先に進む。
	e.advanceToken()
	if token := e.currentToken(); !token.IsKeyword("class") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "class", token.Content())
	} else if e.token != nil {
		t.Errorf("FAILED: expected engine's token to be %#v, got %#v", nil, e.token)
	}
	// e.tokenに中身がある場合は中身を空にしてトークンは先へは進まない。
	e.token = e.currentToken()
	e.advanceToken()
	if token := e.currentToken(); !token.IsKeyword("class") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "class", token.Content())
	} else if e.token != nil {
		t.Errorf("FAILED: expected engine's token to be %#v, got %#v", nil, e.token)
	}
}

func TestPeekToken(t *testing.T) {
	e := newEngine()
	e.advanceToken()
	// e.tokenがnilの場合は現在のトークンをe.tokenにセットして、tokenizerのトークンだけ先へ進める。
	e.peekToken()
	if token := e.currentToken(); !token.IsKeyword("class") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "class", token.Content())
	} else if e.token == nil {
		t.Errorf("FAILED: expected engine's token to be %#v, got %#v", "class", nil)
	} else if token := e.tkz.CurrentToken(); !token.IsIdentifier() {
		t.Errorf("FAILED: expected tokenizer's token to be %#v, got %#v", "Test", token.Content())
	}
	// e.tokenに中身がある場合は何回e.peekToken()を実行しても何も変化しない。
	e.peekToken()
	if token := e.currentToken(); !token.IsKeyword("class") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "class", token.Content())
	} else if e.token == nil {
		t.Errorf("FAILED: expected engine's token to be %#v, got %#v", "class", nil)
	} else if token := e.tkz.CurrentToken(); !token.IsIdentifier() {
		t.Errorf("FAILED: expected tokenizer's token to be %#v, got %#v", "Test", token.Content())
	}
	e.peekToken()
	if token := e.currentToken(); !token.IsKeyword("class") {
		t.Errorf("FAILED: expected current token to be %#v, got %#v", "class", token.Content())
	} else if e.token == nil {
		t.Errorf("FAILED: expected engine's token to be %#v, got %#v", "class", nil)
	} else if token := e.tkz.CurrentToken(); !token.IsIdentifier() {
		t.Errorf("FAILED: expected tokenizer's token to be %#v, got %#v", "Test", token.Content())
	}
}

package engine

import (
	"bufio"
	"bytes"
	"testing"

	"jackcompiler/tokenizer"
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
	return New(tkz, b)
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

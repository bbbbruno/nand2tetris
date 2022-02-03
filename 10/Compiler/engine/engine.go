package engine

import (
	"bufio"
	"compiler/tokenizer"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Engine interface {
	Compile() error
}

type engine struct {
	tokenizer.Tokenizer
	token tokenizer.Token
	*bufio.Writer
	hierarchy int
}

func New(tkz tokenizer.Tokenizer, w io.Writer) Engine {
	return &engine{tkz, nil, bufio.NewWriter(w), 0}
}

func (e *engine) Compile() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
	e.compileClass()

	return err
}

const (
	OPEN  = true
	CLOSE = false
)

func (e *engine) compileClass() {
	e.writeHierarchy("class", OPEN)
	e.writeKeyword("class")
	e.writeIdentifier()
	e.writeSymbol("{")
	e.compileClassVarDec()
	e.compileSubroutineDec()
	e.writeSymbol("}")
}

func (e *engine) compileClassVarDec() {
}

func (e *engine) compileSubroutineDec() {
}

func (e *engine) advanceToken() tokenizer.Token {
	if token := e.token; token != nil {
		e.token = nil
		return token
	}
	if e.HasMoreTokens() {
		e.Advance()
		return e.CurrentToken()
	} else {
		panic(errors.New("no more tokens"))
	}
}

func (e *engine) writeIdentifier() {
	token := e.advanceToken()
	if token.Type() == tokenizer.IDENTIFIER {
		e.writeToken()
	} else {
		panic(errors.New("token is not identifier"))
	}
}

func (e *engine) writeKeyword(keys ...string) {
	token := e.advanceToken()
	if token.Type() != tokenizer.KEYWORD {
		panic(errors.New("token is not keyword"))
	}
	if !contains(keys, token.Content()) {
		panic(errors.New("token is not desired keyword"))
	}

	e.writeToken()
}

func (e *engine) writeSymbol(sym string) {
	token := e.advanceToken()
	if token.Type() == tokenizer.SYMBOL && token.Content() == sym {
		e.writeToken()
	} else {
		panic(errors.New("token is not symbol"))
	}
}

func (e *engine) writeHierarchy(s string, open bool) {
	if open {
		e.hierarchy++
	} else {
		e.hierarchy--
	}
	spaces := strings.Repeat(" ", e.hierarchy*2)
	tag := "<" + s + ">"
	if !open {
		tag = tag[:1] + "/" + tag[1:]
	}
	text := fmt.Sprintf("%s%s", spaces, tag)
	if _, err := fmt.Fprintln(e, text); err != nil {
		panic(err)
	}
}

func (e *engine) writeToken() {
	spaces := strings.Repeat(" ", e.hierarchy*2)
	text := fmt.Sprintf("%s%s", spaces, e.CurrentToken())
	if _, err := fmt.Fprintln(e, text); err != nil {
		panic(err)
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

package engine

import (
	"errors"
	"fmt"
	"strings"
)

func (e *engine) writeKeyword(keys ...string) {
	if token := e.currentToken(); token.IsKeyword(keys...) {
		e.writeToken()
	} else {
		panic(fmt.Errorf("token is not expected keyword, expected %v", keys))
	}
}

func (e *engine) writeSymbol(syms ...string) {
	if token := e.currentToken(); token.IsSymbol(syms...) {
		e.writeToken()
	} else {
		panic(fmt.Errorf("token is not expected symbol, expected %v", syms))
	}
}

func (e *engine) writeIdentifier() {
	if token := e.currentToken(); token.IsIdentifier() {
		e.writeToken()
	} else {
		panic(errors.New("token is not identifier"))
	}
}

func (e *engine) writeIntConst() {
	if token := e.currentToken(); token.IsIdentifier() {
		e.writeToken()
	} else {
		panic(errors.New("token is not integer constant"))
	}
}

func (e *engine) writeStringConst() {
	if token := e.currentToken(); token.IsIdentifier() {
		e.writeToken()
	} else {
		panic(errors.New("token is not string constant"))
	}
}

var primitiveTypes = []string{"int", "char", "boolean"}

func (e *engine) writeType() {
	if token := e.currentToken(); token.IsKeyword() {
		e.writeKeyword(primitiveTypes...)
	} else if token.IsIdentifier() {
		e.writeIdentifier()
	}
}

func (e *engine) writeKeywordOrType(keys ...string) {
	if token := e.currentToken(); token.IsKeyword(keys...) {
		e.writeKeyword(keys...)
	} else {
		e.writeType()
	}
}

func (e *engine) writeHierarchy(s string, open bool) {
	if !open {
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
	if open {
		e.hierarchy++
	}
}

func (e *engine) writeToken() {
	spaces := strings.Repeat(" ", e.hierarchy*2)
	text := fmt.Sprintf("%s%s", spaces, e.currentToken())
	if _, err := fmt.Fprintln(e, text); err != nil {
		panic(err)
	}
	e.advanceToken()
}

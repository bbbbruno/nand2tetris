package engine

import (
	"errors"
	"fmt"
	"strings"
)

func (e *engine) writeKeyword(keys ...string) {
	if token := e.CurrentToken(); token.IsKeyword(keys...) {
		if s := token.Content(); s == "static" || s == "field" || s == "var" {
			e.varBuf[0] = s
		}
		e.writeToken()
	} else {
		panic(fmt.Errorf("token is not expected keyword, expected %v", keys))
	}
}

func (e *engine) writeSymbol(syms ...string) {
	if token := e.CurrentToken(); token.IsSymbol(syms...) {
		e.writeToken()
	} else {
		panic(fmt.Errorf("token is not expected symbol, expected %v", syms))
	}
}

func (e *engine) writeIdentifier() {
	if token := e.CurrentToken(); token.IsIdentifier() {
		e.writeToken()
	} else {
		panic(errors.New("token is not identifier"))
	}
}

func (e *engine) writeIntConst() {
	if token := e.CurrentToken(); token.IsIntConst() {
		e.writeToken()
	} else {
		panic(errors.New("token is not integer constant"))
	}
}

func (e *engine) writeStringConst() {
	if token := e.CurrentToken(); token.IsStringConst() {
		e.writeToken()
	} else {
		panic(errors.New("token is not string constant"))
	}
}

var primitiveTypes = []string{"int", "char", "boolean"}

func (e *engine) writeType() {
	if token := e.CurrentToken(); token.IsKeyword() {
		e.varBuf[1] = token.Content()
		e.writeKeyword(primitiveTypes...)
	} else if token.IsIdentifier() {
		e.varBuf[1] = token.Content()
		e.writeIdentifier()
	}
}

func (e *engine) writeKeywordOrType(keys ...string) {
	if token := e.CurrentToken(); token.IsKeyword(keys...) {
		e.writeKeyword(keys...)
	} else {
		e.writeType()
	}
}

func (e *engine) writeVarName() {
	e.varBuf[2] = e.CurrentToken().Content()
	e.writeIdentifier()
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
	text := fmt.Sprintf("%s%s", spaces, e.CurrentToken())
	if _, err := fmt.Fprintln(e, text); err != nil {
		panic(err)
	}
	e.Advance()
}

func (e *engine) writeVar() {
	name, symtype, kind := e.varBuf[2], e.varBuf[1], e.varBuf[0]
	var (
		sym Symbol
		err error
	)
	if e.scope == CLASS {
		sym, err = e.ClassTable().Define(name, symtype, kind)
	} else {
		sym, err = e.SubroutineTable().Define(name, symtype, kind)
	}
	if err != nil {
		panic(err)
	}

	spaces := strings.Repeat(" ", e.hierarchy*2)
	text := fmt.Sprintf("%s%s", spaces, sym)
	if _, err := fmt.Fprintln(e, text); err != nil {
		panic(err)
	}
}
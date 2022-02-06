package engine

import (
	"errors"
	"fmt"
)

func (e *engine) validateKeyword(keys ...string) {
	if token := e.CurrentToken(); token.IsKeyword(keys...) {
		if s := token.Content(); s == "static" || s == "field" || s == "var" {
			e.variable.kind = s
		} else if s == "constructor" || s == "function" || s == "method" {
			e.subroutine.kind = s
		}
		e.Advance()
	} else {
		panic(fmt.Errorf("token is not expected keyword, expected %v", keys))
	}
}

func (e *engine) validateSymbol(syms ...string) {
	if token := e.CurrentToken(); token.IsSymbol(syms...) {
		e.Advance()
	} else {
		panic(fmt.Errorf("token is not expected symbol, expected %v", syms))
	}
}

func (e *engine) validateIdentifier() {
	if token := e.CurrentToken(); token.IsIdentifier() {
		e.Advance()
	} else {
		panic(errors.New("token is not identifier"))
	}
}

func (e *engine) validateIntConst() {
	if token := e.CurrentToken(); token.IsIntConst() {
		e.term.kind, e.term.value = "intConst", token.Content()
		e.Advance()
	} else {
		panic(errors.New("token is not integer constant"))
	}
}

func (e *engine) validateStringConst() {
	if token := e.CurrentToken(); token.IsStringConst() {
		e.Advance()
	} else {
		panic(errors.New("token is not string constant"))
	}
}

var primitiveTypes = []string{"int", "char", "boolean"}

func (e *engine) validateType() {
	token := e.CurrentToken()
	e.variable.symtype = token.Content()
	if token.IsKeyword() {
		e.validateKeyword(primitiveTypes...)
	} else if token.IsIdentifier() {
		e.validateIdentifier()
	}
}

func (e *engine) validateKeywordOrType(keys ...string) {
	token := e.CurrentToken()
	e.subroutine.returnval = token.Content()
	if token.IsKeyword(keys...) {
		e.validateKeyword(keys...)
	} else {
		e.validateType()
	}
}

func (e *engine) validateVarName() {
	e.variable.name = e.CurrentToken().Content()
	e.addVar()
	e.validateIdentifier()
}

func (e *engine) validateSubroutineName() {
	e.subroutine.name = e.CurrentToken().Content()
	e.validateIdentifier()
}

func (e *engine) validateReceiverName() {
	e.subroutine.receiver = e.subroutine.name
	e.subroutine.name = e.CurrentToken().Content()
	e.validateIdentifier()
}

func (e *engine) validateOperator(keys ...string) {
	e.expression.nextoperator = e.expression.operator
	e.expression.operator = e.CurrentToken().Content()
	e.validateSymbol(keys...)
}

func (e *engine) addVar() {
	var (
		// sym Symbol
		err error
	)
	if e.scope == CLASS {
		_, err = e.ClassTable().Define(e.variable.name, e.variable.symtype, e.variable.kind)
	} else {
		_, err = e.SubroutineTable().Define(e.variable.name, e.variable.symtype, e.variable.kind)
	}
	if err != nil {
		panic(err)
	}
}

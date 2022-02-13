package engine

import (
	"errors"
	"fmt"
	"jackcompiler/symtable"
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
		e.term.value = token.Content()
		e.Advance()
	} else {
		panic(errors.New("token is not integer constant"))
	}
}

func (e *engine) validateStringConst() {
	if token := e.CurrentToken(); token.IsStringConst() {
		e.term.value = token.Content()
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

func (e *engine) validateFunctionName() {
	e.functions = append(e.functions, &function{"", e.CurrentToken().Content(), 0})
	e.validateIdentifier()
}

func (e *engine) validateReceiverName() {
	function := e.functions[len(e.functions)-1]
	function.receiver = function.name
	function.name = e.CurrentToken().Content()
	e.validateIdentifier()
}

func (e *engine) validateOperator(keys ...string) {
	e.operators = append(e.operators, e.CurrentToken().Content())
	e.validateSymbol(keys...)
}

func (e *engine) addVar() {
	if sym, ok := e.SubroutineTable().Find(e.variable.name); ok {
		e.sym = sym
		return
	} else if sym, ok := e.ClassTable().Find(e.variable.name); ok {
		e.sym = sym
		return
	}

	var table symtable.Table
	if e.scope == CLASS {
		table = e.ClassTable()
	} else {
		table = e.SubroutineTable()
	}
	if err := table.Define(e.variable.name, e.variable.symtype, e.variable.kind); err != nil {
		panic(err)
	}
	sym, _ := table.Find(e.variable.name)
	e.sym = sym
}

package engine

import (
	"strconv"
)

func (e *engine) declareSubroutine() {
	receiver, name, nLocals := e.className, e.subroutine.name, e.SubroutineTable().VarCount("argument")
	if err := e.writeFunction(receiver, name, nLocals); err != nil {
		panic(err)
	}
}

func (e *engine) doStatment() {
	receiver, name, nArgs := e.subroutine.receiver, e.subroutine.name, e.expressionCount
	e.expressionCount = 0
	if err := e.writeCall(receiver, name, nArgs); err != nil {
		panic(err)
	}
}

func (e *engine) returnStatment() {
	if e.subroutine.returnval == "void" {
		if err := e.writePop("temp", 0); err != nil {
			panic(err)
		}
		if err := e.writePush("constant", 0); err != nil {
			panic(err)
		}
	}
	if err := e.writeReturn(); err != nil {
		panic(err)
	}
}

var commands = map[string]string{
	"+": "add",
	"-": "sub",
	"*": "call Math.multiply 2",
	"/": "call Math.divide 2",
	"=": "eq",
	">": "gt",
	"<": "lt",
}

func (e *engine) calcExpression() {
	cmd := commands[e.expression.operator]
	e.expression.operator = e.expression.nextoperator
	if err := e.writeArithmethic(cmd); err != nil {
		panic(err)
	}
}

func (e *engine) intTerm() {
	i, err := strconv.Atoi(e.term.value)
	if err != nil {
		panic(err)
	}

	segment, index := "constant", i
	if err := e.writePush(segment, index); err != nil {
		panic(err)
	}
}

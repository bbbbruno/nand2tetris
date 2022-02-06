package engine

import (
	"strconv"
)

func (e *engine) declareSubroutine() {
	receiver, name, nLocals := e.className, e.subroutine.name, e.SubroutineTable().VarCount("var")
	if err := e.writeFunction(receiver, name, nLocals); err != nil {
		panic(err)
	}
}

func (e *engine) startWhile(count int) {
	e.whileCount++
	label := "WHILE_EXP" + strconv.Itoa(count)
	if err := e.writeLabel(label); err != nil {
		panic(err)
	}
}

func (e *engine) jumpWhile(count int) {
	label := "WHILE_END" + strconv.Itoa(count)
	if err := e.writeArithmethic("not"); err != nil {
		panic(err)
	}
	if err := e.writeIf(label); err != nil {
		panic(err)
	}
}

func (e *engine) endWhile(count int) {
	expLabel := "WHILE_EXP" + strconv.Itoa(count)
	endLabel := "WHILE_END" + strconv.Itoa(count)
	if err := e.writeGoto(expLabel); err != nil {
		panic(err)
	}
	if err := e.writeLabel(endLabel); err != nil {
		panic(err)
	}
}

func (e *engine) startIf(count int) {
	trueLabel := "IF_TRUE" + strconv.Itoa(count)
	falseLabel := "IF_FALSE" + strconv.Itoa(count)
	if err := e.writeIf(trueLabel); err != nil {
		panic(err)
	}
	if err := e.writeGoto(falseLabel); err != nil {
		panic(err)
	}
	if err := e.writeLabel(trueLabel); err != nil {
		panic(err)
	}
}

func (e *engine) elseIf(end bool, count int) {
	falseLabel := "IF_FALSE" + strconv.Itoa(count)
	endLabel := "IF_END" + strconv.Itoa(count)
	if !end {
		if err := e.writeGoto(endLabel); err != nil {
			panic(err)
		}
	}
	if err := e.writeLabel(falseLabel); err != nil {
		panic(err)
	}
}

func (e *engine) endIf(count int) {
	endLabel := "IF_END" + strconv.Itoa(count)
	if err := e.writeLabel(endLabel); err != nil {
		panic(err)
	}
}

func (e *engine) letStatement() {
	segment, index := e.sym.Kind.String(), e.sym.Index
	if err := e.writePop(segment, index); err != nil {
		panic(err)
	}
}

func (e *engine) doStatement() {
	receiver, name, nArgs := e.subroutine.receiver, e.subroutine.name, e.expressionCount
	e.expressionCount = 0
	if err := e.writeCall(receiver, name, nArgs); err != nil {
		panic(err)
	}
	if err := e.writePop("temp", 0); err != nil {
		panic(err)
	}
}

func (e *engine) returnStatement() {
	if e.subroutine.returnval == "void" {
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
	"&": "and",
	"|": "or",
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

var unaries = map[string]string{
	"-": "neg",
	"~": "not",
}

func (e *engine) calcUnary() {
	unary := unaries[e.expression.operator]
	e.expression.operator = e.expression.nextoperator
	if err := e.writeArithmethic(unary); err != nil {
		panic(err)
	}
}

func (e *engine) callIntConst() {
	i, err := strconv.Atoi(e.term.value)
	if err != nil {
		panic(err)
	}

	segment, index := "constant", i
	if err := e.writePush(segment, index); err != nil {
		panic(err)
	}
}

func (e *engine) callKeywordConst() {
	segment := "constant"
	if v := e.term.value; v == "true" {
		if err := e.writePush(segment, 0); err != nil {
			panic(err)
		}
		if err := e.writeArithmethic("not"); err != nil {
			panic(err)
		}
	} else if v == "false" {
		if err := e.writePush(segment, 0); err != nil {
			panic(err)
		}
	}
}

func (e *engine) callFunc() {
	receiver, name, nArgs := e.subroutine.receiver, e.subroutine.name, e.expressionCount
	e.expressionCount = 0
	if err := e.writeCall(receiver, name, nArgs); err != nil {
		panic(err)
	}
}

func (e *engine) callVar() {
	segment, index := e.sym.Kind.String(), e.sym.Index
	if err := e.writePush(segment, index); err != nil {
		panic(err)
	}
}

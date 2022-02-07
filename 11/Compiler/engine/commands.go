package engine

import (
	"jackcompiler/symtable"
	"strconv"
)

func (e *engine) declareSubroutine() {
	receiver, name, nLocals := e.className, e.subroutine.name, e.SubroutineTable().VarCount("var")
	if err := e.writeFunction(receiver, name, nLocals); err != nil {
		panic(err)
	}
}

func (e *engine) startWhile(count int) {
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

func (e *engine) letStatement(sym *symtable.Symbol) {
	segment, index := sym.Kind.String(), sym.Index
	if err := e.writePop(segment, index); err != nil {
		panic(err)
	}
}

func (e *engine) letArrayStatment() {
	if err := e.writePop("temp", 0); err != nil {
		panic(err)
	}
	if err := e.writePop("pointer", 1); err != nil {
		panic(err)
	}
	if err := e.writePush("temp", 0); err != nil {
		panic(err)
	}
	if err := e.writePop("that", 0); err != nil {
		panic(err)
	}
}

func (e *engine) doStatement() {
	e.callFunc()
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

func (e *engine) calcArray(sym *symtable.Symbol) {
	e.sym = sym
	e.callVar()
	e.expression.nextoperator = e.expression.operator
	e.expression.operator = "+"
	e.calcExpression()
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

func (e *engine) callStringConst() {
	str := e.term.value
	if err := e.writePush("constant", len(str)); err != nil {
		panic(err)
	}
	if err := e.writeCall("String", "new", 1); err != nil {
		panic(err)
	}
	for _, s := range e.term.value {
		if err := e.writePush("constant", int(s)); err != nil {
			panic(err)
		}
		if err := e.writeCall("String", "appendChar", 2); err != nil {
			panic(err)
		}
	}
}

func (e *engine) callKeywordConst() {
	switch e.term.value {
	case "true":
		if err := e.writePush("constant", 0); err != nil {
			panic(err)
		}
		if err := e.writeArithmethic("not"); err != nil {
			panic(err)
		}
	case "false":
		if err := e.writePush("constant", 0); err != nil {
			panic(err)
		}
	case "this":
		if err := e.writePush("pointer", 0); err != nil {
			panic(err)
		}
	case "null":
		if err := e.writePush("constant", 0); err != nil {
			panic(err)
		}
	}
}

func (e *engine) resetSubroutine() {
	e.subroutine.receiver, e.subroutine.name, e.subroutine.kind = "", "", ""
	e.expressionCount = 0
}

func (e *engine) callFunc() {
	receiver, name, nArgs := e.subroutine.receiver, e.subroutine.name, e.expressionCount
	if receiver == "" {
		nArgs++
		receiver = e.className
		if err := e.writePush("pointer", 0); err != nil {
			panic(err)
		}
	} else if sym, ok := e.SubroutineTable().Find(receiver); ok {
		nArgs++
		receiver = sym.Symtype
		if err := e.writePush(sym.Kind.String(), sym.Index); err != nil {
			panic(err)
		}
	} else if sym, ok := e.ClassTable().Find(receiver); ok {
		nArgs++
		receiver = sym.Symtype
		if err := e.writePush(sym.Kind.String(), sym.Index); err != nil {
			panic(err)
		}
	}
	e.resetSubroutine()
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

func (e *engine) callArray() {
	if err := e.writePop("pointer", 1); err != nil {
		panic(err)
	}
	if err := e.writePush("that", 0); err != nil {
		panic(err)
	}
}

func (e *engine) allocateMemory() {
	size := e.ClassTable().VarCount("field")
	if err := e.writePush("constant", size); err != nil {
		panic(err)
	}
	if err := e.writeCall("Memory", "alloc", 1); err != nil {
		panic(err)
	}
	if err := e.writePop("pointer", 0); err != nil {
		panic(err)
	}
}

func (e *engine) setPointer() {
	if err := e.writePush("argument", 0); err != nil {
		panic(err)
	}
	if err := e.writePop("pointer", 0); err != nil {
		panic(err)
	}
}

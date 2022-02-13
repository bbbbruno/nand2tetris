package engine

import (
	"bufio"
	"errors"
	"io"

	"jackcompiler/symtable"
	"jackcompiler/tokenizer"
)

type scope bool

const CLASS, SUBROUTINE scope = true, false

type function struct {
	receiver, name  string
	expressionCount int
}

type buf struct {
	subroutine struct{ returnval, name, kind string }
	variable   struct{ name, symtype, kind string }
	sym        *symtable.Symbol
	term       struct{ value, kind string }
	functions  []*function
	operators  []string
	whileCount int
	ifCount    int
}

type engine struct {
	className string
	Tokenizer
	SymbolTable
	*bufio.Writer
	hierarchy int
	*buf
	scope
}

type Tokenizer interface {
	Advance()
	Peek() *tokenizer.Token
	CurrentToken() *tokenizer.Token
}

type SymbolTable interface {
	ClassTable() symtable.Table
	SubroutineTable() symtable.Table
	ResetSubroutineTable()
}

func New(className string, tkz Tokenizer, w io.Writer) *engine {
	return &engine{className, tkz, symtable.New(), bufio.NewWriter(w), 0, &buf{}, CLASS}
}

func (e *engine) Compile() (err error) {
	// defer func() {
	// 	if rec := recover(); rec != nil {
	// 		err = rec.(error)
	// 	}
	// }()
	e.Advance()
	e.compileClass()
	if err := e.Flush(); err != nil {
		return err
	}

	return err
}

const (
	OPEN  = true
	CLOSE = false
)

func (e *engine) compileClass() {
	e.validateKeyword("class")
	e.validateIdentifier()
	e.validateSymbol("{")
	for e.CurrentToken().IsKeyword("static", "field") {
		e.compileClassVarDec()
	}
	for e.CurrentToken().IsKeyword("constructor", "function", "method") {
		e.compileSubroutineDec()
	}
	e.validateSymbol("}")
}

func (e *engine) compileClassVarDec() {
	e.validateKeyword("static", "field")
	e.validateType()
	e.validateVarName()
	for e.CurrentToken().IsSymbol(",") {
		e.validateSymbol(",")
		e.validateVarName()
	}
	e.validateSymbol(";")
}

func (e *engine) compileSubroutineDec() {
	e.ResetSubroutineTable()
	e.scope = SUBROUTINE
	defer func() {
		e.scope = CLASS
		e.whileCount = 0
		e.ifCount = 0
	}()
	e.validateKeyword("constructor", "function", "method")
	if e.subroutine.kind == "method" {
		e.variable.name, e.variable.symtype, e.variable.kind = "_", "_", "argument"
		e.addVar()
	}
	e.validateKeywordOrType("void")
	e.validateSubroutineName()
	e.validateSymbol("(")
	e.compileParameterList()
	e.validateSymbol(")")
	e.compileSubroutineBody()
}

func (e *engine) compileParameterList() {
	if token := e.CurrentToken(); !token.IsKeyword(primitiveTypes...) && !token.IsIdentifier() {
		return
	}

	e.variable.kind = "argument"
	e.validateType()
	e.validateVarName()
	for e.CurrentToken().IsSymbol(",") {
		e.validateSymbol(",")
		e.validateType()
		e.validateVarName()
	}
}

func (e *engine) compileSubroutineBody() {
	e.validateSymbol("{")
	for e.CurrentToken().IsKeyword("var") {
		e.compileVarDec()
	}
	e.declareSubroutine()
	if e.subroutine.kind == "constructor" {
		e.allocateMemory()
	} else if e.subroutine.kind == "method" {
		e.setPointer()
	}
	e.compileStatements()
	e.validateSymbol("}")
}

func (e *engine) compileVarDec() {
	e.validateKeyword("var")
	e.validateType()
	e.validateVarName()
	for e.CurrentToken().IsSymbol(",") {
		e.validateSymbol(",")
		e.validateVarName()
	}
	e.validateSymbol(";")
}

func (e *engine) compileStatements() {
	for e.CurrentToken().IsKeyword("let", "if", "while", "do", "return") {
		token := e.CurrentToken()
		switch true {
		case token.IsKeyword("let"):
			e.compileLetStatement()
		case token.IsKeyword("if"):
			e.compileIfStatement()
		case token.IsKeyword("while"):
			e.compileWhileStatement()
		case token.IsKeyword("do"):
			e.compileDoStatement()
		case token.IsKeyword("return"):
			e.compileReturnStatement()
		default:
			panic(errors.New("token is not valid as statement"))
		}
	}
}

func (e *engine) compileLetStatement() {
	e.validateKeyword("let")
	e.validateVarName()
	sym := e.sym
	isArray := e.CurrentToken().IsSymbol("[")
	if isArray {
		e.validateSymbol("[")
		e.compileExpression()
		e.validateSymbol("]")
		e.calcArray(sym)
	}
	e.validateSymbol("=")
	e.compileExpression()
	e.validateSymbol(";")
	if isArray {
		e.letArrayStatment()
	} else {
		e.letStatement(sym)
	}
}

func (e *engine) compileIfStatement() {
	count := e.ifCount
	e.ifCount++
	e.validateKeyword("if")
	e.validateSymbol("(")
	e.compileExpression()
	e.validateSymbol(")")
	e.validateSymbol("{")
	e.startIf(count)
	e.compileStatements()
	e.validateSymbol("}")
	if token := e.CurrentToken(); !token.IsKeyword("else") {
		e.elseIf(true, count)
		return
	}
	e.validateKeyword("else")
	e.validateSymbol("{")
	e.elseIf(false, count)
	e.compileStatements()
	e.validateSymbol("}")
	e.endIf(count)
}

func (e *engine) compileWhileStatement() {
	count := e.whileCount
	e.whileCount++
	e.validateKeyword("while")
	e.startWhile(count)
	e.validateSymbol("(")
	e.compileExpression()
	e.jumpWhile(count)
	e.validateSymbol(")")
	e.validateSymbol("{")
	e.compileStatements()
	e.validateSymbol("}")
	e.endWhile(count)
}

func (e *engine) compileDoStatement() {
	e.validateKeyword("do")
	e.validateFunctionName()
	if token := e.CurrentToken(); token.IsSymbol("(") {
		e.callReceiver()
		e.validateSymbol("(")
		e.compileExpressionList()
		e.validateSymbol(")")
	} else if token.IsSymbol(".") {
		e.validateSymbol(".")
		e.validateReceiverName()
		e.callReceiver()
		e.validateSymbol("(")
		e.compileExpressionList()
		e.validateSymbol(")")
	} else {
		panic(errors.New("token is not valid as do statement"))
	}
	e.validateSymbol(";")
	e.doStatement()
}

func (e *engine) compileReturnStatement() {
	e.validateKeyword("return")
	if token := e.CurrentToken(); token.IsIntConst() || token.IsStringConst() || token.IsKeyword("true", "false", "null", "this") || token.IsIdentifier() || token.IsSymbol("(", "-", "~") {
		e.compileExpression()
	}
	e.validateSymbol(";")
	e.returnStatement()
}

func (e *engine) compileExpressionList() {
	if token := e.CurrentToken(); !token.IsIntConst() && !token.IsStringConst() && !token.IsKeyword("true", "false", "null", "this") && !token.IsIdentifier() && !token.IsSymbol("(", "-", "~") {
		return
	}

	function := e.functions[len(e.functions)-1]
	function.expressionCount++
	e.compileExpression()
	for e.CurrentToken().IsSymbol(",") {
		function.expressionCount++
		e.validateSymbol(",")
		e.compileExpression()
	}
}

var operators = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func (e *engine) compileExpression() {
	e.compileTerm()
	for e.CurrentToken().IsSymbol(operators...) {
		e.validateOperator(operators...)
		e.compileTerm()
		e.calcExpression()
	}
}

func (e *engine) compileTerm() {
	token := e.CurrentToken()
	switch true {
	case token.IsIntConst():
		e.validateIntConst()
		e.callIntConst()
	case token.IsStringConst():
		e.validateStringConst()
		e.callStringConst()
	case token.IsKeyword("true", "false", "null", "this"):
		e.term.value = e.CurrentToken().Content()
		e.validateKeyword("true", "false", "null", "this")
		e.callKeywordConst()
	case token.IsSymbol("("):
		e.validateSymbol("(")
		e.compileExpression()
		e.validateSymbol(")")
	case token.IsSymbol("-", "~"):
		e.validateOperator("-", "~")
		e.compileTerm()
		e.calcUnary()
	case token.IsIdentifier():
		nextToken := e.Peek()
		if nextToken == nil {
			return
		}
		switch true {
		case nextToken.IsSymbol("["):
			e.validateVarName()
			sym := e.sym
			e.validateSymbol("[")
			e.compileExpression()
			e.validateSymbol("]")
			e.calcArray(sym)
			e.callArray()
		case nextToken.IsSymbol("("):
			e.validateFunctionName()
			e.callReceiver()
			e.validateSymbol("(")
			e.compileExpressionList()
			e.validateSymbol(")")
			e.callFunc()
		case nextToken.IsSymbol("."):
			e.validateFunctionName()
			e.validateSymbol(".")
			e.validateReceiverName()
			e.callReceiver()
			e.validateSymbol("(")
			e.compileExpressionList()
			e.validateSymbol(")")
			e.callFunc()
		default:
			e.validateVarName()
			e.callVar()
		}
	default:
		panic(errors.New("token is not valid as term"))
	}
}

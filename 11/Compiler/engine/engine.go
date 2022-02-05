package engine

import (
	"bufio"
	"errors"
	"io"

	"jackcompiler/tokenizer"
)

type engine struct {
	engineTokenizer
	*bufio.Writer
	hierarchy int
}

type engineTokenizer interface {
	Advance()
	Peek() *tokenizer.Token
	CurrentToken() *tokenizer.Token
}

func New(tkz engineTokenizer, w io.Writer) *engine {
	return &engine{tkz, bufio.NewWriter(w), 0}
}

func (e *engine) Compile() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = rec.(error)
		}
	}()
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
	e.writeHierarchy("class", OPEN)
	defer e.writeHierarchy("class", CLOSE)
	e.writeKeyword("class")
	e.writeIdentifier()
	e.writeSymbol("{")
	for e.CurrentToken().IsKeyword("static", "field") {
		e.compileClassVarDec()
	}
	for e.CurrentToken().IsKeyword("constructor", "function", "method") {
		e.compileSubroutineDec()
	}
	e.writeSymbol("}")
}

func (e *engine) compileClassVarDec() {
	e.writeHierarchy("classVarDec", OPEN)
	defer e.writeHierarchy("classVarDec", CLOSE)
	e.writeKeyword("static", "field")
	e.writeType()
	e.writeIdentifier()
	for e.CurrentToken().IsSymbol(",") {
		e.writeSymbol(",")
		e.writeIdentifier()
	}
	e.writeSymbol(";")
}

func (e *engine) compileSubroutineDec() {
	e.writeHierarchy("subroutineDec", OPEN)
	defer e.writeHierarchy("subroutineDec", CLOSE)
	e.writeKeyword("constructor", "function", "method")
	e.writeKeywordOrType("void")
	e.writeIdentifier()
	e.writeSymbol("(")
	e.compileParameterList()
	e.writeSymbol(")")
	e.compileSubroutineBody()
}

func (e *engine) compileParameterList() {
	e.writeHierarchy("parameterList", OPEN)
	defer e.writeHierarchy("parameterList", CLOSE)
	if token := e.CurrentToken(); !token.IsKeyword(primitiveTypes...) && !token.IsIdentifier() {
		return
	}

	e.writeType()
	e.writeIdentifier()
	for e.CurrentToken().IsSymbol(",") {
		e.writeSymbol(",")
		e.writeType()
		e.writeIdentifier()
	}
}

func (e *engine) compileSubroutineBody() {
	e.writeHierarchy("subroutineBody", OPEN)
	defer e.writeHierarchy("subroutineBody", CLOSE)
	e.writeSymbol("{")
	for e.CurrentToken().IsKeyword("var") {
		e.compileVarDec()
	}
	e.compileStatements()
	e.writeSymbol("}")
}

func (e *engine) compileVarDec() {
	e.writeHierarchy("varDec", OPEN)
	defer e.writeHierarchy("varDec", CLOSE)
	e.writeKeyword("var")
	e.writeType()
	e.writeIdentifier()
	for e.CurrentToken().IsSymbol(",") {
		e.writeSymbol(",")
		e.writeIdentifier()
	}
	e.writeSymbol(";")
}

func (e *engine) compileStatements() {
	e.writeHierarchy("statements", OPEN)
	defer e.writeHierarchy("statements", CLOSE)
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
	e.writeHierarchy("letStatement", OPEN)
	defer e.writeHierarchy("letStatement", CLOSE)
	e.writeKeyword("let")
	e.writeIdentifier()
	if token := e.CurrentToken(); token.IsSymbol("[") {
		e.writeSymbol("[")
		e.compileExpression()
		e.writeSymbol("]")
	}
	e.writeSymbol("=")
	e.compileExpression()
	e.writeSymbol(";")
}

func (e *engine) compileIfStatement() {
	e.writeHierarchy("ifStatement", OPEN)
	defer e.writeHierarchy("ifStatement", CLOSE)
	e.writeKeyword("if")
	e.writeSymbol("(")
	e.compileExpression()
	e.writeSymbol(")")
	e.writeSymbol("{")
	e.compileStatements()
	e.writeSymbol("}")
	if token := e.CurrentToken(); !token.IsKeyword("else") {
		return
	}
	e.writeKeyword("else")
	e.writeSymbol("{")
	e.compileStatements()
	e.writeSymbol("}")
}

func (e *engine) compileWhileStatement() {
	e.writeHierarchy("whileStatement", OPEN)
	defer e.writeHierarchy("whileStatement", CLOSE)
	e.writeKeyword("while")
	e.writeSymbol("(")
	e.compileExpression()
	e.writeSymbol(")")
	e.writeSymbol("{")
	e.compileStatements()
	e.writeSymbol("}")
}

func (e *engine) compileDoStatement() {
	e.writeHierarchy("doStatement", OPEN)
	defer e.writeHierarchy("doStatement", CLOSE)
	e.writeKeyword("do")
	e.writeIdentifier()
	if token := e.CurrentToken(); token.IsSymbol("(") {
		e.writeSymbol("(")
		e.compileExpressionList()
		e.writeSymbol(")")
	} else if token.IsSymbol(".") {
		e.writeSymbol(".")
		e.writeIdentifier()
		e.writeSymbol("(")
		e.compileExpressionList()
		e.writeSymbol(")")
	} else {
		panic(errors.New("token is not valid as do statement"))
	}
	e.writeSymbol(";")
}

func (e *engine) compileReturnStatement() {
	e.writeHierarchy("returnStatement", OPEN)
	defer e.writeHierarchy("returnStatement", CLOSE)
	e.writeKeyword("return")
	if token := e.CurrentToken(); token.IsIntConst() || token.IsStringConst() || token.IsKeyword("true", "false", "null", "this") || token.IsIdentifier() || token.IsSymbol("(", "-", "~") {
		e.compileExpression()
	}
	e.writeSymbol(";")
}

var operators = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func (e *engine) compileExpression() {
	e.writeHierarchy("expression", OPEN)
	defer e.writeHierarchy("expression", CLOSE)
	e.compileTerm()
	for e.CurrentToken().IsSymbol(operators...) {
		e.writeSymbol(operators...)
		e.compileTerm()
	}
}

func (e *engine) compileTerm() {
	e.writeHierarchy("term", OPEN)
	defer e.writeHierarchy("term", CLOSE)
	token := e.CurrentToken()
	switch true {
	case token.IsIntConst():
		e.writeIntConst()
	case token.IsStringConst():
		e.writeStringConst()
	case token.IsKeyword("true", "false", "null", "this"):
		e.writeKeyword("true", "false", "null", "this")
	case token.IsSymbol("("):
		e.writeSymbol("(")
		e.compileExpression()
		e.writeSymbol(")")
	case token.IsSymbol("-", "~"):
		e.writeSymbol("-", "~")
		e.compileTerm()
	case token.IsIdentifier():
		nextToken := e.Peek()
		if nextToken == nil {
			return
		}
		switch true {
		case nextToken.IsSymbol("["):
			e.writeIdentifier()
			e.writeSymbol("[")
			e.compileExpression()
			e.writeSymbol("]")
		case nextToken.IsSymbol("("):
			e.writeIdentifier()
			e.writeSymbol("(")
			e.compileExpressionList()
			e.writeSymbol(")")
		case nextToken.IsSymbol("."):
			e.writeIdentifier()
			e.writeSymbol(".")
			e.writeIdentifier()
			e.writeSymbol("(")
			e.compileExpressionList()
			e.writeSymbol(")")
		default:
			e.writeIdentifier()
		}
	default:
		panic(errors.New("token is not valid as term"))
	}
}

func (e *engine) compileExpressionList() {
	e.writeHierarchy("expressionList", OPEN)
	defer e.writeHierarchy("expressionList", CLOSE)
	if token := e.CurrentToken(); !token.IsIntConst() && !token.IsStringConst() && !token.IsKeyword("true", "false", "null", "this") && !token.IsIdentifier() && !token.IsSymbol("(", "-", "~") {
		return
	}

	e.compileExpression()
	for e.CurrentToken().IsSymbol(",") {
		e.writeSymbol(",")
		e.compileExpression()
	}
}

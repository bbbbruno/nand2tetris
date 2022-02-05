package tokenizer

import (
	"regexp"
	"strconv"
	"unicode"
)

type tokentype int

const (
	KEYWORD tokentype = iota + 1
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

var tokentypes = map[tokentype]string{
	KEYWORD:      "keyword",
	SYMBOL:       "symbol",
	IDENTIFIER:   "identifier",
	INT_CONST:    "integerConstant",
	STRING_CONST: "stringConstant",
}

func (tt tokentype) String() string {
	return tokentypes[tt]
}

var keywords = [...]string{
	"class", "constructor", "function", "method", "field", "static",
	"var", "int", "char", "boolean", "void",
	"true", "false", "null", "this", "let", "do",
	"if", "else", "while", "return",
}

var symbols = [...]byte{
	'{', '}', '(', ')', '[', ']',
	'.', ',', ';', '+', '-', '*', '/',
	'&', '|', '<', '>', '=', '~',
}

func isSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

func isComment(b byte, nextb byte) bool {
	return b == '/' && (nextb == '/' || nextb == '*')
}

func isSymbol(b byte) bool {
	for _, sym := range symbols {
		if b == sym {
			return true
		}
	}

	return false
}

func isKeyword(buf []byte) bool {
	for _, key := range keywords {
		if string(buf) == key {
			return true
		}
	}

	return false
}

func isIntConst(buf []byte) bool {
	_, err := strconv.Atoi(string(buf))
	return err == nil
}

func isStringConst(b byte) bool {
	return b == '"'
}

var isIdentifier = regexp.MustCompile(`\w`).MatchString

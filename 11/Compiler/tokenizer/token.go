package tokenizer

import (
	"fmt"
	"html"
)

type Token struct {
	tokentype tokentype
	content   string
}

func (t Token) Type() tokentype {
	return t.tokentype
}

func (t Token) Content() string {
	return t.content
}

func (t Token) String() string {
	return fmt.Sprintf("<%[1]s> %[2]s </%[1]s>", t.tokentype, html.EscapeString(t.content))
}

func (t Token) IsKeyword(keys ...string) bool {
	if len(keys) == 0 {
		return t.tokentype == KEYWORD
	}
	return t.tokentype == KEYWORD && contains(keys, t.content)
}

func (t Token) IsSymbol(syms ...string) bool {
	if len(syms) == 0 {
		return t.tokentype == SYMBOL
	}
	return t.tokentype == SYMBOL && contains(syms, t.content)
}

func (t Token) IsIdentifier() bool {
	return t.tokentype == IDENTIFIER
}

func (t Token) IsIntConst() bool {
	return t.tokentype == INT_CONST
}

func (t Token) IsStringConst() bool {
	return t.tokentype == STRING_CONST
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

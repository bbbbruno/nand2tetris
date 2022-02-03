package tokenizer

import (
	"fmt"
	"html"
)

type Token interface {
	Type() tokentype
	Content() string
	XmlString() string
}

type token struct {
	tokentype tokentype
	content   string
}

func (t *token) Type() tokentype {
	return t.tokentype
}

func (t *token) Content() string {
	return t.content
}

func (t *token) XmlString() string {
	return fmt.Sprintf("<%[1]s> %[2]s </%[1]s>", t.tokentype, html.EscapeString(t.content))
}

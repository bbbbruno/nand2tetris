package tokenizer

type Token interface {
	Type() tokentype
	Content() string
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

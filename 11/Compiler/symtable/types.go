package symtable

import "fmt"

type kind int

const (
	STATIC kind = iota + 1
	FIELD
	ARG
	VAR
)

var kindStrings = [...]string{"", "static", "field", "argument", "local"}
var kinds = map[string]kind{
	"static":   STATIC,
	"field":    FIELD,
	"argument": ARG,
	"var":      VAR,
}

func (k kind) String() string {
	return kindStrings[k]
}

type symbol struct {
	name    string
	symtype string
	kind    kind
	index   int
}

func (s *symbol) String() string {
	return fmt.Sprintf("%s %s %s %d", s.kind, s.symtype, s.name, s.index)
}

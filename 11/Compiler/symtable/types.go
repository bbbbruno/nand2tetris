package symtable

import "fmt"

type kind int

const (
	STATIC kind = iota + 1
	FIELD
	ARG
	VAR
)

var kinds = map[string]kind{
	"static":   STATIC,
	"field":    FIELD,
	"argument": ARG,
	"var":      VAR,
}
var kindStrings = [...]string{"", "static", "field", "argument", "local"}

func (k kind) String() string {
	return kindStrings[k]
}

type Symbol struct {
	Name    string
	Symtype string
	Kind    kind
	Index   int
}

func (s *Symbol) String() string {
	return fmt.Sprintf("%s %s %s %d", s.Kind, s.Symtype, s.Name, s.Index)
}

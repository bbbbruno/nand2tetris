package symtable

import (
	"errors"
)

type Table interface {
	Find(name string) (*Symbol, bool)
	VarCount(kind string) int
	Define(name, symtype, kind string) error
}

type table struct {
	symbols []*Symbol
}

func (t *table) VarCount(kind string) int {
	k, ok := kinds[kind]
	if !ok {
		return 0
	}

	arr := make([]*Symbol, 0)
	for _, v := range t.symbols {
		if v.Kind == k {
			arr = append(arr, v)
		}
	}
	return len(arr)
}

func (t *table) Find(name string) (*Symbol, bool) {
	for _, sym := range t.symbols {
		if sym.Name == name {
			return sym, true
		}
	}

	return nil, false
}

type classTable struct {
	*table
}

func (t *classTable) Define(name string, symtype string, kind string) error {
	k := kinds[kind]
	if k == STATIC || k == FIELD {
		sym := &Symbol{name, symtype, k, t.VarCount(kind)}
		t.symbols = append(t.symbols, sym)
		return nil
	} else {
		return errors.New("invalid type, expected STATIC or FIELD")
	}
}

type subroutineTable struct {
	*table
}

func (t *subroutineTable) Define(name string, symtype string, kind string) error {
	k := kinds[kind]
	if k == ARG || k == VAR {
		sym := &Symbol{name, symtype, k, t.VarCount(kind)}
		t.symbols = append(t.symbols, sym)
		return nil
	} else {
		return errors.New("invalid type, expected ARG or VAR")
	}
}

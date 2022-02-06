package symtable

import (
	"errors"
)

type Table interface {
	Define(name, symtype, kind string) (*Symbol, error)
	VarCount(kind string) int
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

func (t *table) find(name string) (*Symbol, bool) {
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

func (t *classTable) Define(name string, symtype string, kind string) (*Symbol, error) {
	if sym, ok := t.find(name); ok {
		return sym, nil
	}

	k := kinds[kind]
	if k == STATIC || k == FIELD {
		sym := &Symbol{name, symtype, k, t.VarCount(kind)}
		t.symbols = append(t.symbols, sym)
		return sym, nil
	} else {
		return nil, errors.New("invalid type, expected STATIC or FIELD")
	}
}

type subroutineTable struct {
	*table
}

func (t *subroutineTable) Define(name string, symtype string, kind string) (*Symbol, error) {
	if sym, ok := t.find(name); ok {
		return sym, nil
	}

	k := kinds[kind]
	if k == ARG || k == VAR {
		sym := &Symbol{name, symtype, k, t.VarCount(kind)}
		t.symbols = append(t.symbols, sym)
		return sym, nil
	} else {
		return nil, errors.New("invalid type, expected ARG or VAR")
	}
}

package symtable

import (
	"errors"
)

type Table interface {
	Define(name, symtype, kind string) (*symbol, error)
}

type table struct {
	symbols []*symbol
}

func (t *table) varCount(kind kind) int {
	arr := make([]*symbol, 0)
	for _, v := range t.symbols {
		if v.kind == kind {
			arr = append(arr, v)
		}
	}
	return len(arr)
}

func (t *table) find(name string) (*symbol, bool) {
	for _, sym := range t.symbols {
		if sym.name == name {
			return sym, true
		}
	}

	return nil, false
}

type classTable struct {
	*table
}

func (t *classTable) Define(name string, symtype string, kind string) (*symbol, error) {
	if sym, ok := t.find(name); ok {
		return sym, nil
	}

	k := kinds[kind]
	if k == STATIC || k == FIELD {
		sym := &symbol{name, symtype, k, t.varCount(k)}
		t.symbols = append(t.symbols, sym)
		return sym, nil
	} else {
		return nil, errors.New("invalid type, expected STATIC or FIELD")
	}
}

type subroutineTable struct {
	*table
}

func (t *subroutineTable) Define(name string, symtype string, kind string) (*symbol, error) {
	if sym, ok := t.find(name); ok {
		return sym, nil
	}

	k := kinds[kind]
	if k == ARG || k == VAR {
		sym := &symbol{name, symtype, k, t.varCount(k)}
		t.symbols = append(t.symbols, sym)
		return sym, nil
	} else {
		return nil, errors.New("invalid type, expected ARG or VAR")
	}
}
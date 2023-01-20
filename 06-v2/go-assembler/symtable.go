package assemble

type SymTable struct {
	NextAddr int
	Table    map[string]int
}

const initialAddr = 16

var defaultSymbols = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

func NewSymTable() *SymTable {
	st := &SymTable{NextAddr: initialAddr, Table: make(map[string]int)}
	for k, v := range defaultSymbols {
		st.Table[k] = v
	}
	return st
}

func (st *SymTable) AddSymbol(sym string) *SymTable {
	st.Table[sym] = st.NextAddr
	st.NextAddr++
	return st
}

func (st *SymTable) AddSymbolWithAddr(sym string, addr int) *SymTable {
	st.Table[sym] = addr
	return st
}

func (st *SymTable) Contains(sym string) bool {
	_, ok := st.Table[sym]
	return ok
}

func (st *SymTable) Addr(sym string) int {
	return st.Table[sym]
}

package symtable

type symboltable struct {
	classTable      *classTable
	subroutineTable *subroutineTable
}

func New() *symboltable {
	classTable := &classTable{&table{make([]*Symbol, 0)}}
	subroutineTable := &subroutineTable{&table{make([]*Symbol, 0)}}
	return &symboltable{classTable, subroutineTable}
}

func (st *symboltable) ClassTable() Table {
	return st.classTable
}

func (st *symboltable) SubroutineTable() Table {
	return st.subroutineTable
}

func (st *symboltable) ResetSubroutineTable() {
	st.subroutineTable.symbols = make([]*Symbol, 0)
}

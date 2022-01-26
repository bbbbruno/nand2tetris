package main

import "strconv"

type symbol string

func NewSymbol(s string) symbol {
	return symbol(s)
}

func (s symbol) String() string {
	return string(s)
}

func (s symbol) IsConst() bool {
	if _, err := strconv.Atoi(string(s)); err == nil {
		return true
	}

	return false
}

type addr uint

type SymbolTable struct {
	nextAddr addr
	symbols  map[symbol]addr
}

var DefaultSymbols = map[symbol]addr{
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

// 定義済みシンボル含めたシンボルテーブルのインスタンスを返す。
// ユーザー定義シンボルのアドレスは１６から始まる。
func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{16, make(map[symbol]addr)}
	for k, v := range DefaultSymbols {
		st.symbols[k] = v
	}
	return st
}

// 定義済みシンボル含めたシンボルテーブルのインスタンスを返す。
// 開始アドレスと初期シンボルを指定できる。
func NewSymbolTableWithOpts(addr addr, symbols map[symbol]addr) *SymbolTable {
	st := &SymbolTable{addr, symbols}
	for k, v := range DefaultSymbols {
		st.symbols[k] = v
	}
	return st
}

// シンボルテーブルにシンボルを追加する。
func (st *SymbolTable) AddSymbol(sym symbol) {
	st.symbols[sym] = st.nextAddr
	st.nextAddr += 1
}

// アドレスを直接指定してシンボルテーブルにシンボルを追加する。
func (st *SymbolTable) AddSymbolWithAddr(sym symbol, addr addr) {
	st.symbols[sym] = addr
}

// シンボルテーブル内のシンボルのアドレスと見つかったかどうかを返す。
// シンボルが定数値の場合はaddr型に変換して返す。
func (st *SymbolTable) Addr(sym symbol) (addr, bool) {
	if i, ok := st.symbols[sym]; ok {
		return i, true
	}

	if sym.IsConst() {
		i, _ := strconv.Atoi(sym.String())
		return addr(i), true
	}

	return 0, false
}

package main

import (
	"bufio"
	"io"
	"strings"
)

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

// 入力から一行ずつ解析する作業を二度行い、コマンドの配列とシンボルテーブルを作成する。
// 1回目の解析ではLコマンドのみ解析し、次の行のアドレスをシンボルテーブルに格納する。
// 2回目の解析ではAコマンドとCコマンドのみを解析し、ユーザー定義シンボルもシンボルテーブルに格納する。
func (p Parser) Parse(input io.Reader) (cmds []*Command, st *SymbolTable) {
	lines := parseInput(input)

	st = parseSymbols(lines)
	cmds, st = parseLines(lines, st)

	return cmds, st
}

// 入力から一行ずつ読み取り、空白行とコメントアウトを取り除く。
func parseInput(input io.Reader) (lines []string) {
	s := bufio.NewScanner(input)

	for s.Scan() {
		line := s.Text()
		line = strings.ReplaceAll(line, " ", "")         // 空白除去
		if i := strings.LastIndex(line, "//"); i != -1 { // コメントアウトを除去
			line = line[:i]
		}
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return lines
}

// Lコマンドのみ解析し、ユーザー定義シンボルをシンボルテーブルに格納する。
func parseSymbols(lines []string) *SymbolTable {
	st := NewSymbolTable()
	var count addr

	for _, l := range lines {
		ok, sym := IsLCommand(l)
		if !ok {
			count += 1
			continue
		}

		if _, ok := st.Addr(sym); ok {
			continue
		}

		st.AddSymbolWithAddr(sym, count)
	}

	return st
}

// AコマンドとCコマンドのみ解析し、適時ユーザー定義シンボルをシンボルテーブルに格納する。
func parseLines(lines []string, st *SymbolTable) ([]*Command, *SymbolTable) {
	var cmds []*Command
	for _, l := range lines {
		if ok, _ := IsLCommand(l); ok {
			continue
		}

		cmd := new(Command)
		if ok, sym := IsACommand(l); ok {
			cmd.Type = A_COMMAND
			cmd.Symbol = sym
			if _, ok := st.Addr(cmd.Symbol); !ok {
				st.AddSymbol(cmd.Symbol)
			}
			cmds = append(cmds, cmd)
			continue
		}

		cmd.Type = C_COMMAND
		cmd.Comp = l
		if ss := strings.Split(l, "="); len(ss) == 2 {
			cmd.Dest = ss[0]
			cmd.Comp = strings.ReplaceAll(cmd.Comp, cmd.Dest+"=", "")
		}
		if ss := strings.Split(l, ";"); len(ss) == 2 {
			cmd.Jump = ss[1]
			cmd.Comp = strings.ReplaceAll(cmd.Comp, ";"+cmd.Jump, "")
		}
		cmds = append(cmds, cmd)
	}

	return cmds, st
}

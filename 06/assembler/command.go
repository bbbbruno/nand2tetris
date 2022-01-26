package main

import "fmt"

type cmdType uint8

const (
	A_COMMAND cmdType = iota
	C_COMMAND
	L_COMMAND
)

// TODO: オブジェクト指向っぽくする
type Command struct {
	Type             cmdType
	Symbol           symbol
	Dest, Comp, Jump string
}

// TODO: コマンド種別ごとにシンボルをバリデーションする
func IsACommand(s string) (bool, symbol) {
	return s[0] == '@', NewSymbol(s[1:])
}

func IsLCommand(s string) (bool, symbol) {
	return s[0] == '(' && s[len(s)-1] == ')', NewSymbol(s[1 : len(s)-1])
}

type BinaryCommand uint

func (bcmd BinaryCommand) String() string {
	return fmt.Sprintf("%016b", bcmd)
}

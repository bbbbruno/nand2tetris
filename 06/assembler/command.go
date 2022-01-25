package main

import "fmt"

type cmdType uint8

const (
	A_COMMAND cmdType = iota
	C_COMMAND
)

type Command struct {
	Type             cmdType
	Symbol           string
	Dest, Comp, Jump string
}

type BinaryCommand uint

func (bcmd BinaryCommand) String() string {
	return fmt.Sprintf("%016b", bcmd)
}

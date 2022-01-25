package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Code struct{}

func NewCode() Code {
	return Code{}
}

var translatorsMap = map[cmdType]func(*Command) (BinaryCommand, error){
	A_COMMAND: translateACommand,
	C_COMMAND: translateCCommand,
}

func (c Code) Translate(cmds []*Command) (bcmds []BinaryCommand, err error) {
	for _, cmd := range cmds {
		t, ok := translatorsMap[cmd.Type]
		if !ok {
			return nil, errors.New("unknown command type")
		}

		bcmd, err := t(cmd)
		if err != nil {
			return nil, err
		}
		bcmds = append(bcmds, bcmd)
	}

	return bcmds, nil
}

func translateACommand(cmd *Command) (BinaryCommand, error) {
	sym, err := strconv.Atoi(cmd.Symbol)
	if err != nil {
		return 0, errors.New("symbol must be int")
	}
	binary := fmt.Sprintf("%b%015b", cmd.Type, sym)
	i, err := strconv.ParseInt(binary, 2, 0)
	if err != nil {
		return 0, err
	}
	return BinaryCommand(i), nil
}

func translateCCommand(cmd *Command) (BinaryCommand, error) {
	dest, err := translateDest(cmd.Dest)
	if err != nil {
		return 0, err
	}
	comp, err := translateComp(cmd.Comp)
	if err != nil {
		return 0, err
	}
	jump, err := translateJump(cmd.Jump)
	if err != nil {
		return 0, err
	}
	binary := fmt.Sprintf("%b11%07b%03b%03b", cmd.Type, comp, dest, jump)
	i, err := strconv.ParseInt(binary, 2, 0)
	if err != nil {
		return 0, err
	}
	return BinaryCommand(i), nil
}

var destsMap = map[string]uint8{
	"":    0b000,
	"M":   0b001,
	"D":   0b010,
	"MD":  0b011,
	"A":   0b100,
	"AM":  0b101,
	"AD":  0b110,
	"AMD": 0b111,
}

func translateDest(dest string) (uint8, error) {
	binary, ok := destsMap[dest]
	if !ok {
		return 0, errors.New("invalid dest specified")
	}

	return binary, nil
}

var compsMap = map[string]uint8{
	// a = 0
	"0":   0b0101010,
	"1":   0b0111111,
	"-1":  0b0111010,
	"D":   0b0001100,
	"A":   0b0110000,
	"!D":  0b0001101,
	"!A":  0b0110001,
	"-D":  0b0001111,
	"-A":  0b0110011,
	"D+1": 0b0011111,
	"A+1": 0b0110111,
	"D-1": 0b0001110,
	"A-1": 0b0110010,
	"D+A": 0b0000010,
	"D-A": 0b0010011,
	"A-D": 0b0000111,
	"D&A": 0b0000000,
	"D|A": 0b0010101,
	// a = 1
	"M":   0b1110000,
	"!M":  0b1110001,
	"-M":  0b1110011,
	"M+1": 0b1110111,
	"M-1": 0b1110010,
	"D+M": 0b1000010,
	"D-M": 0b1010011,
	"M-D": 0b1000111,
	"D&M": 0b1000000,
	"D|M": 0b1010101,
}

func translateComp(comp string) (uint8, error) {
	binary, ok := compsMap[comp]
	if !ok {
		return 0, errors.New("invalid comp specified")
	}

	return binary, nil
}

var jumpsMap = map[string]uint8{
	"":    0b000,
	"JGT": 0b001,
	"JEQ": 0b010,
	"JGE": 0b011,
	"JLT": 0b100,
	"JNE": 0b101,
	"JLE": 0b110,
	"JMP": 0b111,
}

func translateJump(jump string) (uint8, error) {
	binary, ok := jumpsMap[jump]
	if !ok {
		return 0, errors.New("invalid jump specified")
	}

	return binary, nil
}

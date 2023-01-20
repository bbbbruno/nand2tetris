package assemble

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type InvalidSymbolError struct {
	symbol string
}

func (e InvalidSymbolError) Error() string {
	return fmt.Sprintf("Invalid symbol %s. "+
		"You can only use alphabets, numbers, and symbols _.$: as you like. "+
		"You can't start with numbers.", e.symbol)
}

func IsValidSymbol(s string) bool {
	if match, err := regexp.MatchString(`^[^0-9][a-zA-Z0-9_.$:]*$`, s); err != nil {
		return false
	} else {
		return match
	}
}

type Cmd interface {
	Binary() string
}

type ACmd struct {
	Symbol string
	Value  int
}

var _ Cmd = (*ACmd)(nil)

func NewACmd(s string, st *SymTable) (*ACmd, error) {
	sym := s[1:]
	v, err := strconv.Atoi(sym)
	if err == nil {
		return &ACmd{Symbol: "", Value: v}, nil
	}

	if !IsValidSymbol(sym) {
		return nil, &InvalidSymbolError{symbol: sym}
	}

	if !st.Contains(sym) {
		st.AddSymbol(sym)
	}
	return &ACmd{Symbol: sym, Value: st.Addr(sym)}, nil
}

func (c ACmd) Binary() string {
	return fmt.Sprintf("0%015b", c.Value)
}

type CCmd struct {
	Dest, Comp, Jump string
}

var _ Cmd = (*CCmd)(nil)

var destMap = map[string]int{
	"":    0b000,
	"M":   0b001,
	"D":   0b010,
	"MD":  0b011,
	"A":   0b100,
	"AM":  0b101,
	"AD":  0b110,
	"AMD": 0b111,
}
var compMap = map[string]int{
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
var jumpMap = map[string]int{
	"":    0b000,
	"JGT": 0b001,
	"JEQ": 0b010,
	"JGE": 0b011,
	"JLT": 0b100,
	"JNE": 0b101,
	"JLE": 0b110,
	"JMP": 0b111,
}

func NewCCmd(s string) (*CCmd, error) {
	d, c, j := mnemonics(s)
	if _, ok := destMap[d]; !ok {
		return nil, fmt.Errorf("failed to create CCmd: invalid dest %s", d)
	}
	if _, ok := compMap[c]; !ok {
		return nil, fmt.Errorf("failed to create CCmd: invalid comp %s", c)
	}
	if _, ok := jumpMap[j]; !ok {
		return nil, fmt.Errorf("failed to create CCmd: invalid jump %s", j)
	}
	return &CCmd{Dest: d, Comp: c, Jump: j}, nil
}

func (c CCmd) Binary() string {
	return fmt.Sprintf("111%07b%03b%03b", compMap[c.Comp], destMap[c.Dest], jumpMap[c.Jump])
}

func mnemonics(s string) (d, c, j string) {
	c = s
	if ss := strings.Split(s, "="); len(ss) == 2 {
		d = ss[0]
		c = strings.Replace(c, d+"=", "", 1)
	}
	if ss := strings.Split(s, ";"); len(ss) == 2 {
		j = ss[1]
		c = strings.Replace(c, ";"+j, "", 1)
	}
	return d, c, j
}

package translator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"

	"vmconverter/vmcommand"
)

type Translator interface {
	SetNewFile(io.Writer)
	WriteArithmethic(string) error
	WritePushPop(vmcommand.CommandType, string, int) error
	Close() error
}

type codewriter struct {
	*bufio.Writer
	closed bool
}

func NewCodeWriter() Translator {
	return &codewriter{}
}

func (c *codewriter) SetNewFile(w io.Writer) {
	c.closed = false
	c.Writer = bufio.NewWriter(w)
}

func (c *codewriter) Close() error {
	if _, err := fmt.Fprint(c, end()); err != nil {
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	if c.closed {
		return errors.New("already closed")
	}

	c.closed = true
	return nil
}

var operationMap = map[string]string{
	"add": operateDouble("M+D"),
	"sub": operateDouble("M-D"),
	"neg": operateSingle("-M"),
	"and": operateDouble("M&D"),
	"or":  operateDouble("M|D"),
	"not": operateSingle("!M"),
}
var comparisonMap = map[string]func() string{
	"eq": func() string { return operateCompare("JEQ") },
	"gt": func() string { return operateCompare("JGT") },
	"lt": func() string { return operateCompare("JLT") },
}

// 算術コマンドを書き込む。
func (c *codewriter) WriteArithmethic(instruction string) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	var text string
	switch instruction {
	case "eq", "gt", "lt":
		text = comparisonMap[instruction]()
	default:
		text = operationMap[instruction]
	}

	_, err := fmt.Fprint(c, text)
	return err
}

// PUSHまたはPOPコマンドを書き込む。
func (c *codewriter) WritePushPop(instruction vmcommand.CommandType, segment string, index int) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	switch instruction {
	case vmcommand.C_PUSH:
		return c.writePush(segment, index)
	case vmcommand.C_POP:
		return c.writePop(segment, index)
	default:
		return nil
	}
}

func (c *codewriter) checkClosed() error {
	if c.closed {
		return errors.New("writer is closed")
	} else {
		return nil
	}
}

var segmentMap = map[string]func(int) string{
	"constant": func(index int) string { return constant(index) + push() },
}

func (c *codewriter) writePush(segment string, index int) error {
	text := segmentMap[segment](index)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) writePop(segment string, index int) error {
	return nil
}

func push() string {
	return `@SP
A=M
M=D
@SP
M=M+1
`
}

func pop(comp string) string {
	return fmt.Sprintf(`@SP
M=M-1
A=M
D=%s
M=0
`, comp)
}

func operateDouble(comp string) string {
	return pop("M") + pop(comp) + push()
}

func operateSingle(comp string) string {
	return pop(comp) + push()
}

func operateCompare(jump string) string {
	return pop("M") + pop("M-D") + compare(jump) + push()
}

func compare(jump string) string {
	label := rand.Intn(1000000)
	return fmt.Sprintf(`@TRUE%d
D;%s
D=0
@FINAL%d
0;JMP
(TRUE%d)
D=-1
(FINAL%d)
`, label, jump, label, label, label)
}

func constant(index int) string {
	return fmt.Sprintf(`@%d
D=A
`, index)
}

func end() string {
	return `(END)
@END
0;JMP
`
}

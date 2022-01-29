package translator

import (
	"bufio"
	"errors"
	"fmt"
	"io"

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
var comparisonMap = map[string]string{
	"eq": "JEQ",
	"gt": "JGT",
	"lt": "JLT",
}

// 算術コマンドを書き込む。
func (c *codewriter) WriteArithmethic(instruction string) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	var text string
	switch instruction {
	case "eq", "gt", "lt":
		text = operateCompare(comparisonMap[instruction])
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

func (c *codewriter) writePush(segment string, index int) error {
	var text string
	if segment == "constant" {
		text = constant(index) + push()
	} else {
		text = memoryPush(segment, index)
	}
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) writePop(segment string, index int) error {
	text := memoryPop(segment, index)
	_, err := fmt.Fprint(c, text)
	return err
}

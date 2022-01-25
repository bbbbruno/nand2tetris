package main

import (
	"errors"
	"testing"
)

func TestTranslate(t *testing.T) {
	err := errors.New("unknown command type")
	c := NewCode()
	testCases := []struct {
		in   *Command
		want string
		err  error
	}{
		{&Command{A_COMMAND, "100", "", "", ""}, "0000000001100100", nil},
		{&Command{C_COMMAND, "", "D", "D+1", "JGT"}, "1110011111010001", nil},
		{&Command{3, "", "D", "D+1", "JGT"}, "", err},
	}
	for _, test := range testCases {
		got, err := c.Translate(test.in)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("c.Translate(%#v) expected no error, got %#v", test.in, err)
		} else if test.err == nil && err != nil {
			t.Errorf("c.Translate(%#v) got %#v, want %#v", test.in, err, test.err)
		} else if got != test.want {
			t.Errorf("c.Translate(%v) got %v, want %v", test.in, got, test.want)
		}
	}
}

func TestTranslateACommand(t *testing.T) {
	err := errors.New("symbol must be int")
	testCases := []struct {
		in   *Command
		want string
		err  error
	}{
		{&Command{A_COMMAND, "100", "", "", ""}, "0000000001100100", nil},
		{&Command{A_COMMAND, "100g", "", "", ""}, "", err},
		{&Command{A_COMMAND, "", "", "", ""}, "", err},
	}
	for _, test := range testCases {
		got, err := translateACommand(test.in)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("c.Translate(%#v) expected no error, got %#v", test.in, err)
		} else if test.err == nil && err != nil {
			t.Errorf("c.Translate(%#v) got %#v, want %#v", test.in, err, test.err)
		} else if got != test.want {
			t.Errorf("c.Translate(%v) got %v, want %v", test.in, got, test.want)
		}
	}
}

func TestTranslateCCommand(t *testing.T) {
	// TODO: 実装する
}

func TestTranslateDest(t *testing.T) {
	err := errors.New("invalid dest specified")
	testCases := []struct {
		in   string
		want uint8
		err  error
	}{
		{"B", 0, err},
		{"", 0b000, nil},
		{"M", 0b001, nil},
		{"MD", 0b011, nil},
		{"AM", 0b101, nil},
		{"AMD", 0b111, nil},
	}
	for _, test := range testCases {
		got, err := translateDest(test.in)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("c.translateDest(%#v) expected no error, got %#v", test.in, err)
		} else if test.err == nil && err != nil {
			t.Errorf("c.translateDest(%#v) got %#v, want %#v", test.in, err, test.err)
		} else if got != test.want {
			t.Errorf("c.translateDest(%#v) got %#v, want %#v", test.in, got, test.want)
		}
	}
}

func TestTranslateComp(t *testing.T) {
	err := errors.New("invalid comp specified")
	testCases := []struct {
		in   string
		want uint8
		err  error
	}{
		{"", 0, err},
		{"B-A", 0, err},
		{"0", 0b0101010, nil},
		{"-1", 0b0111010, nil},
		{"D+1", 0b0011111, nil},
		{"D&A", 0b0000000, nil},
		{"M+1", 0b1110111, nil},
		{"M-D", 0b1000111, nil},
	}
	for _, test := range testCases {
		got, err := translateComp(test.in)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("c.translateComp(%#v) expected no error, got %#v", test.in, err)
		} else if test.err == nil && err != nil {
			t.Errorf("c.translateComp(%#v) got %#v, want %#v", test.in, err, test.err)
		} else if got != test.want {
			t.Errorf("c.translateComp(%#v) got %#v, want %#v", test.in, got, test.want)
		}
	}
}

func TestTranslateJump(t *testing.T) {
	err := errors.New("invalid jump specified")
	testCases := []struct {
		in   string
		want uint8
		err  error
	}{
		{"JWW", 0, err},
		{"", 0b000, nil},
		{"JGT", 0b001, nil},
		{"JGE", 0b011, nil},
		{"JNE", 0b101, nil},
		{"JMP", 0b111, nil},
	}
	for _, test := range testCases {
		got, err := translateJump(test.in)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("c.translateDest(%#v) expected no error, got %#v", test.in, err)
		} else if test.err == nil && err != nil {
			t.Errorf("c.translateDest(%#v) got %#v, want %#v", test.in, err, test.err)
		} else if got != test.want {
			t.Errorf("c.translateDest(%#v) got %#v, want %#v", test.in, got, test.want)
		}
	}
}

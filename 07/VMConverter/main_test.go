package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertSimpleAdd(t *testing.T) {
	vmFile, _ := os.ReadFile(filepath.Join("tests", "SimpleAdd", "SimpleAdd.vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("tests", "SimpleAdd", "SimpleAdd.asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

package main

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func execTest(t *testing.T, dirname string) {
	rand.Seed(1)
	root := filepath.Join("fixtures", dirname)
	paths, _, _ := FindVmFiles(root)
	got := bytes.NewBuffer(make([]byte, 0))
	asmFile, _ := os.ReadFile(filepath.Join(root, dirname+".asm"))
	expected := bytes.NewBuffer(asmFile)
	err := VMConvert(paths, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertFibonacciElement(t *testing.T) {
	execTest(t, "FibonacciElement")
}

func TestConvertNestedCall(t *testing.T) {
	execTest(t, "NestedCall")
}

func TestConvertStaticsTest(t *testing.T) {
	execTest(t, "StaticsTest")
}

package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/verybigtuple/hackassembler/code"
)

func TestReadAsmCodeLabels(t *testing.T) {
	asm := `
		(L0)
		@1 		// 0
		(L1)
		D=-1 	// 1
		M=D		// 2
		(L2)

		// Comment
		@1		// 3
	`
	want := map[string]int{"L0": 0, "L1": 1, "L2": 3}

	symTable := code.NewSymbolTable()
	reader := bufio.NewReader(strings.NewReader(asm))
	_, err := readAsmCode(reader, symTable)
	if err != nil {
		t.Errorf("Error from function: %v", err)
		return
	}

	for name, wantAddr := range want {
		actual, err := symTable.Get(name)
		if err != nil {
			t.Errorf("Labed %s not in symtable", name)
			return
		}
		if actual != wantAddr {
			t.Errorf("Address of %s:%v; want %v", name, actual, wantAddr)
		}
	}
}

func TestReadAsmCodeResult(t *testing.T) {
	asm := `
		(L0)
		@1
		(L1)
		D=-1
		M=D
		(L2)

		// Comment
		@1
	`

	want := []string{
		"@1",
		"D=-1",
		"M=D",
		"@1",
	}

	symTable := code.NewSymbolTable()
	reader := bufio.NewReader(strings.NewReader(asm))
	actual, err := readAsmCode(reader, symTable)
	if err != nil {
		t.Errorf("Error from function: %v", err)
		return
	}
	if len(actual) != len(want) {
		t.Errorf("Actual len: %v; want len: %v", len(actual), len(want))
		return
	}
	for i, actualLine := range actual {
		if actualLine != want[i] {
			t.Errorf("Line %v, Actual %v; want %v", i, actualLine, want[i])
			return
		}
	}
}

func TestRun(t *testing.T) {
	asm := `
		@R0
		D=M              // D = first number
		@R1
		D=D-M            // D = first number - second number
		@OUTPUT_FIRST
		D;JGT            // if D>0 (first is greater) goto output_first
		@R1
		D=M              // D = second number
		@OUTPUT_D
		0;JMP            // goto output_d
	(OUTPUT_FIRST)
		@R0             
		D=M              // D = first number
	(OUTPUT_D)
		@R2
		M=D              // M[2] = D (greatest number)
	(INFINITE_LOOP)
		@INFINITE_LOOP
		0;JMP
	`
	want := [...]string{
		"0000000000000000",
		"1111110000010000",
		"0000000000000001",
		"1111010011010000",
		"0000000000001010",
		"1110001100000001",
		"0000000000000001",
		"1111110000010000",
		"0000000000001100",
		"1110101010000111",
		"0000000000000000",
		"1111110000010000",
		"0000000000000010",
		"1110001100001000",
		"0000000000001110",
		"1110101010000111",
	}

	reader := bufio.NewReader(strings.NewReader(asm))
	sb := strings.Builder{}
	writer := bufio.NewWriter(&sb)

	err := run(reader, writer)
	if err != nil {
		t.Errorf("Error from function: %v", err)
		return
	}

	actual := strings.Split(sb.String(), "\n")
	// remove the last elemet as it is may be an emty string
	le := len(actual) - 1
	if len(actual) > 0 && actual[le] == "" {
		actual = actual[:le]
	}

	if len(actual) != len(want) {
		t.Errorf("Actual len: %v; want len: %v", len(actual), len(want))
		return
	}
	for i, actualLine := range actual {
		if actualLine != want[i] {
			t.Errorf("Line %v, Actual %v; want %v", i, actualLine, want[i])
			return
		}
	}
}

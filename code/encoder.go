package code

import (
	"fmt"
	"strconv"

	"github.com/verybigtuple/hackassembler/parser"
)

const (
	minInt = 0
	maxInt = 32767

	ainstrPrefix = "0"
	cinstrPrefix = "111"
)

var destTable = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var jmpTable = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

var cmpTable = map[string]string{
	"0":   "0" + "101" + "010",
	"1":   "0" + "111" + "111",
	"-1":  "0" + "111" + "010",
	"D":   "0" + "001" + "100",
	"A":   "0" + "110" + "000",
	"M":   "1" + "110" + "000",
	"!D":  "0" + "001" + "101",
	"!A":  "0" + "110" + "001",
	"!M":  "1" + "110" + "001",
	"-D":  "0" + "001" + "111",
	"-A":  "0" + "110" + "011",
	"-M":  "1" + "110" + "011",
	"D+1": "0" + "011" + "111",
	"A+1": "0" + "110" + "111",
	"M+1": "1" + "110" + "111",
	"D-1": "0" + "001" + "110",
	"A-1": "0" + "110" + "010",
	"M-1": "1" + "110" + "010",
	"D+A": "0" + "000" + "010",
	"D+M": "1" + "000" + "010",
	"D-A": "0" + "010" + "011",
	"D-M": "1" + "010" + "011",
	"A-D": "0" + "000" + "111",
	"M-D": "1" + "000" + "111",
	"D&A": "0" + "000" + "000",
	"D&M": "1" + "000" + "000",
	"D|A": "0" + "010" + "101",
	"D|M": "1" + "010" + "101",
}

// EncodeNumber returns 15 bit as string. If argument is less than zero
// or more than 2^15-1, it returns error
func EncodeNumber(n int) (string, error) {
	if n < minInt || n > maxInt {
		return "", &EncoderError{Msg: fmt.Sprintf("Cannot decode %v as it is out of bound", n)}
	}
	return fmt.Sprintf("%015b", n), nil
}

func EncodeAInstr(ai parser.AInstruction, st *SymbolTable) (string, error) {
	var n int
	var err error

	getAddr := func(f func(string) (int, error)) {
		var addr int
		if err != nil {
			return
		}
		addr, err = f(ai.Value)
		n = addr
	}

	if ai.IsVar {
		if !st.Exists(ai.Value) {
			getAddr(st.AddVar)
		} else {
			getAddr(st.Get)
		}
	} else {
		getAddr(strconv.Atoi)
	}

	if err != nil {
		return "", err
	}

	sn, err := EncodeNumber(n)
	if err != nil {
		return "", err
	}
	return ainstrPrefix + sn, nil
}

func EncodeCInstr(ci parser.CIntstruction) (string, error) {
	var err error
	var encDest, encComp, encJmp string

	encode := func(tbl map[string]string, val string, dest *string) {
		if err != nil {
			return
		}
		if _, ok := tbl[val]; !ok {
			err = &EncoderError{Msg: fmt.Sprintf("Cannot encode %v", val)}
			return
		}
		*dest = tbl[val]
	}

	encode(destTable, ci.Dest, &encDest)
	encode(cmpTable, ci.Comp, &encComp)
	encode(jmpTable, ci.Jump, &encJmp)
	if err != nil {
		return "", err
	}

	return cinstrPrefix + encComp + encDest + encJmp, nil
}

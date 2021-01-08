package code

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/verybigtuple/hackassembler/parser"
)

func TestCmpTableUnique(t *testing.T) {
	m := map[string]int{}
	for k, v := range cmpTable {
		if _, ok := m[v]; !ok {
			m[v] = 0
		}
		m[v]++
		if m[v] > 1 {
			t.Errorf("Value '%s' is not unique for key '%s'", v, k)
			return
		}
	}
}

func TestCmpTableAMBit(t *testing.T) {
	for k, v := range cmpTable {
		switch {
		case strings.Contains(k, "A"):
			if v[0] != '0' {
				t.Errorf("Value '%s' for key '%s' has wrong first bit", v, k)
				return
			}
		case strings.Contains(k, "M"):
			if v[0] != '1' {
				t.Errorf("Value '%s' for key '%s' has wrong first bit", v, k)
				return
			}
		}
	}
}

func remSp(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func TestDecodeNumberRegular(t *testing.T) {
	testCases := []struct {
		arg  int
		want string
	}{
		{
			arg:  0,
			want: remSp("000 0000 0000 0000"),
		},
		{
			arg:  1,
			want: remSp("000 0000 0000 0001"),
		},
		{
			arg:  2000,
			want: remSp("000 0111 1101 0000"),
		},
		{
			arg:  2000,
			want: remSp("000 0111 1101 0000"),
		},
		{
			arg:  32767,
			want: remSp("111 1111 1111 1111"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Arg %v", tc.arg), func(t *testing.T) {
			actual, err := EncodeNumber(tc.arg)
			if err != nil {
				t.Errorf("EncodeNumber returned unexpected error: %v", err)
				return
			}
			if actual != tc.want {
				t.Errorf("Actual: %v, want: %v", actual, tc.want)
			}
		})
	}
}

func TestEncodeNumberError(t *testing.T) {
	testCases := []int{
		-1,
		-15000,
		32768,
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Arg %v", tc), func(t *testing.T) {
			actual, err := EncodeNumber(tc)
			de := &EncoderError{}
			if err != nil && !errors.As(err, &de) {
				t.Errorf("EncodeNumber retuned unexpected error type: %v", err)
				return
			}

			if err == nil {
				t.Errorf("EncodeNumber did not returned an error: %v", actual)
			}
		})
	}
}

func TestEncodeAInstr(t *testing.T) {
	symTable := NewSymbolTable()

	testCases := []struct {
		instr parser.AInstruction
		want  string
	}{
		{
			instr: parser.AInstruction{IsVar: false, Value: "1"},
			want:  remSp("0000 0000 0000 0001"),
		},
		{
			instr: parser.AInstruction{IsVar: false, Value: "001"},
			want:  remSp("0000 0000 0000 0001"),
		},
		{
			instr: parser.AInstruction{IsVar: true, Value: "R15"},
			want:  remSp("0000 0000 0000 1111"),
		},
		{
			instr: parser.AInstruction{IsVar: true, Value: "i1"},
			want:  remSp("0000 0000 0001 0000"),
		},
		{
			instr: parser.AInstruction{IsVar: true, Value: "i2"},
			want:  remSp("0000 0000 0001 0001"),
		},
		// Repeat the same instruction on purpose
		{
			instr: parser.AInstruction{IsVar: true, Value: "i2"},
			want:  remSp("0000 0000 0001 0001"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.instr), func(t *testing.T) {
			actual, err := EncodeAInstr(tc.instr, symTable)
			if err != nil {
				t.Errorf("DecodeAInstr returned unexpected error: %v", err)
				return
			}
			if actual != tc.want {
				t.Errorf("Actual: %v, want: %v", actual, tc.want)
			}
		})
	}
}

func TestEncodeAInstrError(t *testing.T) {
	symTable := NewSymbolTable()

	testCases := []parser.AInstruction{
		{IsVar: false, Value: "i1"},    // not digit value
		{IsVar: false, Value: "65536"}, // value is very big
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			actual, err := EncodeAInstr(tc, symTable)
			if err == nil {
				t.Errorf("EncodeAInstr did not returned an error: %v", actual)
			}
		})
	}
}

func TestEncodeCInstr(t *testing.T) {
	testCases := []struct {
		instr parser.CIntstruction
		want  string
	}{
		{
			instr: parser.CIntstruction{Dest: "M", Comp: "1"},
			want:  remSp("1110 1111 1100 1000"),
		},
		{
			instr: parser.CIntstruction{Dest: "M", Comp: "D+M"},
			want:  remSp("1111 0000 1000 1000"),
		},
		{
			instr: parser.CIntstruction{Comp: "0"},
			want:  remSp("1110 1010 1000 0000"),
		},
		{
			instr: parser.CIntstruction{Comp: "0", Jump: "JMP"},
			want:  remSp("1110 1010 1000 0111"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.instr), func(t *testing.T) {
			actual, err := EncodeCInstr(tc.instr)
			if err != nil {
				t.Errorf("EncodeCInstr returned unexpected error: %v", err)
				return
			}
			if actual != tc.want {
				t.Errorf("Actual: %v, want: %v", actual, tc.want)
			}
		})
	}
}

func TestEncodeCInstrError(t *testing.T) {
	testCases := []parser.CIntstruction{
		{Dest: "MDA", Comp: "0"},
		{Dest: "D", Comp: "D+1", Jump: "JJJ"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			actual, err := EncodeCInstr(tc)
			if err == nil {
				t.Errorf("EncodeCInstr did not returned an error: %v", actual)
			}
		})
	}
}

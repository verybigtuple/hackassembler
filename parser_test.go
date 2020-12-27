package main

import (
	"testing"
)

type testCase struct {
	operator string
	want     CIntstruction
}

func (tC *testCase) Run(t *testing.T) {
	t.Run(tC.operator, func(t *testing.T) {
		actual, err := ParseCInstrunction(tC.operator)
		if err != nil {
			t.Errorf("Exception from function: %w", err)
		}
		if *actual != tC.want {
			t.Errorf("Parsed: %+v ; want %+v", *actual, tC.want)
		}
	})
}

func TestParseCInstructionRegular(t *testing.T) {
	testCases := []testCase{
		{
			operator: "0",
			want:     CIntstruction{Comp: "0"},
		},
		{
			operator: "M&D",
			want:     CIntstruction{Comp: "M&D"},
		},
		{
			operator: "A=M|D",
			want:     CIntstruction{Dest: "A", Comp: "M|D"},
		},
		{
			operator: "M=D",
			want:     CIntstruction{Dest: "M", Comp: "D"},
		},
		{
			operator: "0;JMP",
			want:     CIntstruction{Comp: "0", Jump: "JMP"},
		},
		{
			operator: "M=D+1;JEQ",
			want:     CIntstruction{Dest: "M", Comp: "D+1", Jump: "JEQ"},
		},
		{
			operator: "AMD=-M;JEQ",
			want:     CIntstruction{Dest: "AMD", Comp: "-M", Jump: "JEQ"},
		},
	}

	for _, tC := range testCases {
		tC.Run(t)
	}
}

func TestParseCInstructionSpaces(t *testing.T) {
	testCases := []testCase{
		{
			operator: "   D",
			want:     CIntstruction{Comp: "D"},
		},
		{
			operator: "\t\tD\t ",
			want:     CIntstruction{Comp: "D"},
		},

		{
			operator: " D = D + 1 ",
			want:     CIntstruction{Dest: "D", Comp: "D+1"},
		},
		{
			operator: " D = D + 1 ; JMP ",
			want:     CIntstruction{Dest: "D", Comp: "D+1", Jump: "JMP"},
		},
		{
			operator: " A  D = D + 1 ; JMP ",
			want:     CIntstruction{Dest: "AD", Comp: "D+1", Jump: "JMP"},
		},
	}

	for _, tC := range testCases {
		tC.Run(t)
	}
}

func TestParseCInstructionComment(t *testing.T) {
	testCases := []testCase{
		{
			operator: "D \t // Comment ",
			want:     CIntstruction{Comp: "D"},
		},
		{
			operator: "D=A+D  //Long Comment",
			want:     CIntstruction{Dest: "D", Comp: "A+D"},
		},
		{
			operator: "D=D+1;JMP // Comment ",
			want:     CIntstruction{Dest: "D", Comp: "D+1", Jump: "JMP"},
		},
	}

	for _, tC := range testCases {
		tC.Run(t)
	}
}

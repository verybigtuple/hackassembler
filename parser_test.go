package main

import (
	"testing"
)

func TestParseCInstructionRegular(t *testing.T) {
	testCases := []struct {
		operator string
		want     CIntstruction
	}{
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
		t.Run(tC.operator, func(t *testing.T) {
			actual := ParseCInstrunction(tC.operator)
			if *actual != tC.want {
				t.Errorf("Parsed: %+v ; want %+v", *actual, tC.want)
			}
		})
	}
}

func TestParseCInstructionSpaces(t *testing.T) {
	testCases := []struct {
		operator string
		want     CIntstruction
	}{
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
		t.Run(tC.operator, func(t *testing.T) {
			actual := ParseCInstrunction(tC.operator)
			if *actual != tC.want {
				t.Errorf("Parsed: %+v ; want %+v", *actual, tC.want)
			}
		})
	}
}

func TestParseCInstructionComment(t *testing.T) {
	testCases := []struct {
		operator string
		want     CIntstruction
	}{
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
		t.Run(tC.operator, func(t *testing.T) {
			actual := ParseCInstrunction(tC.operator)

			if *actual != tC.want {
				t.Errorf("Parsed: %+v ; want %+v", *actual, tC.want)
			}
		})
	}
}

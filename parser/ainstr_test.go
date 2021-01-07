package parser

import (
	"errors"
	"testing"
)

type ainstrTestCase struct {
	operator string
	want     AInstruction
}

func (tC *ainstrTestCase) run(p *AParser, t *testing.T) {
	t.Run(tC.operator, func(t *testing.T) {
		actual, err := p.Parse(tC.operator)
		if err != nil {
			t.Errorf("The test returned an exception: %v", err)
			return
		}
		if *actual != tC.want {
			t.Errorf("Parsed: %+v ; want %+v", *actual, tC.want)
		}
	})
}

func (tC *ainstrTestCase) runParseError(p *AParser, t *testing.T) {
	t.Run(tC.operator, func(t *testing.T) {
		actual, err := p.Parse(tC.operator)
		if err == nil {
			t.Errorf("Error was not arisen as expected. Actual: %+v", *actual)
			return
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Error was arisen, but it is not ParseError: %v", err)
		}
	})
}

func TestAParseRegular(t *testing.T) {
	testCases := []ainstrTestCase{
		{
			operator: "@0",
			want:     AInstruction{IsVar: false, Value: "0"},
		},
		{
			operator: "  @0  ",
			want:     AInstruction{IsVar: false, Value: "0"},
		},
		{
			operator: "@123",
			want:     AInstruction{IsVar: false, Value: "123"},
		},
		{
			operator: "@0123",
			want:     AInstruction{IsVar: false, Value: "0123"},
		},
		{
			operator: "@i",
			want:     AInstruction{IsVar: true, Value: "i"},
		},
		{
			operator: "@variable1",
			want:     AInstruction{IsVar: true, Value: "variable1"},
		},
		{
			operator: "@v123",
			want:     AInstruction{IsVar: true, Value: "v123"},
		},
		{
			operator: "@var //Comment",
			want:     AInstruction{IsVar: true, Value: "var"},
		},
	}

	p := NewAParser()
	for _, tC := range testCases {
		tC.run(p, t)
	}
}

func TestAParseError(t *testing.T) {
	testCases := []ainstrTestCase{
		{
			operator: "A",
		},
		{
			operator: "@",
		},
		{
			operator: "@@0",
		},
		{
			operator: "@0@",
		},
		{
			operator: "@0DD",
		},
		{
			operator: "@123//C",
		},
		{
			operator: "@123/",
		},
		{
			operator: "@A B",
		},
		{
			operator: "@A_B",
		},
		{
			operator: "@A /B",
		},
		{
			operator: "@A-B",
		},
		{
			operator: "@A /_",
		},
	}

	p := NewAParser()
	for _, tC := range testCases {
		tC.runParseError(p, t)
	}
}

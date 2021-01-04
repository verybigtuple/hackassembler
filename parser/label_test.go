package parser

import (
	"errors"
	"testing"
)

type labelTestCase struct {
	operator string
	want     Label
}

func (tC *labelTestCase) run(p *LabelParser, t *testing.T) {
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

func (tC *labelTestCase) runParseError(p *LabelParser, t *testing.T) {
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

func TestParseLabelRegular(t *testing.T) {
	testCases := []labelTestCase{
		{
			operator: "(LABEL)",
			want:     Label("LABEL"),
		},
		{
			operator: " (LABEL)  ",
			want:     Label("LABEL"),
		},
		{
			operator: "(L1)",
			want:     Label("L1"),
		},
		{
			operator: "(lab1)",
			want:     Label("lab1"),
		},

		{
			operator: "(lab1) //Comment",
			want:     Label("lab1"),
		},
	}

	p := NewLabelParser()
	for _, tC := range testCases {
		tC.run(p, t)
	}
}

func TestParseLabelSyntaxError(t *testing.T) {
	testCases := []labelTestCase{
		{
			operator: "(Label",
		},
		{
			operator: "( Label)",
		},
		{
			operator: "(Lab el)",
		},
		{
			operator: "(Label )",
		},
		{
			operator: "(Label(",
		},
		{
			operator: "(1Label)",
		},
		{
			operator: "(Label) a",
		},
		{
			operator: "(Label) /a",
		},
	}

	p := NewLabelParser()
	for _, tC := range testCases {
		tC.runParseError(p, t)
	}
}

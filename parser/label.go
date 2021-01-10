package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Label contains a single string that represents HackAssembler Label
type Label string

// LabelParser for Label
type LabelParser struct {
	label    Label
	reader   *pRuneReader
	nextStep func() error
	strB     strings.Builder
}

// NewLabelParser returns a pointer to a new LabelParser
func NewLabelParser() *LabelParser {
	lp := LabelParser{strB: strings.Builder{}}
	lp.strB.Grow(15)
	return &lp
}

// Parse returns a label parsed from Label or an error
func (p *LabelParser) Parse(str string) (*Label, error) {
	p.reader = newPRuneReader(str)
	p.strB.Reset()

	p.nextStep = p.checkStart
	for err := p.nextStep(); !errors.Is(err, errEOP); err = p.nextStep() {
		if err != nil {
			return nil, err
		}
	}

	return &p.label, nil
}

func (p *LabelParser) checkStart() error {
	rv, _, err := p.reader.ReadAfterSpaces()
	if err != nil {
		return errEOP
	}
	if rv != startLabelLiteral {
		return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected start of label '%c'", rv)}
	}

	p.nextStep = p.checkFirst
	return nil
}

func (p *LabelParser) checkFirst() error {
	rv, _, err := p.reader.ReadRune()
	if err != nil {
		return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Label should finish with '%v'", endLabelLiteral)}
	}

	if !unicode.IsLetter(rv) && rv != '_' {
		return &ParseError{
			Pos: p.reader.Pos,
			Msg: fmt.Sprintf("Unexpected character '%v': Label must begin with a letter", rv),
		}
	}
	p.strB.WriteRune(rv)
	p.nextStep = p.readRest
	return nil
}

func (p *LabelParser) readRest() error {
	for {
		rv, _, err := p.reader.ReadRune()
		if err != nil {
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Label should finish with '%v'", endLabelLiteral)}
		}

		if rv == endLabelLiteral {
			p.label = Label(p.strB.String())
			p.nextStep = p.checkTail
			return nil
		}

		if unicode.IsSpace(rv) {
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Label cannot has spaces")}
		}

		p.strB.WriteRune(rv)
	}
}

func (p *LabelParser) checkTail() error {
	return checkInlineComment(p.reader)
}

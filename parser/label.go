package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	startLabel = '('
	endLabel   = ')'
)

var EOP error = errors.New("EOP")

type Label string

type LabelParser struct {
	label    Label
	reader   *pRuneReader
	nextStep func() error
	strB     strings.Builder
}

func NewLabelParser() *LabelParser {
	lp := LabelParser{strB: strings.Builder{}}
	lp.strB.Grow(15)
	return &lp
}

func (p *LabelParser) Parse(str string) (*Label, error) {
	p.reader = newPRuneReader(str)
	p.strB.Reset()

	p.nextStep = p.checkStart
	for err := p.nextStep(); !errors.Is(err, EOP); err = p.nextStep() {
		if err != nil {
			return nil, err
		}
	}

	return &p.label, nil
}

func (p *LabelParser) checkStart() error {
	rv, _, err := p.reader.ReadAfterSpaces()
	if err != nil {
		return EOP
	}
	if rv != startLabel {
		return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected start of label '%c'", rv)}
	}

	p.nextStep = p.checkFirst
	return nil
}

func (p *LabelParser) checkFirst() error {
	rv, _, err := p.reader.ReadRune()
	if err != nil {
		return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Label should finish with '%v'", endLabel)}
	}

	if !unicode.IsLetter(rv) {
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
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Label should finish with '%v'", endLabel)}
		}

		if rv == endLabel {
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
	for {
		rv, _, err := p.reader.ReadRune()
		if err != nil {
			return EOP
		}

		switch {
		case rv == '/':
			nrv, _, err := p.reader.ReadRune()
			if err != nil {
				return &ParseError{Pos: p.reader.Pos, Msg: "Unexpected character /"}
			}
			if nrv != '/' {
				return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected character '%v'", rv)}
			}
			return EOP
		case !unicode.IsSpace(rv):
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected character '%v'", rv)}
		}
	}
}

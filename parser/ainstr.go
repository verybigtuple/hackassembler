package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const startAInstr = '@'

//AInstruction of HackAssembler
type AInstruction struct {
	IsVar bool
	Value string
}

//AParser is Paraser for A-Instructions
type AParser struct {
	aInstr   AInstruction
	reader   *pRuneReader
	nextStep func() error
	strB     strings.Builder
}

//NewAParser returns a pointer to a created AParser
func NewAParser() *AParser {
	ap := AParser{}
	ap.strB.Grow(15)
	return &ap
}

//Parse function returns parsed AInstruction or error
func (p *AParser) Parse(s string) (*AInstruction, error) {
	p.reader = newPRuneReader(s)
	p.strB.Reset()

	p.nextStep = p.checkStart
	for err := p.nextStep(); !errors.Is(err, errEOP); err = p.nextStep() {
		if err != nil {
			return nil, err
		}
	}

	return &p.aInstr, nil
}

func (p *AParser) checkStart() error {
	rv, _, err := p.reader.ReadAfterSpaces()
	if err != nil {
		return errEOP
	}

	if rv != startAInstr {
		return &ParseError{Pos: p.reader.Pos, Msg: "Unexpected start of A-Intruction"}
	}

	p.nextStep = p.checkFirst
	return nil
}

func (p *AParser) checkFirst() error {
	rv, _, err := p.reader.ReadRune()
	if err != nil {
		return &ParseError{Pos: p.reader.Pos, Msg: "A-Instruction ends unexpectedly"}
	}

	if !unicode.IsDigit(rv) && !unicode.IsLetter(rv) && rv != '_' {
		return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected first symbol '%c'", rv)}
	}

	if unicode.IsDigit(rv) {
		p.nextStep = p.readNumber
	} else {
		p.nextStep = p.readVar
	}

	p.strB.WriteRune(rv)
	return nil
}

func (p *AParser) readNumber() error {
	var e error = nil
	for {
		rv, _, err := p.reader.ReadRune()
		if err != nil {
			e = errEOP
			break
		}
		if unicode.IsSpace(rv) {
			p.nextStep = p.readComment
			break
		}
		if !unicode.IsDigit(rv) {
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected character in A-Instruction '%c'", rv)}
		}
		p.strB.WriteRune(rv)
	}
	p.aInstr.IsVar = false
	p.aInstr.Value = p.strB.String()
	return e
}

func (p *AParser) readVar() error {
	var e error = nil
	for {
		rv, _, err := p.reader.ReadRune()
		if err != nil {
			e = errEOP
			break
		}
		if unicode.IsSpace(rv) {
			p.nextStep = p.readComment
			break
		}
		if !unicode.IsLetter(rv) && !unicode.IsDigit(rv) && rv != '_' {
			return &ParseError{Pos: p.reader.Pos, Msg: fmt.Sprintf("Unexpected character in A-Instruction '%c'", rv)}
		}
		p.strB.WriteRune(rv)
	}
	p.aInstr.IsVar = true
	p.aInstr.Value = p.strB.String()
	return e
}

func (p *AParser) readComment() error {
	return checkInlineComment(p.reader)
}

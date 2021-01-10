package parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// CIntstruction parsed into 3 parts
type CIntstruction struct {
	Dest string
	Comp string
	Jump string
}

type runeReadSeeker interface {
	io.RuneReader
	io.Seeker
}

//CParser for parsing C-Instructions. Implements Prser intrface
type CParser struct {
	cInstr     CIntstruction
	rReader    runeReadSeeker
	cmdBuilder strings.Builder
	pos        int
}

//NewCParser returns ptr to a new CParser
func NewCParser() *CParser {
	p := CParser{cmdBuilder: strings.Builder{}}
	p.cmdBuilder.Grow(3)
	return &p
}

//Parse reurns parsed CInstruction from a line of code or an error
func (p *CParser) Parse(str string) (*CIntstruction, error) {
	p.cInstr = CIntstruction{}
	p.rReader = strings.NewReader(str)
	p.cmdBuilder.Reset()
	p.pos = 0

	setter := p.setComp
	for {
		rv, _, err := p.rReader.ReadRune()
		if err != nil {
			break
		}
		p.pos++

		switch {
		case rv == compDelim:
			err = p.setDest()
		case rv == jumpDelim:
			err = p.setComp()
			setter = p.setJump
		case rv == commentLiteral:
			err = p.parseComment()
		case !unicode.IsSpace(rv):
			p.cmdBuilder.WriteRune(rv)
		}

		if err != nil {
			return nil, err
		}

	}

	err := setter()
	if err != nil {
		return nil, err
	}

	return &p.cInstr, nil
}

func (p *CParser) setDest() error {
	if p.cmdBuilder.Len() == 0 {
		return &ParseError{Pos: p.pos, Msg: fmt.Sprintf("Dest must be set up before '%v'", compDelim)}
	}
	p.cInstr.Dest = p.cmdBuilder.String()
	p.cmdBuilder.Reset()
	return nil
}

func (p *CParser) setComp() error {
	if p.cmdBuilder.Len() == 0 && p.cInstr.Dest != "" {
		return &ParseError{Pos: p.pos, Msg: "Computation operator absent after Destination"}
	}
	p.cInstr.Comp = p.cmdBuilder.String()
	p.cmdBuilder.Reset()
	return nil
}

func (p *CParser) setJump() error {
	if p.cmdBuilder.Len() == 0 {
		return &ParseError{Pos: p.pos, Msg: fmt.Sprintf("Jump must be be set up after '%v'", jumpDelim)}
	}
	if p.cInstr.Comp == "" {
		return &ParseError{Pos: p.pos, Msg: "Computation absent before Jump"}
	}
	p.cInstr.Jump = p.cmdBuilder.String()
	p.cmdBuilder.Reset()
	return nil
}

func (p *CParser) parseComment() error {
	// checking the next slash
	nextSlash, _, err := p.rReader.ReadRune()
	if nextSlash != commentLiteral || err != nil {
		return &ParseError{Pos: p.pos, Msg: "Expected '/' for the inline comment"}
	}
	// skip the rest of string to the end, i.e. we just ignore the comment
	p.rReader.Seek(0, io.SeekEnd)
	return nil
}

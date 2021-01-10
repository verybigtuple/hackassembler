package parser

import (
	"fmt"
)

const (
	ainstrLiteral     = '@'
	startLabelLiteral = '('
	endLabelLiteral   = ')'

	compDelim = '='
	jumpDelim = ';'

	commentPrefix  = "//"
	commentLiteral = '/'
)

//Parser interface
type Parser interface {
	Parse(s string) (*interface{}, error)
}

func checkInlineComment(r *pRuneReader) error {
	rv, _, err := r.ReadAfterSpaces()
	if err != nil {
		return errEOP
	}

	if rv != commentLiteral {
		return &ParseError{Pos: r.Pos, Msg: fmt.Sprintf("Unexpected character '%c'", rv)}
	}

	nrv, _, err := r.ReadRune() // read next rune
	if err != nil {
		return &ParseError{Pos: r.Pos, Msg: fmt.Sprintf("Unexpected end of line after '%c'", rv)}
	}
	if nrv != commentLiteral {
		return &ParseError{Pos: r.Pos, Msg: fmt.Sprintf("Unexpected character '%c'", rv)}
	}
	return errEOP
}

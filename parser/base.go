package parser

import (
	"errors"
	"fmt"
)

const commentLiteral = '/'

//Parser interface
type Parser interface {
	Parse(s string) (*interface{}, error)
}

//ParseError implements error arisen while parsing
type ParseError struct {
	Pos int
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parsing error at position %d: %s", e.Pos, e.Msg)
}

//Special error to stop state machine
var errEOP error = errors.New("EOP")

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

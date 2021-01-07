package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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

type pRuneReader struct {
	reader *strings.Reader
	Pos    int
}

func newPRuneReader(s string) *pRuneReader {
	rr := pRuneReader{reader: strings.NewReader(s)}
	return &rr
}

func (rR *pRuneReader) ReadRune() (rune, int, error) {
	r, s, err := rR.reader.ReadRune()
	if err == nil {
		rR.Pos++
	}
	return r, s, err
}

func (rR *pRuneReader) ReadAfterSpaces() (rune, int, error) {
	for {
		rv, s, err := rR.ReadRune()
		if !unicode.IsSpace(rv) || err != nil {
			return rv, s, err
		}
	}
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

package parser

import (
	"fmt"
	"strings"
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

// IsLabelLine returns true if line starts with a label prefix
func IsLabelLine(line string) bool {
	return strings.HasPrefix(line, string(startLabelLiteral))
}

// IsAInstrLine returns true if line starts with an A-Instruction prefix
func IsAInstrLine(line string) bool {
	return strings.HasPrefix(line, string(ainstrLiteral))
}

// IsCommentLine returns true if line starts with a comment Prefix
func IsCommentLine(line string) bool {
	return strings.HasPrefix(line, commentPrefix)
}

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

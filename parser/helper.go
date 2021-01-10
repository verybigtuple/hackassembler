package parser

import (
	"fmt"
	"strings"
	"unicode"
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

// Symbols that are ok for Label/A-Instruction body,
// i.e @a.name or (LABEL_A$SOME) are ok
var varSymbolSpecCase = map[rune]bool{'_': true, '.': true, '$': true}

// isVarRune returns true if rune can be a rune of a variable name or a lebel name.
// It cannot be applicable for the first rune
func isVarRune(r rune) bool {
	_, ok := varSymbolSpecCase[r]
	return unicode.IsLetter(r) || unicode.IsDigit(r) || ok
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

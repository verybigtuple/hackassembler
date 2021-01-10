package parser

import (
	"errors"
	"fmt"
)

// ParseError implements error arisen while parsing Label, A- or C-Intruction.
type ParseError struct {
	Pos int
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parsing error at position %d: %s", e.Pos, e.Msg)
}

// Special error to stop state machine
var errEOP error = errors.New("EOP")

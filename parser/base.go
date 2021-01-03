package parser

import "fmt"

//Parser interface
type Parser interface {
	Parse(s string) (interface{}, error)
}

//ParseError implements error arisen while parsing
type ParseError struct {
	Pos int
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Paring error at position %d: %s", e.Pos, e.Msg)
}

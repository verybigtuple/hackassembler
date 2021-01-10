package parser

import (
	"strings"
	"unicode"
)

// pRuneReader is struct that has a regular reader. Pos is the current count of symbols:
// the first rune is 1.
type pRuneReader struct {
	reader *strings.Reader
	Pos    int
}

func newPRuneReader(s string) *pRuneReader {
	rr := pRuneReader{reader: strings.NewReader(s)}
	return &rr
}

// ReadRune returns the next rune
func (rR *pRuneReader) ReadRune() (rune, int, error) {
	r, s, err := rR.reader.ReadRune()
	if err == nil {
		rR.Pos++
	}
	return r, s, err
}

// ReadAfterSpaces returns the next rune after all spaces
func (rR *pRuneReader) ReadAfterSpaces() (rune, int, error) {
	for {
		rv, s, err := rR.ReadRune()
		if !unicode.IsSpace(rv) || err != nil {
			return rv, s, err
		}
	}
}

package main

import (
	"strings"
	"unicode"
)

const (
	compDelim          = '='
	jumpDelim          = ';'
	inlineCommentDelim = '/'
)

// CIntstruction parsed into 3 parts
type CIntstruction struct {
	Dest string
	Comp string
	Jump string
}

// ParseCInstrunction returns pointer to a CInstruction struct parsed from a string line
func ParseCInstrunction(input string) *CIntstruction {
	ci := CIntstruction{}
	b := strings.Builder{}
	b.Grow(3)

	expectedPart := &ci.Comp
Loop:
	for _, runeVal := range input {
		switch {
		case runeVal == compDelim:
			ci.Dest = b.String()
			b.Reset()
		case runeVal == jumpDelim:
			ci.Comp = b.String()
			expectedPart = &ci.Jump
			b.Reset()
		case runeVal == inlineCommentDelim:
			break Loop
		case !unicode.IsSpace(runeVal):
			b.WriteRune(runeVal)
		}
	}

	if b.Len() > 0 {
		*expectedPart = b.String()
	}
	return &ci
}

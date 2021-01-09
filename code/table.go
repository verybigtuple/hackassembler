package code

import (
	"fmt"
)

const (
	minUserRAM = 16
	maxUserRAM = 16383
	minROM     = 0
	maxROM     = 32767
)

// Default vars according to the language specification
var initTable = map[string]int{
	// Virtual Machine
	"SP":   0,
	"LCL":  1,
	"ARG":  2,
	"THIS": 3,
	"THAT": 4,

	// R registers
	"R0": 0, "R1": 1, "R2": 2, "R3": 3, "R4": 4,
	"R5": 5, "R6": 6, "R7": 7, "R8": 8, "R9": 9,
	"R10": 10, "R11": 11, "R12": 12, "R13": 13, "R14": 14, "R15": 15,

	// Inputs
	"SCREEN": 16384,
	"KBD":    24576,
}

// SymbolTable is register for ROM labels and RAM variables
type SymbolTable struct {
	Table        map[string]int
	UserRegister int
}

// NewSymbolTable creates a new SymbolTable and init it with predefined vars
func NewSymbolTable() *SymbolTable {
	// Copying initTable
	newTable := make(map[string]int, len(initTable))
	for k, v := range initTable {
		newTable[k] = v
	}

	vt := SymbolTable{Table: newTable, UserRegister: minUserRAM}
	return &vt
}

// Exists returns true if label or var is in the symbol table
func (t *SymbolTable) Exists(name string) bool {
	_, ok := t.Table[name]
	return ok
}

// AddVar adds a new user var. The address if the var is added automatically and is returned
// from the function. If the Var already exists, the error will be returned
func (t *SymbolTable) AddVar(name string) (int, error) {
	if t.UserRegister > maxUserRAM {
		return 0, &EncoderError{Msg: fmt.Sprintf("User RAM ran out. Address %v is reserved", t.UserRegister)}
	}
	if t.Exists(name) {
		return 0, &EncoderError{Msg: fmt.Sprintf("Variable '%v' already exists", name)}
	}
	t.Table[name] = t.UserRegister
	t.UserRegister++
	return t.Table[name], nil
}

// Get returns RAM address for a var and ROM address for a label. If var or label
// does not exist in the symbol table, the error will be returned
func (t *SymbolTable) Get(name string) (int, error) {
	if v, ok := t.Table[name]; ok {
		return v, nil
	}
	return 0, &EncoderError{Msg: fmt.Sprintf("Cannot get var or label '%v' as it does not exist", name)}
}

// AddLabel adds a new label with custom integer ROM address 'val'
func (t *SymbolTable) AddLabel(name string, val int) (int, error) {
	if t.Exists(name) {
		return 0, &EncoderError{Msg: fmt.Sprintf("Cannot add label '%v' as it has alredy existed", name)}
	}

	if val < minROM || val > maxROM {
		return 0, &EncoderError{Msg: fmt.Sprintf("Label '%v' has ROM address %v out of bound", name, val)}
	}

	t.Table[name] = val
	return val, nil
}

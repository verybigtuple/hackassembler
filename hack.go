package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/verybigtuple/hackassembler/code"
	"github.com/verybigtuple/hackassembler/parser"
)

const (
	initCodeSize = 1000
)

func clearLine(s string) string {
	return strings.Trim(s, " \t\r\n")
}

type codeReader struct {
	input     *bufio.Reader
	LineCount int
}

func newCodeReader(in *bufio.Reader) *codeReader {
	return &codeReader{input: in}
}

func (r *codeReader) readNextLine() (string, error) {
	line, err := r.input.ReadString('\n')
	if err != nil {
		return "", err
	}
	r.LineCount++
	normLine := strings.Trim(line, " \t\r\n")
	return normLine, nil
}

// readNextCodeLine skips all comment lines and empty line and read next instruction or label line
func (r *codeReader) readNextCodeLine() (string, error) {
	for {
		line, err := r.readNextLine()
		if err != nil {
			return "", err
		}
		if len(line) > 0 && !parser.IsCommentLine(line) {
			return line, nil
		}
	}
}

// readAsmCode reads code lines, adds all labels to Symbol table and retruns all asm lines
// without spaces and comments as []string
func readAsmCode(in *bufio.Reader, st *code.SymbolTable) ([]string, error) {
	asmLines := make([]string, 0, initCodeSize)

	asmReader := newCodeReader(in)
	labelParser := parser.NewLabelParser()

	romCount := 0
	for {
		line, err := asmReader.readNextCodeLine()
		if err != nil {
			return asmLines, nil
		}

		if parser.IsLabelLine(line) {
			label, err := labelParser.Parse(line)
			if err != nil {
				return nil, err
			}
			st.AddLabel(string(*label), romCount)
		} else {
			asmLines = append(asmLines, line)
			// asmLines[romCount] = line
			romCount++
		}
	}
}

func encodeAsm(asmCode []string, st *code.SymbolTable) ([]string, error) {
	encoded := make([]string, 0, len(asmCode))
	aParcer := parser.NewAParser()
	cParser := parser.NewCParser()

	for _, v := range asmCode {
		if parser.IsAInstrLine(v) {
			//TODO: What the hell?
			ai, err := aParcer.Parse(v)
			if err != nil {
				return nil, err
			}
			encA, err := code.EncodeAInstr(*ai, st)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, encA)
		} else {
			ci, err := cParser.Parse(v)
			if err != nil {
				return nil, err
			}
			encC, err := code.EncodeCInstr(*ci)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, encC)
		}
	}
	return encoded, nil
}

func run(in *bufio.Reader, out *bufio.Writer) error {
	symbolTable := code.NewSymbolTable()

	codeLines, err := readAsmCode(in, symbolTable)
	if err != nil {
		return err
	}

	encLines, err := encodeAsm(codeLines, symbolTable)
	if err != nil {
		return err
	}

	for _, cl := range encLines {
		out.WriteString(cl + "\n")
	}
	out.Flush()
	return nil
}

func main() {
	if len(os.Args) < 3 {
		panic("Not enough args")
	}

	inFilePath := os.Args[1]
	outFile := os.Args[2]

	inF, err := os.OpenFile(inFilePath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer inF.Close()
	inReader := bufio.NewReader(inF)

	outF, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer outF.Close()
	outWriter := bufio.NewWriter(outF)

	err = run(inReader, outWriter)
	if err != nil {
		panic(err)
	}
}

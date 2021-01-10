package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/verybigtuple/hackassembler/code"
	"github.com/verybigtuple/hackassembler/parser"
)

const (
	initCodeSize = 1000

	// Exit Codes
	parserError = -1
	codeError   = -2
	fileError   = -3
	otherError  = -99
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
			st.AddLabel(label.Value, romCount)
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
	inFileFlag := flag.String("in", "", "Input file with hack assembler. Usually has extension *.asm")
	outFileFlag := flag.String("out", "", "Output file with binary code. Usually a file *.hack")

	flag.Parse()

	if *inFileFlag == "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Input file is not set", *inFileFlag))
		os.Exit(fileError)
	}

	if *outFileFlag == "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Output file is not set", *inFileFlag))
		os.Exit(fileError)
	}

	if _, err := os.Stat(*inFileFlag); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Input file %s is not found", *inFileFlag))
		os.Exit(fileError)
	}

	inF, err := os.OpenFile(*inFileFlag, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot open input file: %v", err))
		os.Exit(fileError)
	}
	defer inF.Close()

	outParentDir := filepath.Dir(*outFileFlag)
	if err := os.MkdirAll(outParentDir, os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot create dir %s: %v", outParentDir, err))
		os.Exit(fileError)
	}

	outF, err := os.OpenFile(*outFileFlag, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot open output file: %v", err))
		os.Exit(fileError)
	}
	defer outF.Close()

	inReader := bufio.NewReader(inF)
	outWriter := bufio.NewWriter(outF)
	err = run(inReader, outWriter)
	if err != nil {
		switch e := err.(type) {
		case *parser.ParseError:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Parsing Error: %v", e))
			os.Exit(parserError)
		case *code.EncoderError:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Encoding Error: %v", e))
			os.Exit(codeError)
		default:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Unknown Error: %v", e))
			os.Exit(otherError)
		}
	}
}

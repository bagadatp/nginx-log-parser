package parser

import (
	"io"
	"log"
	"os"
	"testing"

	"gitlab.autoiterative.com/group-lucid-swirles-heyrovsky/log-parser/pkg/validator"
)

const (
	ValidInputFile    = "testdata/nginx-valid0.log"
	InvalidInputFile  = "testdata/nginx-invalid0.log"
	DefaultOutputFile = "testdata/output.json"
)

func TestNewDefaultParser(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	validator := validator.NewRegexValidator(logger)
	parser := NewDefaultParser(logger, validator)
	_, ok := interface{}(parser).(Parser)
	if !ok {
		t.Fatalf("Returned Parser object does not implement Parser interface")
	}
	if parser.l != logger {
		t.Fatalf("Returned Parser object does not have the right logger object: %p != %p", parser.l, logger)
	}
	if parser.v != validator {
		t.Fatalf("Returned Parser object does not have the right logger object: %p != %p", parser.v, validator)
	}
}

func TestParseNginxLog(t *testing.T) {
	for _, tc := range []struct {
		name     string
		inF      *os.File
		outF     *os.File
		maxIps   int
		maxPaths int
		exp      error
	}{
		{"nil input file", nil, nil, 3, 3, ErrInvalidFile},
		{"nil output file", fileOpenHelper(ValidInputFile), nil, 3, 3, nil},
		{"invalid input file", fileOpenHelper(InvalidInputFile), nil, 3, 3, ErrErroneousFile},
		{"valid input file", fileOpenHelper(ValidInputFile), fileOpenHelper(DefaultOutputFile), 3, 3, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.inF != nil {
				defer tc.inF.Close()
			}
			logger := log.New(io.Discard, "", 0)
			validator := validator.NewRegexValidator(logger)
			parser := NewDefaultParser(logger, validator)
			err := parser.ParseNginxLog(tc.inF, tc.outF, tc.maxIps, tc.maxPaths)
			if err != tc.exp {
				t.Fatalf("unexpected error returned: %v != %v", err, tc.exp)
			}
			if tc.outF != nil {
				tc.outF.Close()
				os.Remove(DefaultOutputFile)
			}
		})
	}
}

func fileOpenHelper(file string) *os.File {
	f, _ := os.Open(file)
	return f
}

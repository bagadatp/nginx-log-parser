package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strconv"

	"gitlab.autoiterative.com/group-lucid-swirles-heyrovsky/log-parser/pkg/validator"
)

var (
	ErrInvalidFile   = errors.New("invalid file pointer")
	ErrErroneousFile = errors.New("file with erroneous lines")
)

type LogStats struct {
	TotalLines       uint64            `json:"total_number_of_lines_processed"`
	TotalLinesOk     uint64            `json:"total_number_of_lines_ok"`
	TotalLinesFailed uint64            `json:"total_number_of_lines_failed"`
	Ips              map[string]uint64 `json:"top_client_ips,omitempty"`
	Paths            map[string]string `json:"top_path_avg_seconds,omitempty"`
}

type Parser interface {
	ParseNginxLog(*os.File, *os.File, int, int) error
}

type DefaultParser struct {
	l *log.Logger
	v validator.Validator
}

func NewDefaultParser(logger *log.Logger, val validator.Validator) *DefaultParser {
	return &DefaultParser{
		l: logger,
		v: val,
	}
}

func (p *DefaultParser) ParseNginxLog(inF *os.File, outF *os.File, maxIps, maxPaths int) error {
	if inF == nil {
		return ErrInvalidFile
	}
	scanner := bufio.NewScanner(inF)

	var writer *bufio.Writer
	writer = bufio.NewWriter(io.Discard)
	if outF != nil {
		writer = bufio.NewWriter(outF)
		defer writer.Flush()
	}

	logStats := LogStats{0, 0, 0, nil, nil}
	for {
		if !scanner.Scan() {
			break
		}
		nextLine := scanner.Text()
		logStats.TotalLines++
		tokens := p.v.ValidateLine(nextLine)
		if len(tokens) == 0 {
			logStats.TotalLinesFailed++
			continue
		}
		tokenMap, ok := p.v.ValidateTokens(tokens)
		if !ok {
			logStats.TotalLinesFailed++
			continue
		}
		logStats.TotalLinesOk++

		// get the tokens
		ip := p.v.GetToken(0, tokenMap)
		path := p.v.GetToken(4, tokenMap)
		duration := p.v.GetToken(7, tokenMap)
		durationInt, _ := strconv.ParseUint(duration, 10, 64)
		countIp(ip)
		countPath(path, durationInt)
	}

	logStats.Ips = getTopIps(maxIps)
	logStats.Paths = getTopPaths(maxPaths)

	// encode to JSON string
	byteArr, e := json.Marshal(logStats)
	if e != nil {
		p.l.Printf("error marshalling object to JSON string")
		return e
	}
	// write to output file
	bytesWritten, e := writer.Write(byteArr)
	if bytesWritten != len(byteArr) {
		p.l.Printf("unable to write - written fewer bytes than requested")
		return e
	}

	if logStats.TotalLinesFailed > 0 {
		return ErrErroneousFile
	}
	return scanner.Err()
}

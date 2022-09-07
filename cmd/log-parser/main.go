package main

import (
	"bufio"
	"fmt"
	"gitlab.autoiterative.com/group-lucid-swirles-heyrovsky/log-parser/pkg/validator"
	"log"
	"os"
	"strconv"

	"github.com/pborman/getopt/v2"

	"gitlab.autoiterative.com/group-lucid-swirles-heyrovsky/log-parser/pkg/parser"
)

var (
	stdIOStr           = "-"
	defaultMaxIpsStr   = "10"
	defaultMaxIps, _   = strconv.Atoi(defaultMaxIpsStr)
	defaultMaxPathsStr = defaultMaxIpsStr
	defaultMaxPaths, _ = strconv.Atoi(defaultMaxIpsStr)

	minMaxIps   = 0
	maxMaxIps   = 10000
	minMaxPaths = minMaxIps
	maxMaxPaths = maxMaxIps
)

const (
	version  = "v0.0.1-dev"
	cfgUsage = `
Usage: log-parser -i <INPUT_FILE> -o <OUTPUT_FILE> [-I <MAX_IPS>] [-P <MAX_PATHS>]
   -i, --in=<file>    specifies the input file, required
   -o, --out=<file>   specifies the output JSON file, required.
   -I, --max-client-ips=<IPs>    IPs is the maximum number of results
              to output in the 'top_client_ips' field. Default is 10. Allowed values are 0-10000.
   -P, --max-paths=<PATHs>       PATHs is the maximum number of results
              to output on the 'top_path_avg_seconds' field. Default is 10. Allowed values are 0-10000.
`
)

func dumpInput(logger *log.Logger, prefix, in, out string, maxIps, maxPaths int) {
	logger.Printf("%v: %v, %v, %v, %v", prefix, in, out, maxIps, maxPaths)
	file, e := os.Open(in)
	defer file.Close()
	if e != nil {
		logger.Printf("error opening file '%v'", in)
		return
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		nextLine := fileScanner.Text()
		logger.Printf(":::%v", nextLine)
		lineCount++
	}
	logger.Printf("number of lines: %v", lineCount)
}

func main() {

	l := log.New(os.Stdout, fmt.Sprintf("log-parser/%d ", os.Getpid()), 0)
	//l.Println(os.Args)

	//// Init
	var (
		inFile      = getopt.StringLong("in", 'i', "", "Input nginx log file")
		outFile     = getopt.StringLong("out", 'o', "", "Output JSON file")
		maxIpsStr   = getopt.StringLong("max-client-ips", 'I', defaultMaxIpsStr, "maximum number of results to output in the top_client_ips field")
		maxPathsStr = getopt.StringLong("max-paths", 'P', defaultMaxPathsStr, "maximum number of results to output on the top_path_avg_seconds field")
	)
	getopt.SetUsage(func() {
		fmt.Fprintf(os.Stderr, "%v", cfgUsage)
	})
	getopt.ParseV2()
	maxIps, e := strconv.Atoi(*maxIpsStr)
	if e != nil {
		maxIps = defaultMaxIps
	}
	maxPaths, e := strconv.Atoi(*maxPathsStr)
	if e != nil {
		maxPaths = defaultMaxPaths
	}
	//dumpInput(l, "parsed args", *inFile, *outFile, maxIps, maxPaths)

	if maxIps < minMaxIps || maxIps > maxMaxIps || maxPaths < minMaxPaths || maxPaths > maxMaxPaths {
		getopt.Usage()
		os.Exit(1)
	}
	if inFile == nil || *inFile == "" || outFile == nil || *outFile == "" {
		getopt.Usage()
		os.Exit(1)
	}
	//dumpInput(l, "validated args", *inFile, *outFile, maxIps, maxPaths)
	if _, err := os.Stat(*inFile); *inFile != stdIOStr && err != nil {
		getopt.Usage()
		os.Exit(1)
	}
	var inF, outF *os.File
	if *inFile == stdIOStr {
		inF = os.Stdin
	} else {
		inF, e = os.Open(*inFile)
		if e != nil {
			l.Fatalf("Could not open input file %s", *inFile)
		}
		defer inF.Close()
	}

	if *outFile == stdIOStr {
		outF = os.Stdout
	} else {
		outF, e = os.Create(*outFile)
		if e != nil {
			l.Fatalf("Could not open output file %s", *outFile)
		}
		defer outF.Close()
	}

	//// Create the objects
	rValidator := validator.NewRegexValidator(l)
	dParser := parser.NewDefaultParser(l, rValidator)

	//// Main
	e = dParser.ParseNginxLog(inF, outF, maxIps, maxPaths)
	if e != nil {
		//l.Printf("Error parsing file %v", e)
	}
}

package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"
)

const usageText = `
Usage of %s:
   %s [<inputFileName>]
`
const (
	exitCodeOK = iota
	exitCodeNG
	exitCodeParseFlagErr
	exitCodeFileOpenErr
	exitCodeFlagErr
	exitCodeCsvFormatErr
)

type cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
	delimiter            string
	header               bool
	redmine              bool
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin, delimiter: "\t", header: false, redmine: false}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	flags := flag.NewFlagSet("tsv2mdtbl", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	flags.StringVar(&c.delimiter, "delimiter", "\t", "delimiter character")
	flags.StringVar(&c.delimiter, "d", "\t", "delimiter character")
	flags.BoolVar(&c.header, "header", false, "use header")
	flags.BoolVar(&c.header, "H", false, "use header")
	flags.BoolVar(&c.redmine, "redmine", false, "output redmine table")
	flags.BoolVar(&c.redmine, "r", false, "output redmine table")

	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeParseFlagErr
	}
	param := flags.Args()
	// fmt.Println(param)

	records := validateParam(param, c.inStream, c.delimiter)

	output := csv2MdTbl(records, c.header, c.redmine)
	write(c.outStream, output)

	return exitCodeOK
}

func validateParam(param []string, inStream io.Reader, delimiter string) (records [][]string) {
	var file string
	var reader io.Reader
	var err error
	switch len(param) {
	case 0:
		reader = bufio.NewReader(inStream)
	case 1:
		file = param[0]
		f, err := os.Open(file)
		if err != nil {
			fatal(err, exitCodeFileOpenErr)
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	default:
		fatal(errors.New("failed to read param"), exitCodeFlagErr)
	}

	csvr := csv.NewReader(reader)
	delm, _ := utf8.DecodeLastRuneInString(delimiter)
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

	records, err = csvr.ReadAll()
	if err != nil {
		fatal(err, exitCodeCsvFormatErr)
	}

	return records
}

func csv2MdTbl(records [][]string, header bool, redmine bool) (output []string) {
	for i, r := range records {
		outputLine := "| " + strings.Join(r, " | ") + " |"
		output = append(output, outputLine)
		if i == 0 && header {
			rowLength := len(records[0])
			separator := make([]string, rowLength)
			for j := range separator {
				separator[j] = "---"
			}
			separatorLine := "| " + strings.Join(separator, " | ") + " |"
			output = append(output, separatorLine)
		}
		if i == 0 && header && redmine {
			headerLine := "|_. " + strings.Join(r, " |_. ") + " |"
			output = []string{}
			output = append(output, headerLine)
		}
	}

	return output
}

func write(outStream io.Writer, output []string) {
	for _, o := range output {
		fmt.Fprintln(outStream, o)
	}
}

func fatal(err error, errorCode int) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(errorCode)
}

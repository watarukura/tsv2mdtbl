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
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"
	"strings"
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
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin, delimiter: "\t", header: false}
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

	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeParseFlagErr
	}
	param := flags.Args()
	// fmt.Println(param)

	records := validateParam(param, c.inStream, c.delimiter)

	csv2MdTbl(records, c.header, c.outStream)

	return exitCodeOK
}

func validateParam(param []string, inStream io.Reader, delimiter string) (records [][]string) {
	var file string
	var reader io.Reader
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

	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatal(err, exitCodeCsvFormatErr)
		}

		// セル内改行は<br>タグに置換する
		for i, f := range record {
			if strings.Contains(f, "\n") {
				record[i] = strings.Replace(f, "\n", "<br>", -1)
			}
		}
		records = append(records, record)
	}

	//fmt.Println(records)
	return records
}

func csv2MdTbl(records [][]string, header bool, outStream io.Writer) {
	table := tablewriter.NewWriter(outStream)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	// Header左寄せ
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	if header {
		table.SetHeader(records[0])
		table.AppendBulk(records[1:])
	} else {
		table.AppendBulk(records)
	}
	table.Render()
}

func fatal(err error, errorCode int) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ", os.Args[0], fn, line, err)
	os.Exit(errorCode)
}

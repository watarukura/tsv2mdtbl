package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTsv2MdTblStdInput(t *testing.T) {
	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input      string
		inputStdin string
		want       string
	}{
		{
			input:      "-H",
			inputStdin: "1\t2\n3\t4\n",
			want:       "| 1 | 2 |\n|---|---|\n| 3 | 4 |\n",
		},
		{
			input:      "-header",
			inputStdin: "1\t2\n3\t4\n",
			want:       "| 1 | 2 |\n|---|---|\n| 3 | 4 |\n",
		},
		{
			input:      "--header",
			inputStdin: "1\t2\n3\t4\n",
			want:       "| 1 | 2 |\n|---|---|\n| 3 | 4 |\n",
		},
		{
			input:      "-d ,",
			inputStdin: "1,2\n3,4\n",
			want:       "| 1 | 2 |\n| 3 | 4 |\n",
		},
		{
			input:      "--delimiter=,",
			inputStdin: "1,2\n3,4\n",
			want:       "| 1 | 2 |\n| 3 | 4 |\n",
		},
		// {
		// 	input:      "--header --redmine",
		// 	inputStdin: "1\t2\n3\t4\n",
		// 	want:       "|_. 1 |_. 2 |\n| 3 | 4 |\n",
		// },
		{
			input:      "-H",
			inputStdin: "1\t2\n3\t4567890\n",
			want:       "| 1 | 2       |\n|---|---------|\n| 3 | 4567890 |\n",
		},
		{
			input:      "-H",
			inputStdin: "1\t2\n3\t\"4567\n890\"\n",
			want:       "| 1 | 2           |\n|---|-------------|\n| 3 | 4567<br>890 |\n",
		},
		{
			input:      "-d ,",
			inputStdin: "1,2\n3,\n",
			want:       "| 1 | 2 |\n| 3 |   |\n",
		},
	}
	for i, c := range cases {
		inStream.Reset()
		outStream.Reset()
		errStream.Reset()
		inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream, inStream: inStream}

		args := append([]string{"tsv2mdtbl"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Case: %d, Unexpected output: \n%s, want: \n%s", i, outStream.String(), c.want)
		}
	}
}
func TestTsv2MdTblFileInput(t *testing.T) {
	// outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	outStream, errStream, _ := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	cases := []struct {
		input string
		// inputStdin string
		want string
	}{
		{
			input: "testdata/TEST1.txt",
			// inputStdin: "",
			want: "| 1 | 2 |\n| 3 | 4 |\n",
		},
		{
			input: "-H testdata/TEST1.txt",
			// inputStdin: "",
			want: "| 1 | 2 |\n|---|---|\n| 3 | 4 |\n",
		},
		{
			input: "--delimiter=, --header testdata/TEST2.txt",
			// inputStdin: "",
			want: "| NUM1 | NUM2   |\n|------|--------|\n|    1 |      2 |\n|    3 | 4<br>5 |\n",
		},
	}
	for i, c := range cases {
		// inStream.Reset()
		outStream.Reset()
		errStream.Reset()
		// inStream = bytes.NewBufferString(c.inputStdin)
		cli := &cli{outStream: outStream, errStream: errStream}

		args := append([]string{"tsv2mdtbl"}, strings.Split(c.input, " ")...)
		status := cli.run(args)
		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		if outStream.String() != c.want {
			t.Errorf("Case: %d, Unexpected output: \n%s, want: \n%s", i, outStream.String(), c.want)
		}
	}
}

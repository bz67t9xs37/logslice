// Command logslice extracts a time-range window from a structured log file.
//
// Usage:
//
//	logslice -start <time> -end <time> [-out <file>] <logfile>
package main

import (
	"flag"
	"fmt"
	"os"

	"logslice/internal/splitter"
)

func main() {
	start := flag.String("start", "", "start of time window (required)")
	end := flag.String("end", "", "end of time window (required)")
	out := flag.String("out", "", "output file (default: stdout)")
	flag.Parse()

	if *start == "" || *end == "" {
		fmt.Fprintln(os.Stderr, "error: -start and -end are required")
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "error: exactly one log file argument required")
		flag.Usage()
		os.Exit(1)
	}

	input, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening input: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	var output *os.File
	if *out == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(*out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	res, err := splitter.Run(splitter.Config{
		Input:  input,
		Output: output,
		Start:  *start,
		End:    *end,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "scanned %d lines, wrote %d lines\n",
		res.LinesScanned, res.LinesWritten)
}

// Package splitter provides the high-level pipeline that ties together
// time parsing, line filtering, and output writing.
//
// Typical usage:
//
//	res, err := splitter.Run(splitter.Config{
//		Input:  logFile,
//		Output: os.Stdout,
//		Start:  "2024-01-10T09:00:00Z",
//		End:    "2024-01-10T17:00:00Z",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("wrote %d lines\n", res.LinesWritten)
//
// The pipeline is:
//
//	io.ReadSeeker → scanner → filter → output.Writer
package splitter

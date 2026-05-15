// Package scanner implements sequential log-line scanning with time-window
// filtering.
//
// A Scanner wraps a [timeparse.Parser] and an [io.Reader], emitting only
// lines whose timestamps fall within a caller-supplied [time.Time] range.
// Lines that carry no recognisable timestamp are treated as continuation
// lines and are forwarded to the caller only when the scanner is already
// inside the active window — preserving multi-line log entries such as
// stack traces.
//
// Scanning stops early once a timestamped line is found that lies after
// the end of the requested window, taking advantage of the chronological
// ordering that is typical of log files.
//
// Usage:
//
//	p, _ := timeparse.New(nil)
//	s := scanner.New(p, startTime, endTime)
//	err := s.Scan(file, func(line string) error {
//		_, err := fmt.Fprintln(out, line)
//		return err
//	})
package scanner

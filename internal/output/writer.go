// Package output handles writing filtered log lines to a destination.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer wraps an io.Writer with buffered output and optional line counting.
type Writer struct {
	w          *bufio.Writer
	linesWritten int
	closer     io.Closer
}

// New creates a Writer that writes to the given io.Writer.
// If w also implements io.Closer, it will be closed on Close.
func New(w io.Writer) *Writer {
	var closer io.Closer
	if c, ok := w.(io.Closer); ok {
		closer = c
	}
	return &Writer{
		w:      bufio.NewWriterSize(w, 64*1024),
		closer: closer,
	}
}

// NewFile opens (or creates) a file at path and returns a Writer for it.
// The caller must call Close to flush and close the underlying file.
func NewFile(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: create file %q: %w", path, err)
	}
	return New(f), nil
}

// WriteLine writes a single log line followed by a newline character.
func (w *Writer) WriteLine(line string) error {
	if _, err := fmt.Fprintln(w.w, line); err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	w.linesWritten++
	return nil
}

// LinesWritten returns the total number of lines written so far.
func (w *Writer) LinesWritten() int {
	return w.linesWritten
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return nil
}

// Close flushes buffered data and closes the underlying writer if applicable.
func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	if w.closer != nil {
		if err := w.closer.Close(); err != nil {
			return fmt.Errorf("output: close: %w", err)
		}
	}
	return nil
}

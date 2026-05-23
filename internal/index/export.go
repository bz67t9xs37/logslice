package index

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat controls the serialisation format used by Export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// exportEntry is the JSON/CSV-friendly representation of an Entry.
type exportEntry struct {
	Timestamp string `json:"timestamp"`
	Offset    int64  `json:"offset"`
	Length    int    `json:"length"`
}

func toExportEntries(entries []Entry) []exportEntry {
	out := make([]exportEntry, len(entries))
	for i, e := range entries {
		out[i] = exportEntry{
			Timestamp: e.Timestamp.UTC().Format(time.RFC3339Nano),
			Offset:    e.Offset,
			Length:    e.Length,
		}
	}
	return out
}

// Export writes index entries to w in the requested format.
func Export(w io.Writer, entries []Entry, format ExportFormat) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, entries)
	case FormatCSV:
		return exportCSV(w, entries)
	case FormatText:
		return exportText(w, entries)
	default:
		return fmt.Errorf("index/export: unknown format %q", format)
	}
}

func exportJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(toExportEntries(entries))
}

func exportCSV(w io.Writer, entries []Entry) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "offset", "length"}); err != nil {
		return err
	}
	for _, e := range toExportEntries(entries) {
		row := []string{
			e.Timestamp,
			fmt.Sprintf("%d", e.Offset),
			fmt.Sprintf("%d", e.Length),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportText(w io.Writer, entries []Entry) error {
	for _, e := range toExportEntries(entries) {
		if _, err := fmt.Fprintf(w, "%s\toffset=%d\tlength=%d\n",
			e.Timestamp, e.Offset, e.Length); err != nil {
			return err
		}
	}
	return nil
}

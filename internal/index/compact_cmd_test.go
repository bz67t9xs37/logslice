package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeTempLogForCompact creates a small log file and returns its path.
func writeTempLogForCompact(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "compact.log")

	lines := []string{
		"2024-01-01T00:00:00Z level=info msg=a\n",
		"2024-01-01T00:01:00Z level=info msg=b\n",
		"2024-01-01T00:02:00Z level=info msg=c\n",
		"2024-01-01T00:10:00Z level=info msg=d\n",
		"2024-01-01T00:11:00Z level=info msg=e\n",
	}
	var data []byte
	for _, l := range lines {
		data = append(data, []byte(l)...)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCompactFile_ReducesEntries(t *testing.T) {
	logPath := writeTempLogForCompact(t)
	opts := CompactOptions{BucketSize: 5 * time.Minute}

	out, res, err := CompactFile(logPath, opts)
	if err != nil {
		t.Fatalf("CompactFile: %v", err)
	}
	if res.InputCount == 0 {
		t.Error("expected non-zero input count")
	}
	if len(out) == 0 {
		t.Error("expected at least one compacted entry")
	}
	if len(out) >= res.InputCount {
		t.Errorf("expected fewer entries after compaction: %d >= %d", len(out), res.InputCount)
	}
}

func TestRunCompact_WritesOutput(t *testing.T) {
	logPath := writeTempLogForCompact(t)
	opts := DefaultCompactOptions()

	if err := RunCompact(logPath, opts); err != nil {
		t.Fatalf("RunCompact: %v", err)
	}
}

func TestDefaultCompactOptions(t *testing.T) {
	opts := DefaultCompactOptions()
	if opts.BucketSize != time.Minute {
		t.Errorf("BucketSize: got %v, want 1m", opts.BucketSize)
	}
	if opts.MaxEntries != 0 {
		t.Errorf("MaxEntries: got %d, want 0", opts.MaxEntries)
	}
}

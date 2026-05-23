package index

import (
	"os"
	"testing"
	"time"
)

func writeTempLogForRepair(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "repair-log-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	base := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		fmt.Fprintf(f, "%s log line %d\n", base.Add(time.Duration(i)*time.Second).Format(time.RFC3339), i)
	}
	return f.Name()
}

func TestRepairFile_RoundTrip(t *testing.T) {
	logPath := writeTempLogForRepair(t)
	defer os.Remove(logPath)

	// Build an initial cache.
	_, err := Build(logPath, nil)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	result, err := RepairFile(logPath, DefaultRepairOptions())
	if err != nil {
		t.Fatalf("RepairFile: %v", err)
	}
	if result.OriginalCount == 0 {
		t.Error("expected non-zero OriginalCount")
	}
	if result.RepairedCount > result.OriginalCount {
		t.Error("RepairedCount should not exceed OriginalCount")
	}

	cachePath := CachePath(logPath)
	defer os.Remove(cachePath)
}

func TestDefaultRepairOptions(t *testing.T) {
	opts := DefaultRepairOptions()
	if opts.Verbose {
		t.Error("expected Verbose=false by default")
	}
}

func TestRunRepair_MissingCache(t *testing.T) {
	err := RunRepair("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing log file")
	}
}

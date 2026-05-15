# logslice

Fast log file splitter that extracts time-range windows from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries within a specific time range:

```bash
logslice --input app.log --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"
```

Write the output to a file:

```bash
logslice --input app.log --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --output slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Path to the source log file | required |
| `--from` | Start of the time window (RFC3339) | required |
| `--to` | End of the time window (RFC3339) | required |
| `--output` | Output file path | stdout |
| `--format` | Log timestamp format (`json`, `logfmt`) | `json` |

### Example

```bash
# Extract the last hour of errors from a large JSON log file
logslice --input /var/log/app.log \
         --from "2024-01-15T14:00:00Z" \
         --to "2024-01-15T15:00:00Z" \
         --format json \
         --output errors.log
```

## How It Works

`logslice` uses binary search to locate the start and end positions of the target time window, avoiding the need to scan the entire file. This makes it significantly faster than `grep`-based approaches on large log files.

## Requirements

- Go 1.21 or later

## License

MIT — see [LICENSE](LICENSE) for details.
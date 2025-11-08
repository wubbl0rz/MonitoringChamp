# Directory Size Exporter

This project is a Prometheus exporter that periodically scans a directory and reports the size (in bytes) of each subdirectory. The metrics are exposed over HTTP for Prometheus scraping.

***

### Features

- Exports directory sizes as Prometheus metrics.
- Configurable via command-line arguments or environment variables.
- Supports tab completion for CLI commands.
- Includes verbose logging option for debugging and monitoring.

***

### Usage

Start the exporter with the `start` command:

```bash
./dirsize-exporter start
```

To display help and available options:

```bash
./dirsize-exporter start --help
```

To enable shell autocompletion:

```bash
./dirsize-exporter completion -c zsh
```

***

### Configuration Options

You can configure the exporter using command-line flags or environment variables.

| Option | Description | Environment Variable | Default Value |
|---------|--------------|----------------------|----------------|
| `--verbose` | Enables verbose logging (boolean). | `EXPORTER_VERBOSE` | `0` |
| `--data-dir` | Path to the data directory to monitor. | `EXPORTER_DATA_DIR` | `/data` |
| `--interval`, `-i` | Refresh interval in seconds. | `EXPORTER_INTERVAL` | `5` |
| `--port` | Listening port for the metrics endpoint. | `EXPORTER_PORT` | `9100` |

***

### Metrics Exposed

| Metric Name | Type | Description |
|--------------|------|-------------|
| `dir_size_bytes{dir="..."}` | Gauge | Directory size in bytes for each subdirectory in the data directory. |

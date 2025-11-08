package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alecthomas/kong"
	kongcompletion "github.com/jotaen/kong-completion"
)

var CLI struct {
	Start      StartCommand              `cmd:"" help:"Start exporter."`
	Completion kongcompletion.Completion `cmd:"" help:"Outputs shell code for initialising tab completions"`
}

type StartCommand struct {
	Verbose  bool   `help:"Verbose logging." env:"EXPORTER_VERBOSE" default:"0"`
	DataDir  string `help:"Path to data dir." env:"EXPORTER_DATA_DIR" default:"/data"`
	Interval int8   `help:"Refresh interval in seconds." env:"EXPORTER_INTERVAL" short:"i" default:"5"`
	Port     int16  `help:"Listen port." env:"EXPORTER_PORT" default:"9100"`
}

var dirSizeMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "dir_size_bytes",
	Help: "Directory size in bytes",
}, []string{"dir"})

func DirSize(path string) (int64, error) {
	var size int64

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func Refresh(path string) error {
	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	dirSizeMetric.Reset()

	var results []slog.Attr

	for _, dirEntry := range entries {
		if !dirEntry.IsDir() {
			continue
		}

		path := filepath.Join(path, dirEntry.Name())

		size, err := DirSize(path)

		if err != nil {
			return err
		}

		results = append(results, slog.Attr{
			Key:   path,
			Value: slog.Int64Value(size),
		})

		dirSizeMetric.WithLabelValues(path).Set(float64(size))
	}

	slog.Info("Refreshing", "Sizes", results)

	return nil
}

func (l *StartCommand) Run() error {
	slog.Info("Starting exporter.",
		slog.Group("Config", "Port", l.Port, "DataDir", l.DataDir, "Verbose", l.Verbose, "Interval", l.Interval))

	registry := prometheus.NewRegistry()

	registry.MustRegister(dirSizeMetric)

	go func() {
		for {
			err := Refresh(l.DataDir)

			if err != nil {
				slog.Warn(err.Error())
			}

			time.Sleep(time.Duration(l.Interval) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	return http.ListenAndServe(fmt.Sprintf(":%d", l.Port), nil)
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := kong.Must(&CLI)

	kongcompletion.Register(app)

	ctx, err := app.Parse(os.Args[1:])

	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = ctx.Run()

	if err != nil {
		slog.Error(err.Error())
		return
	}
}

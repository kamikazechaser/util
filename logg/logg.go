// Package logg provides a wrapper around Go's structured logging package slog,
// offering convenient configuration options for different log formats and levels.
//
// The main entry point is NewLogg, which creates a new logger with customizable options.
//
// Example usage:
//
//	logger := logg.NewLogg(logg.Opts{
//		Component: "my-service",
//		LogLevel:  slog.LevelInfo,
//		FormatType: logg.JSON,
//	})
package logg

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type (
	// FormatType represents the different supported log formats.
	FormatType int

	// Opts customize [slog.HandlerOptions] and provides configuration
	// options for creating a new logger instance.
	Opts struct {
		// Writer defaults to os.Stderr if nil
		Writer io.Writer
		// Component enriches each log line with a componenent key/value.
		// Useful for aggregating/filtering with your log collector.
		Component string
		// Group nests individual keys in the format group.child.
		Group string
		// Log format.
		// Logfmt is the default log format
		// Human prints colourized logs useful for CLIs or development
		FormatType FormatType
		// Minimal level to log.
		// Debug level will automatically enable source code location.
		LogLevel slog.Level
	}
)

const (
	// Logfmt format uses Go's default key=value text format for logging.
	Logfmt FormatType = iota
	// Human format provides colorized, human-readable logs suitable for CLI applications
	// and development environments.
	Human
	// JSON format outputs logs in structured JSON format for machine processing and
	// integration with log aggregation systems.
	JSON
)

// NewLogg creates a new logger with the given options.
// Default log format is Logfmt.
func NewLogg(o Opts) *slog.Logger {
	w := o.Writer
	if w == nil {
		w = os.Stderr
	}

	var handler slog.Handler

	addSource := o.LogLevel == slog.LevelDebug

	switch o.FormatType {
	case Human:
		handler = tint.NewHandler(w, &tint.Options{
			Level:      o.LogLevel,
			TimeFormat: time.Kitchen,
			AddSource:  addSource,
		})
	case JSON:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:     o.LogLevel,
			AddSource: addSource,
		})
	default:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     o.LogLevel,
			AddSource: addSource,
		})
	}

	logger := slog.New(handler)

	if o.Component != "" {
		logger = logger.With("component", o.Component)
	}
	if o.Group != "" {
		logger = logger.WithGroup(o.Group)
	}

	return logger
}

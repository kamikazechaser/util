// Package logg provides a wrapper around Go's structured logging package slog,
// offering convenient configuration options and preset handlers for different log formats and levels.
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
		// Component enriches each log line with a componenent value.
		// Useful for aggregating/filtering with your log collector.
		Component string
		// ExtraAttributes adds extra attributes to each log line.
		// Useful for adding information like region, node, environment, etc.
		ExtraAttributes []slog.Attr
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

	addSource := o.LogLevel <= slog.LevelDebug

	var handler slog.Handler
	switch o.FormatType {
	case Human:
		handler = tint.NewHandler(w, &tint.Options{
			Level:      o.LogLevel,
			TimeFormat: time.Kitchen,
			AddSource:  addSource,
		}).WithAttrs(o.ExtraAttributes)
	case JSON:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:     o.LogLevel,
			AddSource: addSource,
		}).WithAttrs(o.ExtraAttributes)
	case Logfmt:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     o.LogLevel,
			AddSource: addSource,
		}).WithAttrs(o.ExtraAttributes)
	}

	logger := slog.New(handler)

	if o.Component != "" {
		logger = logger.With("component", o.Component)
	}

	return logger
}

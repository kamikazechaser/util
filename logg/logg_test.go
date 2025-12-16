package logg

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewLogg_Logfmt(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogg(Opts{
		Writer:    &buf,
		Component: "test-service",
		ExtraAttributes: []slog.Attr{
			slog.String("region", "us-west"),
			slog.String("env", "production"),
		},
		FormatType: Logfmt,
		LogLevel:   slog.LevelDebug,
	})

	logger.Debug("debug message", "debug_key", "debug_value")
	logger.Info("info message", "info_key", "info_value")
	logger.Warn("warn message", "warn_key", "warn_value")
	logger.Error("error message", "error_key", 123, "bool_key", true)

	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Error("expected debug message")
	}
	if !strings.Contains(output, "info message") {
		t.Error("expected info message")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("expected warn message")
	}
	if !strings.Contains(output, "error message") {
		t.Error("expected error message")
	}

	if !strings.Contains(output, "component=test-service") {
		t.Error("expected component=test-service")
	}

	if !strings.Contains(output, "region=us-west") {
		t.Error("expected region=us-west")
	}
	if !strings.Contains(output, "env=production") {
		t.Error("expected env=production")
	}

	if !strings.Contains(output, "debug_key=debug_value") {
		t.Error("expected debug_key=debug_value")
	}
	if !strings.Contains(output, "info_key=info_value") {
		t.Error("expected info_key=info_value")
	}

	if !strings.Contains(output, "error_key=123") {
		t.Error("expected error_key=123")
	}
	if !strings.Contains(output, "bool_key=true") {
		t.Error("expected bool_key=true")
	}

	if !strings.Contains(output, "source=") {
		t.Error("expected source= in debug logs")
	}

	if strings.HasPrefix(output, "{") {
		t.Error("expected Logfmt format, not JSON")
	}
}

func TestNewLogg_NoExtraAttributes(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogg(Opts{
		Writer:     &buf,
		Component:  "test-service",
		FormatType: JSON,
		LogLevel:   slog.LevelInfo,
	})

	logger.Info("test message", "key", "value")

	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Error("expected test message")
	}
	if !strings.Contains(output, "component") {
		t.Error("expected component field")
	}
	if !strings.Contains(output, "test-service") {
		t.Error("expected test-service value")
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Error("expected key:value in JSON")
	}
}

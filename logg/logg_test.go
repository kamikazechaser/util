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
		Writer:     &buf,
		Component:  "test-service",
		Group:      "requests",
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

	if !strings.Contains(output, "requests.debug_key=debug_value") {
		t.Error("expected requests.debug_key=debug_value")
	}
	if !strings.Contains(output, "requests.info_key=info_value") {
		t.Error("expected requests.info_key=info_value")
	}

	if !strings.Contains(output, "requests.error_key=123") {
		t.Error("expected requests.error_key=123")
	}
	if !strings.Contains(output, "requests.bool_key=true") {
		t.Error("expected requests.bool_key=true")
	}

	if !strings.Contains(output, "source=") {
		t.Error("expected source= in debug logs")
	}

	if strings.HasPrefix(output, "{") {
		t.Error("expected Logfmt format, not JSON")
	}
}

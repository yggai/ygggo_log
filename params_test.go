package ygggo_log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerParams_TextFormatter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf) // 默认 TextFormatter

	logger.Info("param test", map[string]any{"a": 1, "b": 1.5, "c": true}, "d=xxx", 42)

	out := buf.String()
	if !strings.Contains(out, "param test") {
		t.Fatalf("expected base message, got: %s", out)
	}
	// Plain formatter should not contain ANSI codes
	if strings.Contains(out, "\u001b[") || strings.Contains(out, "\x1b[") {
		t.Fatalf("expected no ANSI color codes in text formatter output: %s", out)
	}
	for _, expect := range []string{"a=1", "b=1.5", "c=true", "d=xxx", "42"} {
		if !strings.Contains(out, expect) {
			t.Errorf("missing %q in output: %s", expect, out)
		}
	}
}

func TestLoggerParams_ColorFormatter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	logger.formatter = NewColorFormatter()

	logger.Info("color param test", map[string]any{"a": 1, "b": 1.5, "c": true}, "d=xxx", 42)

	out := buf.String()
	if !strings.Contains(out, "color param test") {
		t.Fatalf("expected base message, got: %s", out)
	}
	// Should contain some ANSI codes for colored keys/values
	codes := []string{ColorCyan, ColorGreen, ColorPurple, ColorYellow, ColorWhite}
	foundCodes := 0
	for _, code := range codes {
		if strings.Contains(out, code) {
			foundCodes++
		}
	}
	if foundCodes < 2 { // at least some of them should appear
		t.Errorf("expected ANSI color codes in output, got: %s", out)
	}
}


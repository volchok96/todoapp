package main

import (
	"testing"
)

func Test_getEnv(t *testing.T) {
	t.Setenv("TEST_ENV", "value")
	tests := []struct {
		key      string
		fallback string
		expected string
	}{
		{"TEST_ENV", "fallback", "value"},
		{"MISSING_ENV", "fallback", "fallback"},
	}

	for _, tt := range tests {
		got := getEnv(tt.key, tt.fallback)
		if got != tt.expected {
			t.Errorf("getEnv(%q, %q) = %q; want %q", tt.key, tt.fallback, got, tt.expected)
		}
	}
}

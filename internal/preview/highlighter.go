package preview

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"gopkg.in/yaml.v3"
)

// Highlighter provides syntax highlighting for file content
type Highlighter struct {
	style     string // Chroma style: "monokai", "dracula", "github", etc.
	formatter string // Chroma formatter: "terminal256", "terminal16", "terminal16m"
}

// NewHighlighter creates a new highlighter with default settings
func NewHighlighter() *Highlighter {
	return &Highlighter{
		style:     "github",
		formatter: "terminal256",
	}
}

// NewHighlighterWithStyle creates a highlighter with custom style
func NewHighlighterWithStyle(style, formatter string) *Highlighter {
	return &Highlighter{
		style:     style,
		formatter: formatter,
	}
}

// Highlight detects content type and applies syntax highlighting
// Returns highlighted string or plain text on error
func (h *Highlighter) Highlight(content []byte, filename string) string {
	if len(content) == 0 {
		return ""
	}

	// 1. Try extension-based detection first (fast path)
	if filename != "" {
		if lexer := lexers.Match(filename); lexer != nil {
			if result, err := h.highlightWith(content, lexer.Config().Name); err == nil {
				return result
			}
		}
	}

	// 2. Try content-based detection with heuristics
	if lang := detectLanguage(content); lang != "" {
		if result, err := h.highlightWith(content, lang); err == nil {
			return result
		}
	}

	// 3. Try Chroma's built-in analysis (limited support)
	if lexer := lexers.Analyse(string(content)); lexer != nil {
		if result, err := h.highlightWith(content, lexer.Config().Name); err == nil {
			return result
		}
	}

	// 4. Fallback to plain text (no highlighting)
	return string(content)
}

// HighlightWithLanguage applies highlighting for a specific language
func (h *Highlighter) HighlightWithLanguage(content []byte, language string) string {
	if result, err := h.highlightWith(content, language); err == nil {
		return result
	}
	return string(content)
}

// SetStyle updates the color scheme
func (h *Highlighter) SetStyle(style string) {
	h.style = style
}

// SetFormatter updates the terminal formatter
func (h *Highlighter) SetFormatter(formatter string) {
	h.formatter = formatter
}

// highlightWith applies highlighting with the specified language
func (h *Highlighter) highlightWith(content []byte, language string) (string, error) {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, string(content), language, h.formatter, h.style)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// detectLanguage uses simple heuristics for common formats
// This handles extension-less files common in key-value stores like Consul
func detectLanguage(content []byte) string {
	if len(content) == 0 {
		return ""
	}

	// Trim whitespace for detection
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return ""
	}

	// JSON detection (fast path - very common in Consul/etcd)
	if isJSON(trimmed) {
		return "json"
	}

	// YAML detection
	if isYAML(trimmed) {
		return "yaml"
	}

	// XML detection
	if bytes.HasPrefix(trimmed, []byte("<?xml")) ||
		(bytes.HasPrefix(trimmed, []byte("<")) && bytes.Contains(trimmed, []byte(">"))) {
		return "xml"
	}

	// TOML detection
	if isLikelyTOML(trimmed) {
		return "toml"
	}

	// Shell script detection
	if bytes.HasPrefix(trimmed, []byte("#!/bin/bash")) ||
		bytes.HasPrefix(trimmed, []byte("#!/bin/sh")) ||
		bytes.HasPrefix(trimmed, []byte("#!/usr/bin/env bash")) {
		return "bash"
	}

	// Python detection
	if bytes.HasPrefix(trimmed, []byte("#!/usr/bin/env python")) ||
		bytes.HasPrefix(trimmed, []byte("#!/usr/bin/python")) {
		return "python"
	}

	// Markdown detection
	if isLikelyMarkdown(trimmed) {
		return "markdown"
	}

	// INI/config file detection
	if isLikelyINI(trimmed) {
		return "ini"
	}

	return ""
}

// isJSON checks if content is valid JSON
func isJSON(content []byte) bool {
	if len(content) == 0 {
		return false
	}

	first := content[0]
	if first != '{' && first != '[' {
		return false
	}

	// Try to parse as JSON (fast validation)
	var out any
	return json.Unmarshal(content, &out) == nil
}

// isYAML checks if content is likely YAML
func isYAML(content []byte) bool {
	if len(content) == 0 {
		return false
	}

	str := string(content)

	// Check for YAML document markers
	if strings.HasPrefix(str, "---\n") || strings.HasPrefix(str, "---\r\n") {
		return true
	}

	// Don't try expensive YAML parsing if it's clearly JSON
	if isJSON(content) {
		return false
	}

	// Check for common YAML patterns
	lines := strings.Split(str, "\n")
	yamlPatterns := 0
	for i, line := range lines {
		if i > 20 { // Only check first 20 lines for performance
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Look for key: value patterns (YAML style)
		if strings.Contains(line, ": ") || strings.HasPrefix(line, "- ") {
			yamlPatterns++
		}
	}

	// If we found multiple YAML patterns, validate with parser
	if yamlPatterns >= 2 {
		var out any
		return yaml.Unmarshal(content, &out) == nil
	}

	return false
}

// isLikelyTOML checks if content is likely TOML
func isLikelyTOML(content []byte) bool {
	str := string(content)
	// TOML typically has [section] headers or key = value pairs
	return (strings.Contains(str, " = ") || strings.Contains(str, "=")) &&
		(strings.Contains(str, "[") && strings.Contains(str, "]"))
}

// isLikelyMarkdown checks if content is likely Markdown
func isLikelyMarkdown(content []byte) bool {
	str := string(content)
	// Check for common markdown patterns
	return strings.HasPrefix(str, "# ") ||
		strings.HasPrefix(str, "## ") ||
		strings.Contains(str, "```") ||
		strings.Contains(str, "](") || // Links
		strings.Contains(str, "**") // Bold
}

// isLikelyINI checks if content is likely INI/config format
func isLikelyINI(content []byte) bool {
	str := string(content)
	lines := strings.Split(str, "\n")

	hasSection := false
	hasKeyValue := false

	for i, line := range lines {
		if i > 20 { // Only check first 20 lines
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for [section]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			hasSection = true
		}

		// Check for key=value
		if strings.Contains(line, "=") && !strings.Contains(line, "==") {
			hasKeyValue = true
		}
	}

	return hasSection && hasKeyValue
}

// DetectLanguage exposes the language detection for external use
func DetectLanguage(content []byte, filename string) string {
	// Try filename first
	if filename != "" {
		ext := filepath.Ext(filename)
		if ext != "" {
			if lexer := lexers.Match(filename); lexer != nil {
				return lexer.Config().Name
			}
		}
	}

	// Try content detection
	return detectLanguage(content)
}

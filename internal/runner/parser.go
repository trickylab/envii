package runner

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/trickylab/envii/internal/model"
)

// ParseDotenv parses .env content into vault variables.
func ParseDotenv(r io.Reader) ([]*model.Var, error) {
	s := bufio.NewScanner(r)
	var vars []*model.Var
	lineNo := 0

	for s.Scan() {
		lineNo++
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, raw, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("line %d: expected KEY=value", lineNo)
		}
		key = strings.TrimSpace(key)
		if !validKey(key) {
			return nil, fmt.Errorf("line %d: invalid key %q", lineNo, key)
		}

		value, err := parseValue(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		vars = append(vars, &model.Var{Key: key, Value: value, Secret: model.IsSecret(key)})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return vars, nil
}

func validKey(key string) bool {
	if key == "" {
		return false
	}
	for i, r := range key {
		if i == 0 && unicode.IsDigit(r) {
			return false
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func parseValue(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	switch raw[0] {
	case '"':
		return parseDoubleQuoted(raw)
	case '\'':
		return parseSingleQuoted(raw)
	default:
		return strings.TrimSpace(stripInlineComment(raw)), nil
	}
}

func parseDoubleQuoted(raw string) (string, error) {
	var b strings.Builder
	escaped := false

	for i := 1; i < len(raw); i++ {
		c := raw[i]
		if escaped {
			switch c {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '"', '\\':
				b.WriteByte(c)
			default:
				b.WriteByte(c)
			}
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '"' {
			if rest := strings.TrimSpace(raw[i+1:]); rest != "" && !strings.HasPrefix(rest, "#") {
				return "", fmt.Errorf("unexpected content after quoted value")
			}
			return b.String(), nil
		}
		b.WriteByte(c)
	}
	return "", fmt.Errorf("unterminated quoted value")
}

func parseSingleQuoted(raw string) (string, error) {
	end := strings.IndexByte(raw[1:], '\'')
	if end < 0 {
		return "", fmt.Errorf("unterminated quoted value")
	}
	end++
	if rest := strings.TrimSpace(raw[end+1:]); rest != "" && !strings.HasPrefix(rest, "#") {
		return "", fmt.Errorf("unexpected content after quoted value")
	}
	return raw[1:end], nil
}

func stripInlineComment(raw string) string {
	for i, r := range raw {
		if r == '#' && (i == 0 || unicode.IsSpace(rune(raw[i-1]))) {
			return raw[:i]
		}
	}
	return raw
}

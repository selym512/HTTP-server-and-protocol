package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const newLine = "\r\n"

func NewHeaders() Headers { return make(Headers) }

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	newLineIndex := bytes.Index(data, []byte(newLine))
	if newLineIndex == -1 {
		return 0, false, nil
	}
	if newLineIndex == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:newLineIndex], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("Invalid header token found: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("Key included invalid character")
	}
	h.Set(key, string(value))
	return newLineIndex + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}
func (h Headers) Overwrite(key, value string) {
	h[key] = value
}
func (h Headers) Remove(key string) {
	delete(h, key)
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}
	return slices.Contains(tokenChars, c)
}

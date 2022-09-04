package gonetable

import (
	"errors"
	"strings"
)

var (
	ErrKeyDelimiter  = errors.New("key delimiter used in key segment")
	ErrKeyNoSegments = errors.New("no key segments provided")
	KeyDelimiter     = "#"
)

type CompositeKey struct {
	HashSegments  []string
	RangeSegments []string
}

func JoinKeySegments(segments []string) (string, error) {
	if len(segments) == 0 {
		return "", ErrKeyNoSegments
	}
	for _, s := range segments {
		if strings.Contains(s, KeyDelimiter) {
			return "", ErrKeyDelimiter
		}
	}
	return strings.Join(segments, KeyDelimiter), nil
}

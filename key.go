package gonetable

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (k CompositeKey) Marshal() (map[string]types.AttributeValue, error) {
	spk, err := joinKeySegments(k.HashSegments)
	if err != nil {
		return nil, err
	}
	ssk, err := joinKeySegments(k.RangeSegments)
	if err != nil {
		return nil, err
	}
	pk, err := attributevalue.Marshal(spk)
	if err != nil {
		return nil, err
	}
	sk, err := attributevalue.Marshal(ssk)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{
		"PK": pk,
		"SK": sk,
	}, nil
}

func joinKeySegments(segments []string) (string, error) {
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

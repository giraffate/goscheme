package main

import (
	"bytes"
	"testing"
)

func TestPeek(t *testing.T) {
	cases := []struct {
		input    []byte
		expected string
	}{
		{input: []byte("("), expected: "("},
		{input: []byte(" ("), expected: "("},
		{input: []byte("\t("), expected: "("},
		{input: []byte("\r("), expected: "("},
		{input: []byte("\n("), expected: "("},
	}

	for _, tc := range cases {
		var got []byte
		l := NewLexer(bytes.NewBuffer(tc.input))
		peek := l.Peek()
		got = append(got, byte(peek))
		if string(got) != tc.expected {
			t.Errorf("Got %s, expected %s\n", got, tc.expected)
		}
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		input    []byte
		expected string
	}{
		{input: []byte("("), expected: "("},
		{input: []byte(" ("), expected: "("},
		{input: []byte("\t("), expected: "("},
		{input: []byte("\r("), expected: "("},
		{input: []byte("\n("), expected: "("},
	}

	for _, tc := range cases {
		var got []byte
		l := NewLexer(bytes.NewBuffer(tc.input))
		next := l.Next()
		got = append(got, byte(next))
		if string(got) != tc.expected {
			t.Errorf("Got %s, expected %s\n", got, tc.expected)
		}
	}
}

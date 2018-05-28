package main

import (
	"bytes"
	"text/scanner"
)

type Lexer struct {
	s     *scanner.Scanner
	token string
}

func NewLexer(b *bytes.Buffer) *Lexer {
	var s scanner.Scanner
	s.Init(b)
	s.Mode = scanner.ScanChars
	return &Lexer{s: &s}
}

func (l *Lexer) Peek() rune {
	p := l.s.Peek()
	switch p {
	case ' ', '\t', '\r', '\n':
		l.s.Next()
		return l.Peek()
	default:
		return p
	}
}

func (l *Lexer) Next() rune {
	n := l.s.Next()
	switch n {
	case ' ', '\t', '\r', '\n':
		return l.Next()
	default:
		return n
	}
}

func (l *Lexer) TextToken() string { return l.token }

func (l *Lexer) Scan() {
	var bs []byte
	b := l.s.Next()
	for {
		switch b {
		case ' ', '\t', '\r', '\n', scanner.EOF:
			b = l.s.Next()
		case '(', ')':
			l.token = string(b)
			return
		default:
			bs = append(bs, byte(b))
			for {
				p := l.s.Peek()
				if isLetter(p) {
					bs = append(bs, byte(l.s.Next()))
				} else {
					break
				}
			}
			l.token = string(bs)
			return
		}
	}
}

func isLetter(b rune) bool {
	switch b {
	case ' ', '\t', '\r', '\n', '(', ')', scanner.EOF:
		return false
	default:
		return true
	}
}

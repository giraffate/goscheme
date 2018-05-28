package main

import (
	"errors"
	"strconv"
	"text/scanner"
)

type Parser struct {
	l *Lexer
}

func NewParser(l *Lexer) *Parser {
	return &Parser{l: l}
}

func (p *Parser) Parse() (Exp, error) {
	for {
		p.l.Scan()
		s := p.l.TextToken()
		switch s {
		case "(":
			return p.parseExps()
		default:
			if n, err := strconv.Atoi(s); err == nil {
				return Number(n), nil
			}
			if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
				return s, nil
			}
			if s == "#t" || s == "#f" {
				return Boolean(s == "#t"), nil
			}
			return Symbol(s), nil
		}
	}
}

func (p *Parser) parseExps() ([]Exp, error) {
	var list []Exp
	for {
		s := p.l.Peek()

		if s == ')' {
			p.l.Next()
			break
		} else if s == scanner.EOF {
			return nil, errors.New("Syntax error, no ')'")
		}

		exp, err := p.Parse()
		if err != nil {
			return nil, err
		}
		list = append(list, exp)
	}
	return list, nil
}

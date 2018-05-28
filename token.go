package main

import (
	"fmt"
)

type Exp interface{}
type Number int
type Symbol string
type Boolean bool

func (b Boolean) String() string {
	if b {
		return "#t"
	}
	return "#f"
}

type Pair struct {
	Car Exp
	Cdr Exp
}

func (p *Pair) isNull() bool {
	if p.Car == nil && p.Cdr == nil {
		return true
	}
	return false
}

func (p *Pair) isList() bool {
	if p.Cdr == nil {
		return true
	}

	switch t := p.Cdr.(type) {
	case *Pair:
		return t.isList()
	default:
		return false
	}
}

func (p *Pair) String() string {
	return fmt.Sprint("<", p.Car, ", ", p.Cdr, ">")
}

type Lambda struct {
	Args Exp
	Body Exp
	Env  *Env
}

func (l Lambda) String() string {
	return ""
}

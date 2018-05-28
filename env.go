package main

type Env struct {
	binds  map[Symbol]Exp
	parent *Env
}

func SetupEnv() *Env {
	binds := make(map[Symbol]Exp)
	binds[Symbol("null")] = &Pair{}

	return &Env{binds: binds}
}

func ExpandEnv(e *Env) *Env {
	return &Env{
		binds:  make(map[Symbol]Exp),
		parent: e,
	}
}

func (e *Env) Get(sym Symbol) Exp {
	v, ok := e.binds[sym]
	if !ok && e.parent != nil {
		v = e.parent.Get(sym)
	}
	return v
}

func (e *Env) Put(k Symbol, v Exp) {
	e.binds[k] = v
}

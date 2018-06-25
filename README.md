# goscheme
goscheme is the tiny scheme interpreter (REPL) written in Go. It's based on [SICP](https://mitpress.mit.edu/sites/default/files/sicp/index.html) chapter 4.

## Example
```
$ ./goscheme
> (+ 1 1)
2
> (define x 3)
> (* x 2)
6
> (define (doubled x) (* x x))
> (doubled 5)
25
```

## Feature
- `+`, `-`, `*`, `/`
- `cons`, `car`, `cdr`, `list`
- `null?`, `eq?`
- `define`, `if`, `cond`, `lambda`, `begin`

## References
These references are so good to make a scheme interpreter in Go.
- [suzuken/gigue](https://github.com/suzuken/gigue)
- [k0kubun/gosick](https://github.com/k0kubun/gosick)
- [chrisbutcher/goscheme](https://github.com/chrisbutcher/goscheme)
- [chenyukang/GoScheme](https://github.com/chenyukang/GoScheme)

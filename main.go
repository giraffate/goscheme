package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	env := SetupEnv()
	var buf []byte

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			if scanner.Err() != nil {
				fmt.Println("input error: ", scanner.Err())
			}
		}
		buf = append(buf, scanner.Bytes()...)
		input := strings.TrimSpace(string(buf))
		output, err, ok := Run(bytes.NewBufferString(input), env)
		if ok {
			if err != nil {
				fmt.Println(err)

			} else if output != nil {
				fmt.Println(output)
			}
			buf = make([]byte, 0, 0)
		}
	}
}

func Run(b *bytes.Buffer, env *Env) (Exp, error, bool) {
	l := NewLexer(b)
	p := NewParser(l)
	exps, err := p.Parse()
	if err != nil {
		return nil, err, false
	}

	ret, err := Eval(exps, env)
	return ret, err, true
}

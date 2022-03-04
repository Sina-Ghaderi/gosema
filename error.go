package main

import (
	"fmt"
	"os"
)

const (
	format = "\033[1;31mfatal --\033[0m %s\n"
	infomt = "\033[1;32mcmlog --\033[0m"
	warnmt = "\033[1;33mwarng --\033[0m"
)

type mainerr error

func defexit() {
	r := recover()
	switch x := r.(type) {
	case mainerr:
		fmt.Printf(format, x.Error())
		os.Exit(-1)
	case nil:
		os.Exit(0)
	default:
		panic(x)
	}
}

func fatal(err error) { panic(mainerr(err)) }
func exits()          { fatal(nil) }
func logst(d ...interface{}) {
	d = append([]interface{}{infomt}, d...)
	fmt.Println(d...)
}

func warng(d ...interface{}) {
	d = append([]interface{}{warnmt}, d...)
	fmt.Println(d...)
}

func usage(st string) {
	fmt.Printf("usage: %v [client|server] .. \n", st)
	exits()
}

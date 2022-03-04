package main

import "errors"

var (
	errNotSame = errors.New("shared memory size is not same as our request")
)

const (
	client = "client"
	server = "server"
)

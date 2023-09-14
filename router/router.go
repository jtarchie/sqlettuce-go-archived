package router

import (
	"fmt"
)

var (
	ErrNoCommandFound  = fmt.Errorf("could not determine command, none were sent")
	ErrIncorrectTokens = fmt.Errorf("received incorrect tokens")
)

const (
	OKResponse   = "+OK\r\n"
	NullResponse = "$-1\r\n"
)

type Router interface {
	Lookup(tokens []string) (Callback, error)
}

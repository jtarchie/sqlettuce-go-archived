package router

import "io"

type Callback func([]string, io.Writer) error

type CallbackRouter Callback

func (c CallbackRouter) Lookup(_ []string) (Callback, bool) {
	return Callback(c), true
}

var _ Router = CallbackRouter(func(s []string, w io.Writer) error {
	return nil
})

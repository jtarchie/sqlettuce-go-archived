package router

import (
	"fmt"
	"io"
)

func staticResponseCallback(response string) Callback {
	return func(_ []string, w io.Writer) error {
		_, err := io.WriteString(w, response)
		if err != nil {
			return fmt.Errorf("could not write static response %q: %w", response, err)
		}

		return nil
	}
}

type StaticResponseRouter string

func (s StaticResponseRouter) Lookup(_ []string) (Callback, bool) {
	return staticResponseCallback(string(s)), true
}

var _ Router = StaticResponseRouter("")

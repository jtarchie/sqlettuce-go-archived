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

type staticResponseRouter string

func (s staticResponseRouter) Lookup(_ []string) (Callback, error) {
	return staticResponseCallback(string(s)), nil
}

var _ Router = staticResponseRouter("")

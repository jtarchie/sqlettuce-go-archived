package handler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jtarchie/sqlettus/db"
)

var (
	ErrNoCommandFound  = fmt.Errorf("could not determine command, none were sent")
	ErrIncorrectTokens = fmt.Errorf("received incorrect tokens")
)

type Callback func([]string, io.Writer) error

type Router interface {
	Lookup(tokens []string) (Callback, error)
}

type Command map[string]Router

func (c Command) Lookup(tokens []string) (Callback, error) {
	next, ok := c[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("could not find command %q: %w", tokens[0], ErrNoCommandFound)
	}

	callback, err := next.Lookup(tokens[1:])
	if err != nil {
		return nil, fmt.Errorf("could not execute command %q: %w", tokens[0], err)
	}

	return callback, nil
}

var _ Router = Command{}

type staticResponse string

func (s staticResponse) Lookup(_ []string) (Callback, error) {
	return func(_ []string, w io.Writer) error {
		_, err := io.WriteString(w, string(s))
		if err != nil {
			return fmt.Errorf("could not write static response %q: %w", s, err)
		}

		return nil
	}, nil
}

var _ Router = staticResponse("")

type tokensLimits struct {
	min, max int
	callback Callback
}

// Lookup implements Router.
func (t *tokensLimits) Lookup(tokens []string) (Callback, error) {
	if 0 < t.min && len(tokens) < t.min {
		return nil, fmt.Errorf("expected minimum number of tokens %d received %d: %w", t.min, len(tokens), ErrIncorrectTokens)
	}

	if 0 < t.max && len(tokens) > t.max {
		return nil, fmt.Errorf("expected maximum number of tokens %d received %d: %w", t.max, len(tokens), ErrIncorrectTokens)
	}

	return t.callback, nil
}

var _ Router = &tokensLimits{}

func SetupRouter(client *db.Client) Command {
	return Command{
		"COMMAND": Command{
			"DOCS": staticResponse("+\r\n"),
		},
		"CONFIG": Command{
			"GET": Command{
				"save":       staticResponse("+\r\n"),
				"appendonly": staticResponse("+no\r\n"),
			},
		},
		"SET": minMaxTokens(2, 0, func(tokens []string, conn io.Writer) error {
			err := client.Set(tokens[1], tokens[2])
			if err != nil {
				return fmt.Errorf("could not execute SET: %w", err)
			}

			_, err = io.WriteString(conn, "+OK\r\n")
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}),
		"GET": minMaxTokens(1, 0, func(tokens []string, conn io.Writer) error {
			value, err := client.Get(tokens[1])
			if err != nil {
				return fmt.Errorf("could not execute GET: %w", err)
			}

			if value == nil {
				_, err = io.WriteString(conn, "+\r\n")
				if err != nil {
					return fmt.Errorf("could not send reply: %w", err)
				}
			} else {
				_, _ = conn.Write([]byte("$"))
				_, _ = io.WriteString(conn, strconv.Itoa(len(*value)))
				_, _ = io.WriteString(conn, "\r\n")
				_, _ = io.WriteString(conn, *value)
				_, err = io.WriteString(conn, "\r\n")
				if err != nil {
					return fmt.Errorf("could not send reply: %w", err)
				}
			}

			return nil
		}),
	}
}

func minMaxTokens(min, max int, callback Callback) *tokensLimits {
	return &tokensLimits{
		min:      min,
		max:      max,
		callback: callback,
	}
}

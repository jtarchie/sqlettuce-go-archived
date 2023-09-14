package handler

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

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
	command := strings.ToUpper(tokens[0])

	next, ok := c[command]
	if !ok {
		slog.Debug("could not find command",
			slog.String("command", command),
		)

		return staticResponseRouter(fmt.Sprintf("-Unsupported command %q\r\n", command)).Lookup(tokens)
	}

	callback, err := next.Lookup(tokens[1:])
	if err != nil {
		slog.Debug("could not lookup command",
			slog.String("command", command),
			slog.String("error", err.Error()),
		)

		return staticResponseRouter("-Unprocessable command\r\n").Lookup(tokens)
	}

	return callback, nil
}

var _ Router = Command{}

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

type tokensLimits struct {
	min, max int
	callback Callback
}

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

type CallbackRouter Callback

func (c CallbackRouter) Lookup(_ []string) (Callback, error) {
	return Callback(c), nil
}

var _ Router = CallbackRouter(func(s []string, w io.Writer) error {
	return nil
})

const OKResponse = "+OK\r\n"

//nolint:funlen
func SetupRouter(client *db.Client) Command {
	return Command{
		"COMMAND": Command{
			"DOCS": staticResponseRouter("+\r\n"),
		},
		"CONFIG": Command{
			"GET": Command{
				"save":       staticResponseRouter("+\r\n"),
				"appendonly": staticResponseRouter("+no\r\n"),
			},
		},
		"FLUSHALL": CallbackRouter(func(_ []string, conn io.Writer) error {
			err := client.FlushAll()
			if err != nil {
				slog.Error("could not FLUSHALL", slog.String("error", err.Error()))
			}

			_, err = io.WriteString(conn, OKResponse)
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}),
		"PING": staticResponseRouter("+PONG\r\n"),
		"SET": minMaxTokens(2, 0, func(tokens []string, conn io.Writer) error {
			err := client.Set(tokens[0], tokens[1])
			if err != nil {
				return fmt.Errorf("could not execute SET: %w", err)
			}

			_, err = io.WriteString(conn, OKResponse)
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}),
		"GET": minMaxTokens(1, 0, func(tokens []string, conn io.Writer) error {
			value, err := client.Get(tokens[0])
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

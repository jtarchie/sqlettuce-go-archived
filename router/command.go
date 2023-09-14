package router

import (
	"fmt"
	"log/slog"
	"strings"
)

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

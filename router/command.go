package router

import (
	"fmt"
	"log/slog"
	"strings"
)

type Command map[string]Router

func (c Command) Lookup(tokens []string) (Callback, bool) {
	command := strings.ToUpper(tokens[0])

	next, ok := c[command]
	if !ok {
		slog.Debug("could not find command",
			slog.String("command", command),
		)

		return staticResponseCallback(fmt.Sprintf("-Unsupported command %q\r\n", command)), false
	}

	return next.Lookup(tokens[1:])
}

var _ Router = Command{}

package router

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/db"
)

//nolint:funlen
func New(
	ctx context.Context,
	client *db.Client,
) Command {
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
		"PING": staticResponseRouter("+PONG\r\n"),
		"ECHO": minMaxTokens(1, 0, func(tokens []string, conn io.Writer) error {
			err := writeBulkString(conn, tokens[1])
			if err != nil {
				return fmt.Errorf("could not echo message: %w", err)
			}

			return nil
		}),
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
		"SET": minMaxTokens(2, 0, func(tokens []string, conn io.Writer) error {
			err := client.Set(ctx, tokens[0], tokens[1])
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
			value, err := client.Get(ctx, tokens[0])
			if err != nil {
				return fmt.Errorf("could not execute GET: %w", err)
			}

			if value == nil {
				_, err = io.WriteString(conn, "+\r\n")
				if err != nil {
					return fmt.Errorf("could not send reply: %w", err)
				}
			} else {
				err := writeBulkString(conn, *value)
				if err != nil {
					return fmt.Errorf("could not write value: %w", err)
				}

				return nil
			}

			return nil
		}),
	}
}

func writeBulkString(conn io.Writer, value string) error {
	_, _ = conn.Write([]byte("$"))
	_, _ = io.WriteString(conn, strconv.Itoa(len(value)))
	_, _ = io.WriteString(conn, "\r\n")
	_, _ = io.WriteString(conn, value)

	_, err := io.WriteString(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not send reply: %w", err)
	}

	return nil
}

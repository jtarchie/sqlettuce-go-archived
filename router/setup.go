package router

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/db"
)

//nolint:funlen,cyclop
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
		"DEL": minMaxTokens(1, 0, func(tokens []string, conn io.Writer) error {
			count := 0

			for _, name := range tokens[1:] {
				err := client.Delete(ctx, name)
				if err != nil {
					_ = writeInt(conn, count)

					return fmt.Errorf("could not execute all DEL: %w", err)
				}
				count++
			}

			err := writeInt(conn, count)
			if err != nil {
				return fmt.Errorf("could not execute DEL: %w", err)
			}

			return nil
		}),
		"SET": minMaxTokens(2, 0, func(tokens []string, conn io.Writer) error {
			err := client.Set(ctx, tokens[1], tokens[2])
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
			value, err := client.Get(ctx, tokens[1])
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

func writeInt(conn io.Writer, value int) error {
	_, _ = io.WriteString(conn, ":")

	if value < 0 {
		_, _ = io.WriteString(conn, "-")
	} else {
		_, _ = io.WriteString(conn, "+")
	}

	_, _ = io.WriteString(conn, strconv.Itoa(value))

	_, err := io.WriteString(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not send int: %w", err)
	}

	return nil
}

func writeBulkString(conn io.Writer, value string) error {
	_, _ = conn.Write([]byte("$"))
	_, _ = io.WriteString(conn, strconv.Itoa(len(value)))
	_, _ = io.WriteString(conn, "\r\n")
	_, _ = io.WriteString(conn, value)

	_, err := io.WriteString(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not send bulk string: %w", err)
	}

	return nil
}

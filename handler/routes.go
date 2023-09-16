package handler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/jtarchie/sqlettus/db"
	"github.com/jtarchie/sqlettus/router"
)

//nolint:funlen,cyclop,gocognit,maintidx,gocyclo
func NewRoutes(
	ctx context.Context,
	client *db.Client,
) router.Command {
	commands := router.Command{
		"COMMAND": router.Command{
			"DOCS": router.StaticResponseRouter(router.EmptyStringResponse),
		},
		"CONFIG": router.Command{
			"GET": router.Command{
				"save":       router.StaticResponseRouter(router.EmptyStringResponse),
				"appendonly": router.StaticResponseRouter("+no\r\n"),
			},
		},
		"PING": router.StaticResponseRouter("+PONG\r\n"),
		"ECHO": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			err := writeBulkString(conn, tokens[1])
			if err != nil {
				return fmt.Errorf("could not echo message: %w", err)
			}

			return nil
		}),
		"FLUSHALL": router.CallbackRouter(func(_ []string, conn io.Writer) error {
			err := client.FlushAll()
			if err != nil {
				slog.Error("could not FLUSHALL", slog.String("error", err.Error()))
			}

			_, err = io.WriteString(conn, router.OKResponse)
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}),
		"DEL": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			count := int64(0)

			for _, name := range tokens[1:] {
				_, found, err := client.Delete(ctx, name)
				if err != nil {
					_ = writeInt(conn, count)

					return fmt.Errorf("could not execute all DEL: %w", err)
				}
				if found {
					count++
				}
			}

			err := writeInt(conn, count)
			if err != nil {
				return fmt.Errorf("could not execute DEL: %w", err)
			}

			return nil
		}),
		"SET": router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
			err := client.Set(ctx, tokens[1], tokens[2])
			if err != nil {
				return fmt.Errorf("could not execute SET: %w", err)
			}

			_, err = io.WriteString(conn, router.OKResponse)
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}),
		"GET": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			value, found, err := client.Get(ctx, tokens[1])
			if err != nil {
				return fmt.Errorf("could not execute GET: %w", err)
			}

			if !found {
				_, err = io.WriteString(conn, router.NullResponse)
				if err != nil {
					return fmt.Errorf("could not send reply: %w", err)
				}

				return nil
			}

			err = writeBulkString(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"GETDEL": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			value, found, err := client.Delete(ctx, tokens[1])
			if err != nil {
				return fmt.Errorf("could not execute GET: %w", err)
			}

			if !found {
				_, err = io.WriteString(conn, router.NullResponse)
				if err != nil {
					return fmt.Errorf("could not send reply: %w", err)
				}

				return nil
			}

			err = writeBulkString(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"APPEND": router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
			value, err := client.Append(ctx, tokens[1], tokens[2])
			if err != nil {
				return fmt.Errorf("could not execute APPEND: %w", err)
			}

			err = writeInt(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"DECR": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			value, err := client.AddInt(ctx, tokens[1], -1)
			if err != nil {
				return fmt.Errorf("could not execute DECR: %w", err)
			}

			err = writeInt(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"INCR": router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
			value, err := client.AddInt(ctx, tokens[1], 1)
			if err != nil {
				return fmt.Errorf("could not execute INCR: %w", err)
			}

			err = writeInt(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"INCRBY": router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
			incr, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				_, _ = io.WriteString(conn, "-Expected integer value to increment\r\n")

				return fmt.Errorf("could not parse integer: %w", err)
			}

			value, err := client.AddInt(ctx, tokens[1], incr)
			if err != nil {
				return fmt.Errorf("could not execute INCRBY: %w", err)
			}

			err = writeInt(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"DECRBY": router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
			incr, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				_, _ = io.WriteString(conn, "-Expected integer value to increment\r\n")

				return fmt.Errorf("could not parse integer: %w", err)
			}

			value, err := client.AddInt(ctx, tokens[1], -incr)
			if err != nil {
				return fmt.Errorf("could not execute DECRBY: %w", err)
			}

			err = writeInt(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"INCRBYFLOAT": router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
			value, err := strconv.ParseFloat(tokens[2], 64)
			if err != nil {
				_, _ = io.WriteString(conn, "-Expected float value to increment\r\n")

				return fmt.Errorf("could not parse float: %w", err)
			}

			value, err = client.AddFloat(ctx, tokens[1], value)
			if err != nil {
				return fmt.Errorf("could not execute INCRBYFLOAT: %w", err)
			}

			err = writeFloat(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
		"GETRANGE": router.MinMaxTokensRouter(3, 0, func(tokens []string, conn io.Writer) error {
			start, _ := strconv.Atoi(tokens[2])
			end, _ := strconv.Atoi(tokens[3])

			value, err := client.Substr(ctx, tokens[1], int64(start), int64(end))
			if err != nil {
				return fmt.Errorf("could not execute GETRANGE: %w", err)
			}

			err = writeBulkString(conn, value)
			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}

			return nil
		}),
	}

	commands["FLUSHDB"] = commands["FLUSHALL"]

	return commands
}

func writeFloat(conn io.Writer, value float64) error {
	_, _ = io.WriteString(conn, ",")
	_, _ = io.WriteString(conn, strconv.FormatFloat(value, 'f', 17, 64))

	_, err := io.WriteString(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not send int: %w", err)
	}

	return nil
}

func writeInt(conn io.Writer, value int64) error {
	_, _ = io.WriteString(conn, ":")
	_, _ = io.WriteString(conn, strconv.FormatInt(value, 10))

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

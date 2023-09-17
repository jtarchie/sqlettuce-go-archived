//nolint:ireturn
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

func strlenRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		value, _, err := client.Get(ctx, tokens[1])
		if err != nil {
			return fmt.Errorf("could not execute GET: %w", err)
		}

		err = writeInt(conn, int64(len(value)))
		if err != nil {
			return fmt.Errorf("could not send reply: %w", err)
		}

		return nil
	})
}

func rpushRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		value, found, err := client.ListRightPushUpsert(ctx, tokens[1], tokens[2:]...)
		if err != nil {
			return fmt.Errorf("could not execute RPUSH: %w", err)
		}

		if !found {
			err = writeError(conn, "Not an array value")
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

			return nil
		}

		err = writeInt(conn, value)
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func rpushXRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		value, err := client.ListRightPush(ctx, tokens[1], tokens[2:]...)
		if err != nil {
			_ = writeError(conn, "WRONGTYPE Operation against a key holding the wrong kind of value")

			return fmt.Errorf("could not execute RPUSHX: %w", err)
		}

		err = writeInt(conn, value)
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func echoRouter() router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		err := writeBulkString(conn, tokens[1])
		if err != nil {
			return fmt.Errorf("could not echo message: %w", err)
		}

		return nil
	})
}

func getRangeRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(3, 0, func(tokens []string, conn io.Writer) error {
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
	})
}

func incrByFloatRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		value, err := strconv.ParseFloat(tokens[2], 64)
		if err != nil {
			_ = writeError(conn, "Expected float value to increment")

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
	})
}

func decrByRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		incr, err := strconv.ParseInt(tokens[2], 10, 64)
		if err != nil {
			_ = writeError(conn, "Expected integer value to increment")

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
	})
}

func incrByRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		incr, err := strconv.ParseInt(tokens[2], 10, 64)
		if err != nil {
			_ = writeError(conn, "Expected integer value to increment")

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
	})
}

func incrRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		value, err := client.AddInt(ctx, tokens[1], 1)
		if err != nil {
			return fmt.Errorf("could not execute INCR: %w", err)
		}

		err = writeInt(conn, value)
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func decrRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		value, err := client.AddInt(ctx, tokens[1], -1)
		if err != nil {
			_ = writeError(conn, "value is not an integer or out of range")

			return fmt.Errorf("could not execute DECR: %w", err)
		}

		err = writeInt(conn, value)
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func lrangeRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(3, 0, func(tokens []string, conn io.Writer) error {
		start, _ := strconv.ParseInt(tokens[2], 10, 64)
		end, _ := strconv.ParseInt(tokens[3], 10, 64)

		values, err := client.ListRange(ctx, tokens[1], start, end)
		if err != nil {
			_ = writeError(conn, "value is not an array")

			return fmt.Errorf("could not execute LRANGE: %w", err)
		}

		_, _ = conn.Write([]byte(fmt.Sprintf("*%d\r\n", len(values))))

		for _, value := range values {
			if value == "" {
				_, err = io.WriteString(conn, router.NullResponse)
			} else {
				err = writeBulkString(conn, value)
			}

			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}
		}

		return nil
	})
}

func appendRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		value, err := client.Append(ctx, tokens[1], tokens[2])
		if err != nil {
			return fmt.Errorf("could not execute APPEND: %w", err)
		}

		err = writeInt(conn, value)
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func getDelRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		values, found, err := client.Delete(ctx, tokens[1])
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

		err = writeBulkString(conn, values[0])
		if err != nil {
			return fmt.Errorf("could not write value: %w", err)
		}

		return nil
	})
}

func mgetRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		values, err := client.MGet(ctx, tokens[1:]...)
		if err != nil {
			return fmt.Errorf("could not execute GET: %w", err)
		}

		_, _ = conn.Write([]byte(fmt.Sprintf("*%d\r\n", len(values))))

		for _, value := range values {
			if value == "" {
				_, err = io.WriteString(conn, router.NullResponse)
			} else {
				err = writeBulkString(conn, value)
			}

			if err != nil {
				return fmt.Errorf("could not write value: %w", err)
			}
		}

		return nil
	})
}

func getRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
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
	})
}

func msetRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		if len(tokens[1:])%2 != 0 {
			// require even number of tokens for key-value pairs
			_ = writeError(conn, "Expected key-value pair, not enough tokens")

			return fmt.Errorf("require even number of tokens: %w", ErrIncorrectTokens)
		}

		err := client.MSet(ctx, tokens[1:]...)
		if err != nil {
			return fmt.Errorf("could not execute MST: %w", err)
		}

		_, err = io.WriteString(conn, router.OKResponse)
		if err != nil {
			return fmt.Errorf("could not send reply: %w", err)
		}

		return nil
	})
}

func setRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(2, 0, func(tokens []string, conn io.Writer) error {
		err := client.Set(ctx, tokens[1], tokens[2])
		if err != nil {
			return fmt.Errorf("could not execute SET: %w", err)
		}

		_, err = io.WriteString(conn, router.OKResponse)
		if err != nil {
			return fmt.Errorf("could not send reply: %w", err)
		}

		return nil
	})
}

func delRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.MinMaxTokensRouter(1, 0, func(tokens []string, conn io.Writer) error {
		values, _, err := client.Delete(ctx, tokens[1:]...)
		if err != nil {
			_ = writeInt(conn, int64(len(values)))

			return fmt.Errorf("could not execute all DEL: %w", err)
		}

		err = writeInt(conn, int64(len(values)))
		if err != nil {
			return fmt.Errorf("could not execute DEL: %w", err)
		}

		return nil
	})
}

func flushAllRouter(
	ctx context.Context,
	client *db.Client,
) router.Router {
	return router.CallbackRouter(func(_ []string, conn io.Writer) error {
		err := client.FlushAll(ctx)
		if err != nil {
			slog.Error("could not FLUSHALL", slog.String("error", err.Error()))
		}

		_, err = io.WriteString(conn, router.OKResponse)
		if err != nil {
			return fmt.Errorf("could not send reply: %w", err)
		}

		return nil
	})
}

package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/jtarchie/sqlettuce/db"
	"github.com/jtarchie/sqlettuce/tcp"
)

type Handler struct {
	client *db.Client
}

func New(client *db.Client) *Handler {
	return &Handler{
		client: client,
	}
}

var _ tcp.Handler = &Handler{}

var (
	ErrIncorrectTokens = fmt.Errorf("received incorrect tokens")
	ErrNoCommandFound  = fmt.Errorf("could not determine command, none were sent")
)

func (h *Handler) OnConnection(ctx context.Context, conn io.ReadWriter) error {
	reader := bufio.NewReader(conn)
	routes := NewRoutes(ctx, h.client)

	for {
		var tokens []string

		lineCount, err := readNumber('*', reader)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read line number: %w", err)
		}

		for i := int64(0); i < lineCount; i++ {
			token, err := readString(reader)
			if err != nil {
				return fmt.Errorf("could not read token: %w", err)
			}

			tokens = append(tokens, string(token))
		}

		tokensLength := len(tokens)

		if tokensLength == 0 {
			return ErrNoCommandFound
		}

		callback, found := routes.Lookup(tokens)
		if !found {
			slog.Debug("could not found route", slog.String("tokens", strings.Join(tokens, " ")))
		}

		err = callback(tokens, conn)
		if err != nil {
			return fmt.Errorf("could not process callback: %w", err)
		}
	}
}

func readString(reader *bufio.Reader) ([]byte, error) {
	expectedLength, err := readNumber('$', reader)
	if err != nil {
		return nil, fmt.Errorf("could not read string length: %w", err)
	}

	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("could not read string: %w", err)
	}

	if int64(len(line)) != expectedLength {
		return nil, fmt.Errorf("could not read string of expected length: %w", ErrIncorrectTokens)
	}

	return line, nil
}

func readNumber(prefix byte, rw *bufio.Reader) (int64, error) {
	line, _, err := rw.ReadLine()
	if err != nil {
		return 0, fmt.Errorf("could not read line with prefix %q: %w", prefix, err)
	}

	if line[0] != prefix {
		//nolint:goerr113
		return 0, fmt.Errorf("expected line to have prefix %q, actual %q", prefix, line[0])
	}

	count := atoi(string(line[1:]))

	return count, nil
}

func atoi(num string) int64 {
	var total, index int64
loop:
	total = total*10 + int64(num[index]-'0')
	index++

	if index < int64(len(num)) {
		goto loop // avoid for loop so this function can be inlined
	}

	return total
}

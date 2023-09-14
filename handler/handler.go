package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/jtarchie/sqlettus/db"
	"github.com/jtarchie/sqlettus/router"
	"github.com/jtarchie/sqlettus/tcp"
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

func (h *Handler) OnConnection(conn io.ReadWriter) error {
	reader := bufio.NewReader(conn)
	rootRouter := router.New(h.client)

	for {
		var tokens []string

		lineCount, err := readNumber('*', reader)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read line number: %w", err)
		}

		for i := 0; i < lineCount; i++ {
			token, err := readString(reader)
			if err != nil {
				return fmt.Errorf("could not read token: %w", err)
			}

			tokens = append(tokens, string(token))
		}

		tokensLength := len(tokens)

		if tokensLength == 0 {
			return router.ErrNoCommandFound
		}

		callback, err := rootRouter.Lookup(tokens)
		if err != nil {
			return fmt.Errorf("could not execute router with tokens: %w", err)
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

	if len(line) != expectedLength {
		return nil, fmt.Errorf("could not read string of expected length: %w", router.ErrIncorrectTokens)
	}

	return line, nil
}

func readNumber(prefix byte, rw *bufio.Reader) (int, error) {
	line, _, err := rw.ReadLine()
	if err != nil {
		return 0, fmt.Errorf("could not read line with prefix %q: %w", prefix, err)
	}

	if line[0] != prefix {
		return 0, fmt.Errorf("expected line to have prefix %q, actual %q: %w", prefix, line[0], err)
	}

	count := atoi(string(line[1:]))

	return count, nil
}

func atoi(num string) int {
	var total, index int
loop:
	total = total*10 + int(num[index]-'0')
	index++

	if index < len(num) {
		goto loop // avoid for loop so this function can be inlined
	}

	return total
}

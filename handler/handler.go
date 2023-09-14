package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/jtarchie/sqlettus/db"
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

var (
	ErrNoCommandFound  = fmt.Errorf("could not determine command, none were sent")
	ErrIncorrectTokens = fmt.Errorf("received incorrect tokens")
)

//nolint:funlen,gocognit,cyclop
func (h *Handler) OnConnection(conn io.ReadWriter) error {
	reader := bufio.NewReader(conn)

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
			return ErrNoCommandFound
		}

		command := tokens[0]
		switch command {
		case "COMMAND":
			_, err = io.WriteString(conn, "*0\r\n")
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

		case "CONFIG":
			switch tokens[1] {
			case "GET":
				switch tokens[2] {
				case "save":
					_, err = io.WriteString(conn, "+\r\n")
					if err != nil {
						return fmt.Errorf("could not send reply: %w", err)
					}
				case "appendonly":
					_, err = io.WriteString(conn, "+no\r\n")
					if err != nil {
						return fmt.Errorf("could not send reply: %w", err)
					}
				default:
					return fmt.Errorf("could not CONFIG GET %q: %w", tokens[2], ErrIncorrectTokens)
				}
			default:
				return fmt.Errorf("could not CONFIG %q: %w", tokens[1], err)
			}
		case "GET":
			if tokensLength < 2 {
				return fmt.Errorf("GET requires key token: %w", err)
			}

			value, err := h.client.Get(tokens[1])
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
		case "SET":
			if tokensLength < 3 {
				return fmt.Errorf("SET requires key-value tokens: %w", err)
			}

			err := h.client.Set(tokens[1], tokens[2])
			if err != nil {
				return fmt.Errorf("could not execute SET: %w", err)
			}

			_, err = io.WriteString(conn, "+OK\r\n")
			if err != nil {
				return fmt.Errorf("could not send reply: %w", err)
			}

		default:
			return fmt.Errorf("could not find command %q: %w", command, err)
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
		return nil, fmt.Errorf("could not read string of expected length: %w", ErrIncorrectTokens)
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

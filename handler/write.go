package handler

import (
	"fmt"
	"io"
	"strconv"
)

func writeError(conn io.Writer, message string) error {
	_, _ = io.WriteString(conn, "-")
	_, _ = io.WriteString(conn, message)

	_, err := io.WriteString(conn, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write error: %w", err)
	}

	return nil
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

package tcp

import (
	"bufio"
	"fmt"
	"net"
)

func Write(port int, message string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return "", fmt.Errorf("could resolve address (localhost:%d): %w", port, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return "", fmt.Errorf("could not dial (localhost:%d): %w", port, err)
	}

	reader := bufio.NewReader(conn)

	_, err = conn.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("could not write message (localhost:%d): %w", port, err)
	}

	contents, _, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("could read from (localhost:%d): %w", port, err)
	}

	err = conn.Close()
	if err != nil {
		return "", fmt.Errorf("could not close connection (localhost:%d): %w", port, err)
	}

	return string(contents), nil
}

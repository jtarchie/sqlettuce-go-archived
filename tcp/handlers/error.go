package handlers

import (
	"errors"
	"fmt"
	"io"

	"github.com/jtarchie/sqlettus/tcp"
)

type Error struct{}

var _ tcp.Handler = &Error{}

var ErrOnConnection = errors.New("this always occurs")

func (*Error) OnConnection(_ io.ReadWriter) error {
	return fmt.Errorf("something happened: %w", ErrOnConnection)
}

package handlers

import (
	"context"
	"fmt"
	"io"

	"github.com/jtarchie/sqlettus/tcp"
)

type Echo struct{}

var _ tcp.Handler = &Echo{}

func (*Echo) OnConnection(_ context.Context, rw io.ReadWriter) error {
	_, err := io.Copy(rw, rw)
	if err != nil {
		return fmt.Errorf("could not echo: %w", err)
	}

	return nil
}
